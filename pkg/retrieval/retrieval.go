// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package retrieval provides the retrieval protocol
// implementation. The protocol is used to retrieve
// chunks over the network using forwarding-kademlia
// routing.
package retrieval

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/penguintop/penguin_bsc/pkg/accounting"
	"github.com/penguintop/penguin_bsc/pkg/cac"
	"github.com/penguintop/penguin_bsc/pkg/logging"
	"github.com/penguintop/penguin_bsc/pkg/p2p"
	"github.com/penguintop/penguin_bsc/pkg/p2p/protobuf"
	"github.com/penguintop/penguin_bsc/pkg/penguin"
	"github.com/penguintop/penguin_bsc/pkg/postage"
	"github.com/penguintop/penguin_bsc/pkg/pricer"
	pb "github.com/penguintop/penguin_bsc/pkg/retrieval/pb"
	"github.com/penguintop/penguin_bsc/pkg/soc"
	"github.com/penguintop/penguin_bsc/pkg/storage"
	"github.com/penguintop/penguin_bsc/pkg/topology"
	"github.com/penguintop/penguin_bsc/pkg/tracing"
	"golang.org/x/sync/singleflight"
)

type requestSourceContextKey struct{}

const (
	protocolName    = "retrieval"
	protocolVersion = "1.0.0"
	streamName      = "retrieval"
)

var _ Interface = (*Service)(nil)

type Interface interface {
	RetrieveChunk(ctx context.Context, addr penguin.Address) (chunk penguin.Chunk, err error)
}

type Service struct {
	addr          penguin.Address
	streamer      p2p.Streamer
	peerSuggester topology.EachPeerer
	storer        storage.Storer
	singleflight  singleflight.Group
	logger        logging.Logger
	accounting    accounting.Interface
	metrics       metrics
	pricer        pricer.Interface
	tracer        *tracing.Tracer
}

func New(addr penguin.Address, storer storage.Storer, streamer p2p.Streamer, chunkPeerer topology.EachPeerer, logger logging.Logger, accounting accounting.Interface, pricer pricer.Interface, tracer *tracing.Tracer) *Service {
	return &Service{
		addr:          addr,
		streamer:      streamer,
		peerSuggester: chunkPeerer,
		storer:        storer,
		logger:        logger,
		accounting:    accounting,
		pricer:        pricer,
		metrics:       newMetrics(),
		tracer:        tracer,
	}
}

func (s *Service) Protocol() p2p.ProtocolSpec {
	return p2p.ProtocolSpec{
		Name:    protocolName,
		Version: protocolVersion,
		StreamSpecs: []p2p.StreamSpec{
			{
				Name:    streamName,
				Handler: s.handler,
			},
		},
	}
}

const (
	maxPeers             = 5
	retrieveChunkTimeout = 10 * time.Second

	retrieveRetryIntervalDuration = 5 * time.Second
)

func (s *Service) RetrieveChunk(ctx context.Context, addr penguin.Address) (penguin.Chunk, error) {
	s.metrics.RequestCounter.Inc()

	v, err, _ := s.singleflight.Do(addr.String(), func() (interface{}, error) {
		span, logger, ctx := s.tracer.StartSpanFromContext(ctx, "retrieve-chunk", s.logger, opentracing.Tag{Key: "address", Value: addr.String()})
		defer span.Finish()

		sp := newSkipPeers()

		ticker := time.NewTicker(retrieveRetryIntervalDuration)
		defer ticker.Stop()

		var (
			peerAttempt  int
			peersResults int
			resultC      = make(chan penguin.Chunk, maxPeers)
			errC         = make(chan error, maxPeers)
		)

		for {
			if peerAttempt < maxPeers {
				peerAttempt++

				s.metrics.PeerRequestCounter.Inc()

				go func() {
					chunk, peer, err := s.retrieveChunk(ctx, addr, sp)
					if err != nil {
						if !peer.IsZero() {
							logger.Debugf("retrieval: failed to get chunk %s from peer %s: %v", addr, peer, err)
						}

						errC <- err
						return
					}

					resultC <- chunk
				}()
			} else {
				ticker.Stop()
			}

			select {
			case <-ticker.C:
				// break
			case chunk := <-resultC:
				return chunk, nil
			case <-errC:
				peersResults++
			case <-ctx.Done():
				logger.Tracef("retrieval: failed to get chunk %s: %v", addr, ctx.Err())
				return nil, fmt.Errorf("retrieval: %w", ctx.Err())
			}

			// all results received
			if peersResults >= maxPeers {
				logger.Tracef("retrieval: failed to get chunk %s", addr)
				return nil, storage.ErrNotFound
			}
		}
	})
	if err != nil {
		return nil, err
	}

	return v.(penguin.Chunk), nil
}

func (s *Service) retrieveChunk(ctx context.Context, addr penguin.Address, sp *skipPeers) (chunk penguin.Chunk, peer penguin.Address, err error) {
	startTimer := time.Now()

	v := ctx.Value(requestSourceContextKey{})
	sourcePeerAddr := penguin.Address{}
	// allow upstream requests if this node is the source of the request
	// i.e. the request was not forwarded, to improve retrieval
	// if this node is the closest to he chunk but still does not contain it
	allowUpstream := true
	if src, ok := v.(string); ok {
		sourcePeerAddr, err = penguin.ParseHexAddress(src)
		if err == nil {
			sp.Add(sourcePeerAddr)
		}
		// do not allow upstream requests if the request was forwarded to this node
		// to avoid the request loops
		allowUpstream = false
	}

	ctx, cancel := context.WithTimeout(ctx, retrieveChunkTimeout)
	defer cancel()
	peer, err = s.closestPeer(addr, sp.All(), allowUpstream)
	if err != nil {
		return nil, peer, fmt.Errorf("get closest for address %s, allow upstream %v: %w", addr.String(), allowUpstream, err)
	}

	peerPO := penguin.Proximity(s.addr.Bytes(), peer.Bytes())

	if !sourcePeerAddr.IsZero() {
		// is forwarded request
		sourceAddrPO := penguin.Proximity(sourcePeerAddr.Bytes(), addr.Bytes())
		addrPO := penguin.Proximity(peer.Bytes(), addr.Bytes())

		poGain := int(addrPO) - int(sourceAddrPO)

		s.metrics.RetrieveChunkPOGainCounter.
			WithLabelValues(strconv.Itoa(poGain)).
			Inc()
	}

	sp.Add(peer)

	// compute the peer's price for this chunk for price header
	chunkPrice := s.pricer.PeerPrice(peer, addr)

	s.logger.Tracef("retrieval: requesting chunk %s from peer %s", addr, peer)
	stream, err := s.streamer.NewStream(ctx, peer, nil, protocolName, protocolVersion, streamName)
	if err != nil {
		s.metrics.TotalErrors.Inc()
		return nil, peer, fmt.Errorf("new stream: %w", err)
	}

	defer func() {
		if err != nil {
			_ = stream.Reset()
		} else {
			go stream.FullClose()
		}
	}()

	// Reserve to see whether we can request the chunk
	err = s.accounting.Reserve(ctx, peer, chunkPrice)
	if err != nil {
		return nil, peer, err
	}
	defer s.accounting.Release(peer, chunkPrice)

	w, r := protobuf.NewWriterAndReader(stream)
	if err := w.WriteMsgWithContext(ctx, &pb.Request{
		Addr: addr.Bytes(),
	}); err != nil {
		s.metrics.TotalErrors.Inc()
		return nil, peer, fmt.Errorf("write request: %w peer %s", err, peer.String())
	}

	var d pb.Delivery
	if err := r.ReadMsgWithContext(ctx, &d); err != nil {
		s.metrics.TotalErrors.Inc()
		return nil, peer, fmt.Errorf("read delivery: %w peer %s", err, peer.String())
	}
	s.metrics.RetrieveChunkPeerPOTimer.
		WithLabelValues(strconv.Itoa(int(peerPO))).
		Observe(time.Since(startTimer).Seconds())
	s.metrics.TotalRetrieved.Inc()

	stamp := new(postage.Stamp)
	err = stamp.UnmarshalBinary(d.Stamp)
	if err != nil {
		return nil, peer, fmt.Errorf("stamp unmarshal: %w", err)
	}
	chunk = penguin.NewChunk(addr, d.Data).WithStamp(stamp)
	if !cac.Valid(chunk) {
		if !soc.Valid(chunk) {
			s.metrics.InvalidChunkRetrieved.Inc()
			s.metrics.TotalErrors.Inc()
			return nil, peer, penguin.ErrInvalidChunk
		}
	}

	// credit the peer after successful delivery
	err = s.accounting.Credit(peer, chunkPrice)
	if err != nil {
		return nil, peer, err
	}
	s.metrics.ChunkPrice.Observe(float64(chunkPrice))

	return chunk, peer, err
}

// closestPeer returns address of the peer that is closest to the chunk with
// provided address addr. This function will ignore peers with addresses
// provided in skipPeers and if allowUpstream is true, peers that are further of
// the chunk than this node is, could also be returned, allowing the upstream
// retrieve request.
func (s *Service) closestPeer(addr penguin.Address, skipPeers []penguin.Address, allowUpstream bool) (penguin.Address, error) {
	closest := penguin.Address{}
	err := s.peerSuggester.EachPeerRev(func(peer penguin.Address, po uint8) (bool, bool, error) {
		for _, a := range skipPeers {
			if a.Equal(peer) {
				return false, false, nil
			}
		}
		if closest.IsZero() {
			closest = peer
			return false, false, nil
		}
		dcmp, err := penguin.DistanceCmp(addr.Bytes(), closest.Bytes(), peer.Bytes())
		if err != nil {
			return false, false, fmt.Errorf("distance compare error. addr %s closest %s peer %s: %w", addr.String(), closest.String(), peer.String(), err)
		}
		switch dcmp {
		case 0:
			// do nothing
		case -1:
			// current peer is closer
			closest = peer
		case 1:
			// closest is already closer to chunk
			// do nothing
		}
		return false, false, nil
	})
	if err != nil {
		return penguin.Address{}, err
	}

	// check if found
	if closest.IsZero() {
		return penguin.Address{}, topology.ErrNotFound
	}
	if allowUpstream {
		return closest, nil
	}

	dcmp, err := penguin.DistanceCmp(addr.Bytes(), closest.Bytes(), s.addr.Bytes())
	if err != nil {
		return penguin.Address{}, fmt.Errorf("distance compare addr %s closest %s base address %s: %w", addr.String(), closest.String(), s.addr.String(), err)
	}
	if dcmp != 1 {
		return penguin.Address{}, topology.ErrNotFound
	}

	return closest, nil
}

func (s *Service) handler(ctx context.Context, p p2p.Peer, stream p2p.Stream) (err error) {
	w, r := protobuf.NewWriterAndReader(stream)
	defer func() {
		if err != nil {
			_ = stream.Reset()
		} else {
			_ = stream.FullClose()
		}
	}()
	var req pb.Request
	if err := r.ReadMsgWithContext(ctx, &req); err != nil {
		return fmt.Errorf("read request: %w peer %s", err, p.Address.String())
	}

	span, _, ctx := s.tracer.StartSpanFromContext(ctx, "handle-retrieve-chunk", s.logger, opentracing.Tag{Key: "address", Value: penguin.NewAddress(req.Addr).String()})
	defer span.Finish()

	ctx = context.WithValue(ctx, requestSourceContextKey{}, p.Address.String())
	addr := penguin.NewAddress(req.Addr)
	chunk, err := s.storer.Get(ctx, storage.ModeGetRequest, addr)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			// forward the request
			chunk, err = s.RetrieveChunk(ctx, addr)
			if err != nil {
				return fmt.Errorf("retrieve chunk: %w", err)
			}
		} else {
			return fmt.Errorf("get from store: %w", err)
		}
	}

	stamp, err := chunk.Stamp().MarshalBinary()
	if err != nil {
		return fmt.Errorf("stamp marshal: %w", err)
	}

	chunkPrice := s.pricer.Price(chunk.Address())
	debit := s.accounting.PrepareDebit(p.Address, chunkPrice)
	defer debit.Cleanup()

	if err := w.WriteMsgWithContext(ctx, &pb.Delivery{
		Data:  chunk.Data(),
		Stamp: stamp,
	}); err != nil {
		return fmt.Errorf("write delivery: %w peer %s", err, p.Address.String())
	}

	s.logger.Tracef("retrieval protocol debiting peer %s", p.Address.String())

	// debit price from p's balance
	return debit.Apply()
}

// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hive_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"runtime/debug"
	"strconv"
	"testing"
	"time"

	ma "github.com/multiformats/go-multiaddr"

	ab "github.com/penguintop/penguin_bsc/pkg/addressbook"
	"github.com/penguintop/penguin_bsc/pkg/crypto"
	"github.com/penguintop/penguin_bsc/pkg/hive"
	"github.com/penguintop/penguin_bsc/pkg/hive/pb"
	"github.com/penguintop/penguin_bsc/pkg/logging"
	"github.com/penguintop/penguin_bsc/pkg/p2p/protobuf"
	"github.com/penguintop/penguin_bsc/pkg/p2p/streamtest"
	"github.com/penguintop/penguin_bsc/pkg/pen"
	"github.com/penguintop/penguin_bsc/pkg/penguin"
	"github.com/penguintop/penguin_bsc/pkg/penguin/test"
	"github.com/penguintop/penguin_bsc/pkg/statestore/mock"
)

func TestHandlerRateLimit(t *testing.T) {

	logger := logging.New(ioutil.Discard, 0)
	statestore := mock.NewStateStore()
	addressbook := ab.New(statestore)
	networkID := uint64(1)

	addressbookclean := ab.New(mock.NewStateStore())

	// create a hive server that handles the incoming stream
	server := hive.New(nil, addressbookclean, networkID, logger)

	serverAddress := test.RandomAddress()

	// setup the stream recorder to record stream data
	serverRecorder := streamtest.New(
		streamtest.WithProtocols(server.Protocol()),
		streamtest.WithBaseAddr(serverAddress),
	)

	peers := make([]penguin.Address, hive.LimitBurst+1)
	for i := range peers {

		underlay, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/" + strconv.Itoa(i))
		if err != nil {
			t.Fatal(err)
		}
		pk, err := crypto.GenerateSecp256k1Key()
		if err != nil {
			t.Fatal(err)
		}
		signer := crypto.NewDefaultSigner(pk)
		overlay, err := crypto.NewOverlayAddress(pk.PublicKey, networkID)
		if err != nil {
			t.Fatal(err)
		}
		penAddr, err := pen.NewAddress(signer, underlay, overlay, networkID)
		if err != nil {
			t.Fatal(err)
		}

		err = addressbook.Put(penAddr.Overlay, *penAddr)
		if err != nil {
			t.Fatal(err)
		}
		peers[i] = penAddr.Overlay
	}

	// create a hive client that will do broadcast
	client := hive.New(serverRecorder, addressbook, networkID, logger)
	err := client.BroadcastPeers(context.Background(), serverAddress, peers...)
	if err != nil {
		t.Fatal(err)
	}

	// // get a record for this stream
	rec, err := serverRecorder.Records(serverAddress, "hive", "1.0.0", "peers")
	if err != nil {
		t.Fatal(err)
	}

	lastRec := rec[len(rec)-1]
	if !errors.Is(lastRec.Err(), hive.ErrRateLimitExceeded) {
		t.Fatal(err)
	}
}

func TestBroadcastPeers(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	logger := logging.New(ioutil.Discard, 0)
	statestore := mock.NewStateStore()
	addressbook := ab.New(statestore)
	networkID := uint64(1)

	// populate all expected and needed random resources for 2 full batches
	// tests cases that uses fewer resources can use sub-slices of this data
	var penAddresses []pen.Address
	var overlays []penguin.Address
	var wantMsgs []pb.Peers

	for i := 0; i < 2; i++ {
		wantMsgs = append(wantMsgs, pb.Peers{Peers: []*pb.BzzAddress{}})
	}

	for i := 0; i < 2*hive.MaxBatchSize; i++ {
		underlay, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/" + strconv.Itoa(i))
		if err != nil {
			t.Fatal(err)
		}
		pk, err := crypto.GenerateSecp256k1Key()
		if err != nil {
			t.Fatal(err)
		}
		signer := crypto.NewDefaultSigner(pk)
		overlay, err := crypto.NewOverlayAddress(pk.PublicKey, networkID)
		if err != nil {
			t.Fatal(err)
		}
		penAddr, err := pen.NewAddress(signer, underlay, overlay, networkID)
		if err != nil {
			t.Fatal(err)
		}

		penAddresses = append(penAddresses, *penAddr)
		overlays = append(overlays, penAddr.Overlay)
		err = addressbook.Put(penAddr.Overlay, *penAddr)
		if err != nil {
			t.Fatal(err)
		}

		wantMsgs[i/hive.MaxBatchSize].Peers = append(wantMsgs[i/hive.MaxBatchSize].Peers, &pb.BzzAddress{
			Overlay:   penAddresses[i].Overlay.Bytes(),
			Underlay:  penAddresses[i].Underlay.Bytes(),
			Signature: penAddresses[i].Signature,
		})
	}

	testCases := map[string]struct {
		addresee         penguin.Address
		peers            []penguin.Address
		wantMsgs         []pb.Peers
		wantOverlays     []penguin.Address
		wantBzzAddresses []pen.Address
	}{
		"OK - single record": {
			addresee:         penguin.MustParseHexAddress("ca1e9f3938cc1425c6061b96ad9eb93e134dfe8734ad490164ef20af9d1cf59c"),
			peers:            []penguin.Address{overlays[0]},
			wantMsgs:         []pb.Peers{{Peers: wantMsgs[0].Peers[:1]}},
			wantOverlays:     []penguin.Address{overlays[0]},
			wantBzzAddresses: []pen.Address{penAddresses[0]},
		},
		"OK - single batch - multiple records": {
			addresee:         penguin.MustParseHexAddress("ca1e9f3938cc1425c6061b96ad9eb93e134dfe8734ad490164ef20af9d1cf59c"),
			peers:            overlays[:15],
			wantMsgs:         []pb.Peers{{Peers: wantMsgs[0].Peers[:15]}},
			wantOverlays:     overlays[:15],
			wantBzzAddresses: penAddresses[:15],
		},
		"OK - single batch - max number of records": {
			addresee:         penguin.MustParseHexAddress("ca1e9f3938cc1425c6061b96ad9eb93e134dfe8734ad490164ef20af9d1cf59c"),
			peers:            overlays[:hive.MaxBatchSize],
			wantMsgs:         []pb.Peers{{Peers: wantMsgs[0].Peers[:hive.MaxBatchSize]}},
			wantOverlays:     overlays[:hive.MaxBatchSize],
			wantBzzAddresses: penAddresses[:hive.MaxBatchSize],
		},
		"OK - multiple batches": {
			addresee:         penguin.MustParseHexAddress("ca1e9f3938cc1425c6061b96ad9eb93e134dfe8734ad490164ef20af9d1cf59c"),
			peers:            overlays[:hive.MaxBatchSize+10],
			wantMsgs:         []pb.Peers{{Peers: wantMsgs[0].Peers}, {Peers: wantMsgs[1].Peers[:10]}},
			wantOverlays:     overlays[:hive.MaxBatchSize+10],
			wantBzzAddresses: penAddresses[:hive.MaxBatchSize+10],
		},
		"OK - multiple batches - max number of records": {
			addresee:         penguin.MustParseHexAddress("ca1e9f3938cc1425c6061b96ad9eb93e134dfe8734ad490164ef20af9d1cf59c"),
			peers:            overlays[:2*hive.MaxBatchSize],
			wantMsgs:         []pb.Peers{{Peers: wantMsgs[0].Peers}, {Peers: wantMsgs[1].Peers}},
			wantOverlays:     overlays[:2*hive.MaxBatchSize],
			wantBzzAddresses: penAddresses[:2*hive.MaxBatchSize],
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			addressbookclean := ab.New(mock.NewStateStore())

			// create a hive server that handles the incoming stream
			server := hive.New(nil, addressbookclean, networkID, logger)

			// setup the stream recorder to record stream data
			recorder := streamtest.New(
				streamtest.WithProtocols(server.Protocol()),
			)

			// create a hive client that will do broadcast
			client := hive.New(recorder, addressbook, networkID, logger)
			if err := client.BroadcastPeers(context.Background(), tc.addresee, tc.peers...); err != nil {
				t.Fatal(err)
			}

			// get a record for this stream
			records, err := recorder.Records(tc.addresee, "hive", "1.0.0", "peers")
			if err != nil {
				t.Fatal(err)
			}
			if l := len(records); l != len(tc.wantMsgs) {
				t.Fatalf("got %v records, want %v", l, len(tc.wantMsgs))
			}

			// there is a one record per batch (wantMsg)
			for i, record := range records {
				messages, err := readAndAssertPeersMsgs(record.In(), 1)
				if err != nil {
					t.Fatal(err)
				}

				if fmt.Sprint(messages[0]) != fmt.Sprint(tc.wantMsgs[i]) {
					t.Errorf("Messages got %v, want %v", messages, tc.wantMsgs)
				}
			}

			expectOverlaysEventually(t, addressbookclean, tc.wantOverlays)
			expectBzzAddresessEventually(t, addressbookclean, tc.wantBzzAddresses)
		})
	}
}

func expectOverlaysEventually(t *testing.T, exporter ab.Interface, wantOverlays []penguin.Address) {
	var (
		overlays []penguin.Address
		err      error
		isIn     = func(a penguin.Address, addrs []penguin.Address) bool {
			for _, v := range addrs {
				if a.Equal(v) {
					return true
				}
			}
			return false
		}
	)

	for i := 0; i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
		overlays, err = exporter.Overlays()
		if err != nil {
			t.Fatal(err)
		}

		if len(overlays) == len(wantOverlays) {
			break
		}
	}
	if len(overlays) != len(wantOverlays) {
		debug.PrintStack()
		t.Fatal("timed out waiting for overlays")
	}

	for _, v := range wantOverlays {
		if !isIn(v, overlays) {
			t.Errorf("overlay %s expected but not found", v.String())
		}
	}

	if t.Failed() {
		t.Errorf("overlays got %v, want %v", overlays, wantOverlays)
	}
}

func expectBzzAddresessEventually(t *testing.T, exporter ab.Interface, wantBzzAddresses []pen.Address) {
	var (
		addresses []pen.Address
		err       error

		isIn = func(a pen.Address, addrs []pen.Address) bool {
			for _, v := range addrs {
				if a.Equal(&v) {
					return true
				}
			}
			return false
		}
	)

	for i := 0; i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
		addresses, err = exporter.Addresses()
		if err != nil {
			t.Fatal(err)
		}

		if len(addresses) == len(wantBzzAddresses) {
			break
		}
	}
	if len(addresses) != len(wantBzzAddresses) {
		debug.PrintStack()
		t.Fatal("timed out waiting for pen addresses")
	}

	for _, v := range wantBzzAddresses {
		if !isIn(v, addresses) {
			t.Errorf("address %s expected but not found", v.Overlay.String())
		}
	}

	if t.Failed() {
		t.Errorf("pen addresses got %v, want %v", addresses, wantPenAddresses)
	}
}

func readAndAssertPeersMsgs(in []byte, expectedLen int) ([]pb.Peers, error) {
	messages, err := protobuf.ReadMessages(
		bytes.NewReader(in),
		func() protobuf.Message {
			return new(pb.Peers)
		},
	)

	if err != nil {
		return nil, err
	}

	if len(messages) != expectedLen {
		return nil, fmt.Errorf("got %v messages, want %v", len(messages), expectedLen)
	}

	var peers []pb.Peers
	for _, m := range messages {
		peers = append(peers, *m.(*pb.Peers))
	}

	return peers, nil
}

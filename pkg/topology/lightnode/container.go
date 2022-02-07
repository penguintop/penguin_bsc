// Copyright 2021 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lightnode

import (
	"context"
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/penguintop/penguin_bsc/pkg/p2p"
	"github.com/penguintop/penguin_bsc/pkg/penguin"
	"github.com/penguintop/penguin_bsc/pkg/topology"
	"github.com/penguintop/penguin_bsc/pkg/topology/pslice"
)

type Container struct {
	base              penguin.Address
	connectedPeers    *pslice.PSlice
	disconnectedPeers *pslice.PSlice
	peerMu            sync.Mutex
}

func NewContainer(base penguin.Address) *Container {
	return &Container{
		base:              base,
		connectedPeers:    pslice.New(1, base),
		disconnectedPeers: pslice.New(1, base),
	}
}

func (c *Container) Connected(ctx context.Context, peer p2p.Peer) {
	c.peerMu.Lock()
	defer c.peerMu.Unlock()

	addr := peer.Address
	c.connectedPeers.Add(addr)
	c.disconnectedPeers.Remove(addr)
}

func (c *Container) Disconnected(peer p2p.Peer) {
	c.peerMu.Lock()
	defer c.peerMu.Unlock()

	addr := peer.Address
	if found := c.connectedPeers.Exists(addr); found {
		c.connectedPeers.Remove(addr)
		c.disconnectedPeers.Add(addr)
	}
}

func (c *Container) Count() int {
	return c.connectedPeers.Length()
}

func (c *Container) RandomPeer(not penguin.Address) (penguin.Address, error) {
	c.peerMu.Lock()
	defer c.peerMu.Unlock()
	var (
		cnt   = big.NewInt(int64(c.Count()))
		addr  = penguin.ZeroAddress
		count = int64(0)
	)

PICKPEER:
	i, e := rand.Int(rand.Reader, cnt)
	if e != nil {
		return penguin.ZeroAddress, e
	}
	i64 := i.Int64()

	count = 0
	_ = c.connectedPeers.EachBinRev(func(peer penguin.Address, _ uint8) (bool, bool, error) {
		if count == i64 {
			addr = peer
			return true, false, nil
		}
		count++
		return false, false, nil
	})

	if addr.Equal(not) {
		goto PICKPEER
	}

	return addr, nil
}

func (c *Container) PeerInfo() topology.BinInfo {
	return topology.BinInfo{
		BinPopulation:     uint(c.connectedPeers.Length()),
		BinConnected:      uint(c.connectedPeers.Length()),
		DisconnectedPeers: peersInfo(c.disconnectedPeers),
		ConnectedPeers:    peersInfo(c.connectedPeers),
	}
}

func peersInfo(s *pslice.PSlice) []*topology.PeerInfo {
	if s.Length() == 0 {
		return nil
	}
	peers := make([]*topology.PeerInfo, 0, s.Length())
	_ = s.EachBin(func(addr penguin.Address, po uint8) (bool, bool, error) {
		peers = append(peers, &topology.PeerInfo{Address: addr})
		return false, false, nil
	})
	return peers
}

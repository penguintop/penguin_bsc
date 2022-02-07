// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debugapi

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/multiformats/go-multiaddr"
	"github.com/penguintop/penguin_bsc/pkg/jsonhttp"
	"github.com/penguintop/penguin_bsc/pkg/p2p"
	"github.com/penguintop/penguin_bsc/pkg/penguin"
)

type peerConnectResponse struct {
	Address string `json:"address"`
}

func (s *Service) peerConnectHandler(w http.ResponseWriter, r *http.Request) {
	addr, err := multiaddr.NewMultiaddr("/" + mux.Vars(r)["multi-address"])
	if err != nil {
		s.logger.Debugf("debug api: peer connect: parse multiaddress: %v", err)
		jsonhttp.BadRequest(w, err)
		return
	}

	penAddr, err := s.p2p.Connect(r.Context(), addr)
	if err != nil {
		s.logger.Debugf("debug api: peer connect %s: %v", addr, err)
		s.logger.Errorf("unable to connect to peer %s", addr)
		jsonhttp.InternalServerError(w, err)
		return
	}

	jsonhttp.OK(w, peerConnectResponse{
		Address: penAddr.Overlay.String(),
	})
}

func (s *Service) peerDisconnectHandler(w http.ResponseWriter, r *http.Request) {
	addr := mux.Vars(r)["address"]
	penguinAddr, err := penguin.ParseHexAddress(addr)
	if err != nil {
		s.logger.Debugf("debug api: parse peer address %s: %v", addr, err)
		jsonhttp.BadRequest(w, "invalid peer address")
		return
	}

	if err := s.p2p.Disconnect(penguinAddr); err != nil {
		s.logger.Debugf("debug api: peer disconnect %s: %v", addr, err)
		if errors.Is(err, p2p.ErrPeerNotFound) {
			jsonhttp.BadRequest(w, "peer not found")
			return
		}
		s.logger.Errorf("unable to disconnect peer %s", addr)
		jsonhttp.InternalServerError(w, err)
		return
	}

	jsonhttp.OK(w, nil)
}

type peersResponse struct {
	Peers []p2p.Peer `json:"peers"`
}

func (s *Service) peersHandler(w http.ResponseWriter, r *http.Request) {
	jsonhttp.OK(w, peersResponse{
		Peers: s.p2p.Peers(),
	})
}

func (s *Service) blocklistedPeersHandler(w http.ResponseWriter, r *http.Request) {
	peers, err := s.p2p.BlocklistedPeers()
	if err != nil {
		s.logger.Debugf("debug api: blocklisted peers: %v", err)
		jsonhttp.InternalServerError(w, nil)
		return
	}

	jsonhttp.OK(w, peersResponse{
		Peers: peers,
	})
}

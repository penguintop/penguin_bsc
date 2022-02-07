// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debugapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/penguintop/penguin_bsc/pkg/jsonhttp"
	"github.com/penguintop/penguin_bsc/pkg/penguin"
	"github.com/penguintop/penguin_bsc/pkg/storage"
)

func (s *Service) hasChunkHandler(w http.ResponseWriter, r *http.Request) {
	addr, err := penguin.ParseHexAddress(mux.Vars(r)["address"])
	if err != nil {
		s.logger.Debugf("debug api: parse chunk address: %v", err)
		jsonhttp.BadRequest(w, "bad address")
		return
	}

	has, err := s.storer.Has(r.Context(), addr)
	if err != nil {
		s.logger.Debugf("debug api: localstore has: %v", err)
		jsonhttp.BadRequest(w, err)
		return
	}

	if !has {
		jsonhttp.NotFound(w, nil)
		return
	}
	jsonhttp.OK(w, nil)
}

func (s *Service) removeChunk(w http.ResponseWriter, r *http.Request) {
	addr, err := penguin.ParseHexAddress(mux.Vars(r)["address"])
	if err != nil {
		s.logger.Debugf("debug api: parse chunk address: %v", err)
		jsonhttp.BadRequest(w, "bad address")
		return
	}

	has, err := s.storer.Has(r.Context(), addr)
	if err != nil {
		s.logger.Debugf("debug api: localstore remove: %v", err)
		jsonhttp.BadRequest(w, err)
		return
	}

	if !has {
		jsonhttp.OK(w, nil)
		return
	}

	err = s.storer.Set(r.Context(), storage.ModeSetRemove, addr)
	if err != nil {
		s.logger.Debugf("debug api: localstore remove: %v", err)
		jsonhttp.InternalServerError(w, err)
		return
	}
	jsonhttp.OK(w, nil)
}
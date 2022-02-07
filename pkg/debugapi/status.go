// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debugapi

import (
	"net/http"

	"github.com/penguintop/penguin_bsc"
	"github.com/penguintop/penguin_bsc/pkg/jsonhttp"
)

type statusResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	jsonhttp.OK(w, statusResponse{
		Status:  "ok",
		Version: pen_bsc.Version,
	})
}

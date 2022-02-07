// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package soc

import (
	"github.com/penguintop/penguin_bsc/pkg/penguin"
)

// Valid checks if the chunk is a valid single-owner chunk.
func Valid(ch penguin.Chunk) bool {
	s, err := FromChunk(ch)
	if err != nil {
		return false
	}

	address, err := s.address()
	if err != nil {
		return false
	}
	return ch.Address().Equal(address)
}

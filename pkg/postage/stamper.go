// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package postage

import (
	"errors"

	"github.com/penguintop/penguin_bsc/pkg/crypto"
	"github.com/penguintop/penguin_bsc/pkg/penguin"
)

var (
	// ErrBucketFull is the error when a collision bucket is full.
	ErrBucketFull = errors.New("bucket full")
)

// Stamper can issue stamps from the given address.
type Stamper interface {
	Stamp(penguin.Address) (*Stamp, error)
}

// stamper connects a stampissuer with a signer.
// A stamper is created for each upload session.
type stamper struct {
	issuer *StampIssuer
	signer crypto.Signer
}

// NewStamper constructs a Stamper.
func NewStamper(st *StampIssuer, signer crypto.Signer) Stamper {
	return &stamper{st, signer}
}

// Stamp takes chunk, see if the chunk can included in the batch and
// signs it with the owner of the batch of this Stamp issuer.
func (st *stamper) Stamp(addr penguin.Address) (*Stamp, error) {
	toSign, err := toSignDigest(addr, st.issuer.batchID)
	if err != nil {
		return nil, err
	}
	sig, err := st.signer.Sign(toSign)
	if err != nil {
		return nil, err
	}
	if err := st.issuer.inc(addr); err != nil {
		return nil, err
	}
	return NewStamp(st.issuer.batchID, sig), nil
}

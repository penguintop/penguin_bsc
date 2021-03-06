// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstore_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"sync/atomic"
	"testing"
	"time"

	"github.com/penguintop/penguin_bsc/pkg/logging"
	"github.com/penguintop/penguin_bsc/pkg/netstore"
	"github.com/penguintop/penguin_bsc/pkg/penguin"
	postagetesting "github.com/penguintop/penguin_bsc/pkg/postage/testing"
	"github.com/penguintop/penguin_bsc/pkg/pss"
	"github.com/penguintop/penguin_bsc/pkg/recovery"
	"github.com/penguintop/penguin_bsc/pkg/sctx"
	"github.com/penguintop/penguin_bsc/pkg/storage"
	"github.com/penguintop/penguin_bsc/pkg/storage/mock"
)

var chunkData = []byte("mockdata")
var chunkStamp = postagetesting.MustNewStamp()

// TestNetstoreRetrieval verifies that a chunk is asked from the network whenever
// it is not found locally
func TestNetstoreRetrieval(t *testing.T) {
	retrieve, store, nstore := newRetrievingNetstore(nil, noopValidStamp)
	addr := penguin.MustParseHexAddress("000001")
	_, err := nstore.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}
	if !retrieve.called {
		t.Fatal("retrieve request not issued")
	}
	if retrieve.callCount != 1 {
		t.Fatalf("call count %d", retrieve.callCount)
	}
	if !retrieve.addr.Equal(addr) {
		t.Fatalf("addresses not equal. got %s want %s", retrieve.addr, addr)
	}

	// store should have the chunk now
	d, err := store.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(d.Data(), chunkData) {
		t.Fatal("chunk data not equal to expected data")
	}

	// check that the second call does not result in another retrieve request
	d, err = nstore.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}

	if retrieve.callCount != 1 {
		t.Fatalf("call count %d", retrieve.callCount)
	}
	if !bytes.Equal(d.Data(), chunkData) {
		t.Fatal("chunk data not equal to expected data")
	}

}

// TestNetstoreNoRetrieval verifies that a chunk is not requested from the network
// whenever it is found locally.
func TestNetstoreNoRetrieval(t *testing.T) {
	retrieve, store, nstore := newRetrievingNetstore(nil, noopValidStamp)
	addr := penguin.MustParseHexAddress("000001")

	// store should have the chunk in advance
	_, err := store.Put(context.Background(), storage.ModePutUpload, penguin.NewChunk(addr, chunkData))
	if err != nil {
		t.Fatal(err)
	}

	c, err := nstore.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}
	if retrieve.called {
		t.Fatal("retrieve request issued but shouldn't")
	}
	if retrieve.callCount != 0 {
		t.Fatalf("call count %d", retrieve.callCount)
	}
	if !bytes.Equal(c.Data(), chunkData) {
		t.Fatal("chunk data mismatch")
	}
}

func TestRecovery(t *testing.T) {
	callbackWasCalled := make(chan bool, 1)
	rec := &mockRecovery{
		callbackC: callbackWasCalled,
	}

	retrieve, _, nstore := newRetrievingNetstore(rec.recovery, noopValidStamp)
	addr := penguin.MustParseHexAddress("deadbeef")
	retrieve.failure = true
	ctx := context.Background()
	ctx = sctx.SetTargets(ctx, "be, cd")

	_, err := nstore.Get(ctx, storage.ModeGetRequest, addr)
	if err != nil && !errors.Is(err, netstore.ErrRecoveryAttempt) {
		t.Fatal(err)
	}

	select {
	case <-callbackWasCalled:
		break
	case <-time.After(100 * time.Millisecond):
		t.Fatal("recovery callback was not called")
	}
}

func TestInvalidRecoveryFunction(t *testing.T) {
	retrieve, _, nstore := newRetrievingNetstore(nil, noopValidStamp)
	addr := penguin.MustParseHexAddress("deadbeef")
	retrieve.failure = true
	ctx := context.Background()
	ctx = sctx.SetTargets(ctx, "be, cd")

	_, err := nstore.Get(ctx, storage.ModeGetRequest, addr)
	if err != nil && err.Error() != "chunk not found" {
		t.Fatal(err)
	}
}

func TestInvalidPostageStamp(t *testing.T) {
	f := func(c penguin.Chunk, _ []byte) (penguin.Chunk, error) {
		return nil, errors.New("invalid postage stamp")
	}
	retrieve, store, nstore := newRetrievingNetstore(nil, f)
	addr := penguin.MustParseHexAddress("000001")
	_, err := nstore.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}
	if !retrieve.called {
		t.Fatal("retrieve request not issued")
	}
	if retrieve.callCount != 1 {
		t.Fatalf("call count %d", retrieve.callCount)
	}
	if !retrieve.addr.Equal(addr) {
		t.Fatalf("addresses not equal. got %s want %s", retrieve.addr, addr)
	}

	// store should have the chunk now
	d, err := store.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(d.Data(), chunkData) {
		t.Fatal("chunk data not equal to expected data")
	}

	if mode := store.GetModePut(addr); mode != storage.ModePutRequestCache {
		t.Fatalf("wanted ModePutRequestCache but got %s", mode)
	}

	// check that the second call does not result in another retrieve request
	d, err = nstore.Get(context.Background(), storage.ModeGetRequest, addr)
	if err != nil {
		t.Fatal(err)
	}

	if retrieve.callCount != 1 {
		t.Fatalf("call count %d", retrieve.callCount)
	}
	if !bytes.Equal(d.Data(), chunkData) {
		t.Fatal("chunk data not equal to expected data")
	}
}

// returns a mock retrieval protocol, a mock local storage and a netstore
func newRetrievingNetstore(rec recovery.Callback, validStamp func(penguin.Chunk, []byte) (penguin.Chunk, error)) (ret *retrievalMock, mockStore *mock.MockStorer, ns storage.Storer) {
	retrieve := &retrievalMock{}
	store := mock.NewStorer()
	logger := logging.New(ioutil.Discard, 0)
	return retrieve, store, netstore.New(store, validStamp, rec, retrieve, logger)
}

type retrievalMock struct {
	called    bool
	callCount int32
	failure   bool
	addr      penguin.Address
}

func (r *retrievalMock) RetrieveChunk(ctx context.Context, addr penguin.Address) (chunk penguin.Chunk, err error) {
	if r.failure {
		return nil, fmt.Errorf("chunk not found")
	}
	r.called = true
	atomic.AddInt32(&r.callCount, 1)
	r.addr = addr
	return penguin.NewChunk(addr, chunkData).WithStamp(chunkStamp), nil
}

type mockRecovery struct {
	callbackC chan bool
}

// Send mocks the pss Send function
func (mr *mockRecovery) recovery(chunkAddress penguin.Address, targets pss.Targets) {
	mr.callbackC <- true
}

func (r *mockRecovery) RetrieveChunk(ctx context.Context, addr penguin.Address) (chunk penguin.Chunk, err error) {
	return nil, fmt.Errorf("chunk not found")
}

var noopValidStamp = func(c penguin.Chunk, _ []byte) (penguin.Chunk, error) {
	return c, nil
}

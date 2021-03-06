// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pusher_test

import (
	"context"
	"errors"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/penguintop/penguin_bsc/pkg/crypto"
	statestore "github.com/penguintop/penguin_bsc/pkg/statestore/mock"

	"github.com/penguintop/penguin_bsc/pkg/localstore"
	"github.com/penguintop/penguin_bsc/pkg/logging"
	"github.com/penguintop/penguin_bsc/pkg/penguin"
	"github.com/penguintop/penguin_bsc/pkg/pusher"
	"github.com/penguintop/penguin_bsc/pkg/pushsync"
	pushsyncmock "github.com/penguintop/penguin_bsc/pkg/pushsync/mock"
	"github.com/penguintop/penguin_bsc/pkg/storage"
	testingc "github.com/penguintop/penguin_bsc/pkg/storage/testing"
	"github.com/penguintop/penguin_bsc/pkg/tags"
	"github.com/penguintop/penguin_bsc/pkg/topology/mock"
)

// no of times to retry to see if we have received response from pushsync
var noOfRetries = 20

// Wrap the actual storer to intercept the modeSet that the pusher will call when a valid receipt is received
type Store struct {
	storage.Storer
	internalStorer storage.Storer
	modeSet        map[string]storage.ModeSet
	modeSetMu      *sync.Mutex

	closed              bool
	setBeforeCloseCount int
	setAfterCloseCount  int
}

// Override the Set function to capture the ModeSetSync
func (s *Store) Set(ctx context.Context, mode storage.ModeSet, addrs ...penguin.Address) error {
	s.modeSetMu.Lock()
	defer s.modeSetMu.Unlock()

	if s.closed {
		s.setAfterCloseCount++
	} else {
		s.setBeforeCloseCount++
	}

	for _, addr := range addrs {
		s.modeSet[addr.String()] = mode
	}
	return nil
}

func (s *Store) Close() error {
	s.modeSetMu.Lock()
	defer s.modeSetMu.Unlock()

	s.closed = true
	return s.internalStorer.Close()
}

// TestSendChunkToPushSync sends a chunk to pushsync to be sent ot its closest peer and get a receipt.
// once the receipt is got this check to see if the localstore is updated to see if the chunk is set
// as ModeSetSync status.
func TestSendChunkToSyncWithTag(t *testing.T) {
	// create a trigger  and a closestpeer
	triggerPeer := penguin.MustParseHexAddress("6000000000000000000000000000000000000000000000000000000000000000")
	closestPeer := penguin.MustParseHexAddress("f000000000000000000000000000000000000000000000000000000000000000")

	key, _ := crypto.GenerateSecp256k1Key()
	signer := crypto.NewDefaultSigner(key)

	pushSyncService := pushsyncmock.New(func(ctx context.Context, chunk penguin.Chunk) (*pushsync.Receipt, error) {
		signature, _ := signer.Sign(chunk.Address().Bytes())
		receipt := &pushsync.Receipt{
			Address:   penguin.NewAddress(chunk.Address().Bytes()),
			Signature: signature,
		}
		return receipt, nil
	})

	mtags, p, storer := createPusher(t, triggerPeer, pushSyncService, mock.WithClosestPeer(closestPeer))
	defer storer.Close()
	defer p.Close()

	ta, err := mtags.Create(1)
	if err != nil {
		t.Fatal(err)
	}

	chunk := testingc.GenerateTestRandomChunk().WithTagID(ta.Uid)

	_, err = storer.Put(context.Background(), storage.ModePutUpload, chunk)
	if err != nil {
		t.Fatal(err)
	}

	// Check is the chunk is set as synced in the DB.
	for i := 0; i < noOfRetries; i++ {
		// Give some time for chunk to be pushed and receipt to be received
		time.Sleep(10 * time.Millisecond)

		err = checkIfModeSet(chunk.Address(), storage.ModeSetSync, storer)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatal(err)
	}

	if ta.Get(tags.StateSynced) != 1 {
		t.Fatalf("tags error")
	}
}

// TestSendChunkToPushSyncWithoutTag is similar to TestSendChunkToPushSync, excep that the tags are not
// present to simulate pen api withotu splitter condition
func TestSendChunkToPushSyncWithoutTag(t *testing.T) {
	chunk := testingc.GenerateTestRandomChunk()

	// create a trigger  and a closestpeer
	triggerPeer := penguin.MustParseHexAddress("6000000000000000000000000000000000000000000000000000000000000000")
	closestPeer := penguin.MustParseHexAddress("f000000000000000000000000000000000000000000000000000000000000000")

	key, _ := crypto.GenerateSecp256k1Key()
	signer := crypto.NewDefaultSigner(key)

	pushSyncService := pushsyncmock.New(func(ctx context.Context, chunk penguin.Chunk) (*pushsync.Receipt, error) {
		signature, _ := signer.Sign(chunk.Address().Bytes())
		receipt := &pushsync.Receipt{
			Address:   penguin.NewAddress(chunk.Address().Bytes()),
			Signature: signature,
		}
		return receipt, nil
	})

	_, p, storer := createPusher(t, triggerPeer, pushSyncService, mock.WithClosestPeer(closestPeer))
	defer storer.Close()
	defer p.Close()

	_, err := storer.Put(context.Background(), storage.ModePutUpload, chunk)
	if err != nil {
		t.Fatal(err)
	}

	// Check is the chunk is set as synced in the DB.
	for i := 0; i < noOfRetries; i++ {
		// Give some time for chunk to be pushed and receipt to be received
		time.Sleep(50 * time.Millisecond)

		err = checkIfModeSet(chunk.Address(), storage.ModeSetSync, storer)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}

// TestSendChunkAndReceiveInvalidReceipt sends a chunk to pushsync to be sent ot its closest peer and
// get a invalid receipt (not with the address of the chunk sent). The test makes sure that this error
// is received and the ModeSetSync is not set for the chunk.
func TestSendChunkAndReceiveInvalidReceipt(t *testing.T) {
	chunk := testingc.GenerateTestRandomChunk()

	// create a trigger  and a closestpeer
	triggerPeer := penguin.MustParseHexAddress("6000000000000000000000000000000000000000000000000000000000000000")
	closestPeer := penguin.MustParseHexAddress("f000000000000000000000000000000000000000000000000000000000000000")

	pushSyncService := pushsyncmock.New(func(ctx context.Context, chunk penguin.Chunk) (*pushsync.Receipt, error) {
		return nil, errors.New("invalid receipt")
	})

	_, p, storer := createPusher(t, triggerPeer, pushSyncService, mock.WithClosestPeer(closestPeer))
	defer storer.Close()
	defer p.Close()

	_, err := storer.Put(context.Background(), storage.ModePutUpload, chunk)
	if err != nil {
		t.Fatal(err)
	}

	// Check is the chunk is set as synced in the DB.
	for i := 0; i < noOfRetries; i++ {
		// Give some time for chunk to be pushed and receipt to be received
		time.Sleep(10 * time.Millisecond)

		err = checkIfModeSet(chunk.Address(), storage.ModeSetSync, storer)
		if err != nil {
			continue
		}
	}
	if err == nil {
		t.Fatalf("chunk not syned error expected")
	}
}

// TestSendChunkAndTimeoutinReceivingReceipt sends a chunk to pushsync to be sent ot its closest peer and
// expects a timeout to get instead of getting a receipt. The test makes sure that timeout error
// is received and the ModeSetSync is not set for the chunk.
func TestSendChunkAndTimeoutinReceivingReceipt(t *testing.T) {
	chunk := testingc.GenerateTestRandomChunk()

	// create a trigger  and a closestpeer
	triggerPeer := penguin.MustParseHexAddress("6000000000000000000000000000000000000000000000000000000000000000")
	closestPeer := penguin.MustParseHexAddress("f000000000000000000000000000000000000000000000000000000000000000")

	key, _ := crypto.GenerateSecp256k1Key()
	signer := crypto.NewDefaultSigner(key)

	pushSyncService := pushsyncmock.New(func(ctx context.Context, chunk penguin.Chunk) (*pushsync.Receipt, error) {
		time.Sleep(1 * time.Second)
		signature, _ := signer.Sign(chunk.Address().Bytes())
		receipt := &pushsync.Receipt{
			Address:   penguin.NewAddress(chunk.Address().Bytes()),
			Signature: signature,
		}
		return receipt, nil
	})

	_, p, storer := createPusher(t, triggerPeer, pushSyncService, mock.WithClosestPeer(closestPeer))
	defer storer.Close()
	defer p.Close()

	_, err := storer.Put(context.Background(), storage.ModePutUpload, chunk)
	if err != nil {
		t.Fatal(err)
	}

	// Check is the chunk is set as synced in the DB.
	for i := 0; i < noOfRetries; i++ {
		// Give some time for chunk to be pushed and receipt to be received
		time.Sleep(10 * time.Millisecond)

		err = checkIfModeSet(chunk.Address(), storage.ModeSetSync, storer)
		if err != nil {
			continue
		}
	}
	if err == nil {
		t.Fatalf("chunk not syned error expected")
	}
}

func TestPusherClose(t *testing.T) {
	// create a trigger  and a closestpeer
	triggerPeer := penguin.MustParseHexAddress("6000000000000000000000000000000000000000000000000000000000000000")
	closestPeer := penguin.MustParseHexAddress("f000000000000000000000000000000000000000000000000000000000000000")

	var (
		goFuncStartedC    = make(chan struct{})
		pusherClosedC     = make(chan struct{})
		goFuncAfterCloseC = make(chan struct{})
	)

	defer func() {
		close(goFuncStartedC)
		close(pusherClosedC)
		close(goFuncAfterCloseC)
	}()

	key, _ := crypto.GenerateSecp256k1Key()
	signer := crypto.NewDefaultSigner(key)

	pushSyncService := pushsyncmock.New(func(ctx context.Context, chunk penguin.Chunk) (*pushsync.Receipt, error) {
		goFuncStartedC <- struct{}{}
		<-goFuncAfterCloseC
		signature, _ := signer.Sign(chunk.Address().Bytes())
		receipt := &pushsync.Receipt{
			Address:   penguin.NewAddress(chunk.Address().Bytes()),
			Signature: signature,
		}
		return receipt, nil
	})

	_, p, storer := createPusher(t, triggerPeer, pushSyncService, mock.WithClosestPeer(closestPeer))

	chunk := testingc.GenerateTestRandomChunk()

	_, err := storer.Put(context.Background(), storage.ModePutUpload, chunk)
	if err != nil {
		t.Fatal(err)
	}

	storer.modeSetMu.Lock()
	if storer.closed == true {
		t.Fatal("store should not be closed")
	}
	if storer.setBeforeCloseCount != 0 {
		t.Fatalf("store 'Set' called %d times before close, expected 0", storer.setBeforeCloseCount)
	}
	if storer.setAfterCloseCount != 0 {
		t.Fatalf("store 'Set' called %d times after close, expected 0", storer.setAfterCloseCount)
	}
	storer.modeSetMu.Unlock()

	select {
	case <-goFuncStartedC:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting to start worker job")
	}

	// close in the background
	go func() {
		p.Close()
		storer.Close()
		pusherClosedC <- struct{}{}
	}()

	select {
	case <-pusherClosedC:
	case <-time.After(2 * time.Second):
		// internal 5 second timeout that waits for all pending push operations to terminate
	}

	storer.modeSetMu.Lock()
	if storer.setBeforeCloseCount != 0 {
		t.Fatalf("store 'Set' called %d times before close, expected 0", storer.setBeforeCloseCount)
	}
	if storer.setAfterCloseCount != 0 {
		t.Fatalf("store 'Set' called %d times after close, expected 0", storer.setAfterCloseCount)
	}
	storer.modeSetMu.Unlock()

	select {
	case goFuncAfterCloseC <- struct{}{}:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for chunk")
	}

	// we need this to allow some goroutines to complete
	time.Sleep(100 * time.Millisecond)

	storer.modeSetMu.Lock()
	if storer.closed != true {
		t.Fatal("store should be closed")
	}
	if storer.setBeforeCloseCount != 1 {
		t.Fatalf("store 'Set' called %d times before close, expected 1", storer.setBeforeCloseCount)
	}
	if storer.setAfterCloseCount != 0 {
		t.Fatalf("store 'Set' called %d times after close, expected 0", storer.setAfterCloseCount)
	}
	storer.modeSetMu.Unlock()

	// should be closed by now
	select {
	case <-pusherClosedC:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting to close pusher")
	}
}

func createPusher(t *testing.T, addr penguin.Address, pushSyncService pushsync.PushSyncer, mockOpts ...mock.Option) (*tags.Tags, *pusher.Service, *Store) {
	t.Helper()
	logger := logging.New(ioutil.Discard, 0)
	storer, err := localstore.New("", addr.Bytes(), nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	mockStatestore := statestore.NewStateStore()
	mtags := tags.NewTags(mockStatestore, logger)
	pusherStorer := &Store{
		Storer:         storer,
		internalStorer: storer,
		modeSet:        make(map[string]storage.ModeSet),
		modeSetMu:      &sync.Mutex{},
	}
	peerSuggester := mock.NewTopologyDriver(mockOpts...)

	pusherService := pusher.New(1, pusherStorer, peerSuggester, pushSyncService, mtags, logger, nil)
	return mtags, pusherService, pusherStorer
}

func checkIfModeSet(addr penguin.Address, mode storage.ModeSet, storer *Store) error {
	var found bool
	storer.modeSetMu.Lock()
	defer storer.modeSetMu.Unlock()

	for k, v := range storer.modeSet {
		if addr.String() == k {
			found = true
			if v != mode {
				return errors.New("chunk mode is not properly set as synced")
			}
		}
	}
	if !found {
		return errors.New("Chunk not synced")
	}
	return nil
}

// To avoid timeout during race testing
// cd pkg/pusher
// go test -race -count 1000 -timeout 60m .

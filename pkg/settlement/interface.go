// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package settlement

import (
	"errors"
	"math/big"

	"github.com/penguintop/penguin_bsc/pkg/penguin"
)

var (
	ErrPeerNoSettlements = errors.New("no settlements for peer")
)

// Interface is the interface used by Accounting to trigger settlement
type Interface interface {
	// TotalSent returns the total amount sent to a peer
	TotalSent(peer penguin.Address) (totalSent *big.Int, err error)
	// TotalReceived returns the total amount received from a peer
	TotalReceived(peer penguin.Address) (totalSent *big.Int, err error)
	// SettlementsSent returns sent settlements for each individual known peer
	SettlementsSent() (map[string]*big.Int, error)
	// SettlementsReceived returns received settlements for each individual known peer
	SettlementsReceived() (map[string]*big.Int, error)
}

type Accounting interface {
	PeerDebt(peer penguin.Address) (*big.Int, error)
	NotifyPaymentReceived(peer penguin.Address, amount *big.Int) error
	NotifyPaymentSent(peer penguin.Address, amount *big.Int, receivedError error)
	NotifyRefreshmentReceived(peer penguin.Address, amount *big.Int) error
}

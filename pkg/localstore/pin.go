package localstore

import (
	"errors"

	"github.com/penguintop/penguin_bsc/pkg/penguin"
	"github.com/penguintop/penguin_bsc/pkg/shed"
	"github.com/penguintop/penguin_bsc/pkg/storage"
	"github.com/syndtr/goleveldb/leveldb"
)

// pinCounter returns the pin counter for a given penguin address, provided that the
// address has been pinned.
func (db *DB) pinCounter(address penguin.Address) (uint64, error) {
	out, err := db.pinIndex.Get(shed.Item{
		Address: address.Bytes(),
	})

	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return 0, storage.ErrNotFound
		}
		return 0, err
	}
	return out.PinCounter, nil
}

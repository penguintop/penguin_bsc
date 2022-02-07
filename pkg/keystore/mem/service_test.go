// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem_test

import (
	"testing"

	"github.com/penguintop/penguin_bsc/pkg/keystore/mem"
	"github.com/penguintop/penguin_bsc/pkg/keystore/test"
)

func TestService(t *testing.T) {
	test.Service(t, mem.New())
}

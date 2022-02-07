// Copyright 2020 The Penguin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package penguin_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"

	"github.com/penguintop/penguin_bsc/pkg/penguin"
)

func TestAddress(t *testing.T) {
	for _, tc := range []struct {
		name    string
		hex     string
		want    penguin.Address
		wantErr error
	}{
		{
			name: "blank",
			hex:  "",
			want: penguin.ZeroAddress,
		},
		{
			name:    "odd",
			hex:     "0",
			wantErr: hex.ErrLength,
		},
		{
			name: "zero",
			hex:  "00",
			want: penguin.NewAddress([]byte{0}),
		},
		{
			name: "one",
			hex:  "01",
			want: penguin.NewAddress([]byte{1}),
		},
		{
			name: "arbitrary",
			hex:  "35a26b7bb6455cbabe7a0e05aafbd0b8b26feac843e3b9a649468d0ea37a12b2",
			want: penguin.NewAddress([]byte{0x35, 0xa2, 0x6b, 0x7b, 0xb6, 0x45, 0x5c, 0xba, 0xbe, 0x7a, 0xe, 0x5, 0xaa, 0xfb, 0xd0, 0xb8, 0xb2, 0x6f, 0xea, 0xc8, 0x43, 0xe3, 0xb9, 0xa6, 0x49, 0x46, 0x8d, 0xe, 0xa3, 0x7a, 0x12, 0xb2}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a, err := penguin.ParseHexAddress(tc.hex)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got error %v, want %v", err, tc.wantErr)
			}
			if a.String() != tc.want.String() {
				t.Errorf("got address %#v, want %#v", a, tc.want)
			}
			if !a.Equal(tc.want) {
				t.Errorf("address %v not equal to %v", a, tc.want)
			}
			if a.IsZero() != tc.want.IsZero() {
				t.Errorf("got address as zero=%v, want zero=%v", a.IsZero(), tc.want.IsZero())
			}
		})
	}
}

func TestAddress_jsonMarshalling(t *testing.T) {
	a1 := penguin.MustParseHexAddress("24798dd5a470e927fa")

	b, err := json.Marshal(a1)
	if err != nil {
		t.Fatal(err)
	}

	var a2 penguin.Address
	if err := json.Unmarshal(b, &a2); err != nil {
		t.Fatal(err)
	}

	if !a1.Equal(a2) {
		t.Error("unmarshalled address is not equal to the original")
	}
}

func TestAddress_MemberOf(t *testing.T) {
	a1 := penguin.MustParseHexAddress("24798dd5a470e927fa")
	a2 := penguin.MustParseHexAddress("24798dd5a470e927fa")
	a3 := penguin.MustParseHexAddress("24798dd5a470e927fb")
	a4 := penguin.MustParseHexAddress("24798dd5a470e927fc")

	set1 := []penguin.Address{a2, a3}
	if !a1.MemberOf(set1) {
		t.Fatal("expected addr as member")
	}

	set2 := []penguin.Address{a3, a4}
	if a1.MemberOf(set2) {
		t.Fatal("expected addr not member")
	}

}

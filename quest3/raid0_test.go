package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRAID0(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "Test RAID0",
			input: `ABCDEFGHIJKLMNOP`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			raid := NewRAID0(4)
			bs := []byte(tc.input)
			err := raid.Write(bs)
			assert.NoError(t, err)

			ret, err := raid.Read(len(bs))
			t.Log(string(ret))
			assert.NoError(t, err)
			assert.Equal(t, bs, ret)
		})
	}
}

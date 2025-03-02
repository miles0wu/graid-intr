package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRAID(t *testing.T) {
	testCases := []struct {
		name           string
		raidBuilder    func(*testing.T) RAID
		inputStr       string
		failDiskIdx    int
		reconstructErr error
	}{
		{
			name: "Test RAID0",
			raidBuilder: func(t *testing.T) RAID {
				return NewRAID0(2)
			},
			inputStr:       `ABCDEFGHIJKLMNOP`,
			failDiskIdx:    1,
			reconstructErr: errors.New("raid 0 not support reconstruct"),
		},
		{
			name: "Test RAID1",
			raidBuilder: func(t *testing.T) RAID {
				return NewRAID1(4)
			},
			failDiskIdx: 1,
			inputStr:    `ABCDEFGHIJKLMNOP`,
		},
		{
			name: "Test RAID5",
			raidBuilder: func(t *testing.T) RAID {
				raid, err := NewRAID5(3)
				assert.NoError(t, err)
				return raid
			},
			failDiskIdx: 1,
			inputStr:    `ABCDEFGHIJKLMNOP`,
		},
		{
			name: "Test RAID10",
			raidBuilder: func(t *testing.T) RAID {
				raid, err := NewRAID10(4, 2)
				assert.NoError(t, err)
				return raid
			},
			failDiskIdx: 1,
			inputStr:    `ABCDEFGHIJKLMNOP`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			raid := tc.raidBuilder(t)
			bs := []byte(tc.inputStr)

			// write
			err = raid.Write(bs)
			assert.NoError(t, err)

			// read
			ret, err := raid.Read(len(bs))
			t.Log(string(ret))
			assert.NoError(t, err)
			assert.Equal(t, bs, ret)

			// fail the disk
			err = raid.FailDisk(tc.failDiskIdx)
			assert.NoError(t, err)

			// reconstruct
			err = raid.Reconstruct()
			assert.Equal(t, tc.reconstructErr, err)
			if err != nil {
				return
			}

			// read after reconstruct
			ret, err = raid.Read(len(bs))
			t.Log(string(ret))
			assert.NoError(t, err)
			assert.Equal(t, bs, ret)
		})
	}
}

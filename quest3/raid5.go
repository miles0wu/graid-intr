package main

import (
	"errors"
	"fmt"
)

type RAID5 struct {
	numDisks int
	disks    []*Disk
}

func NewRAID5(numDisks int) (RAID, error) {
	if numDisks < 3 {
		return nil, fmt.Errorf("numDisks must be at least 3: %d", numDisks)
	}
	disks := make([]*Disk, numDisks)
	for i := range disks {
		disks[i] = NewDisk()
	}
	return &RAID5{
		numDisks: numDisks,
		disks:    disks,
	}, nil
}

func (r *RAID5) CalculateParityDisk(stripeIndex int) int {
	return r.numDisks - ((stripeIndex % r.numDisks) + 1)
}

func (r *RAID5) CalculateParity(blocks []byte) byte {
	var ret byte
	for _, ch := range blocks {
		ret ^= ch
	}
	return ret
}

func (r *RAID5) Write(data []byte) error {
	dataLen := len(data)
	if dataLen == 0 {
		return errors.New("no data to write")
	}

	blocksPerStripe := r.numDisks - 1

	for stripeIndex := 0; stripeIndex*blocksPerStripe < dataLen; stripeIndex++ {
		start := stripeIndex * blocksPerStripe
		end := start + blocksPerStripe
		if end > dataLen {
			end = dataLen
		}

		stripeData := make([]byte, end-start)
		copy(stripeData, data[start:end])
		for len(stripeData) < blocksPerStripe {
			stripeData = append(stripeData, '\x00')
		}
		parity := r.CalculateParity(stripeData)
		parityDisk := r.CalculateParityDisk(stripeIndex)

		// combine data blocks and parity as full stripe data blocks
		fullStripe := make([]byte, 0, r.numDisks)
		dataIndex := 0
		for i := 0; i < r.numDisks; i++ {
			if i == parityDisk {
				fullStripe = append(fullStripe, parity)
			} else {
				fullStripe = append(fullStripe, stripeData[dataIndex])
				dataIndex++
			}
		}
		// write full stripe data blocks
		for idx, block := range fullStripe {
			if r.disks[idx] == nil {
				return fmt.Errorf("disk %d is missing", idx)
			}
			_ = r.disks[idx].WriteChar(block)
		}
	}

	return nil
}

func (r *RAID5) Read(length int) ([]byte, error) {
	ret := make([]byte, 0, length)

	blocksPerStripe := r.numDisks - 1
	for stripeIndex := 0; stripeIndex*blocksPerStripe < length; stripeIndex++ {
		parityDisk := r.CalculateParityDisk(stripeIndex)
		for idx, disk := range r.disks {
			if idx == parityDisk {
				continue
			}
			b, _ := disk.ReadChar(stripeIndex)
			ret = append(ret, b)
		}
	}

	return ret, nil
}

func (r *RAID5) FailDisk(diskNum int) error {
	if diskNum < 0 || diskNum > r.numDisks {
		return errors.New("fail disk error: index out of range")
	}
	r.disks[diskNum] = nil
	return nil
}

func (r *RAID5) Reconstruct() error {
	missingDiskIdx := -1
	missingCnt := 0
	for idx, disk := range r.disks {
		if disk == nil {
			missingDiskIdx = idx
			missingCnt++
		}
	}
	if missingCnt > 1 {
		return errors.New("cannot reconstruct disk: too many disks missing")
	}

	if missingDiskIdx == -1 {
		return nil
	}

	newDisk := NewDisk()
	var numStripes int
	for _, disk := range r.disks {
		if disk != nil {
			numStripes = len(disk.data)
			break
		}
	}
	for stripeIndex := 0; stripeIndex < numStripes; stripeIndex++ {
		var recovered byte
		for idx, disk := range r.disks {
			if idx == missingDiskIdx {
				continue
			}
			b, _ := disk.ReadChar(stripeIndex)
			recovered ^= b
		}
		_ = newDisk.WriteChar(recovered)
	}
	r.disks[missingDiskIdx] = newDisk
	return nil
}

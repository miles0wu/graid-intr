package main

import (
	"errors"
)

type RAID10 struct {
	mirrorGroups []*RAID1
	numGroups    int
	groupSize    int
}

func NewRAID10(totalDisks int, groupSize int) (*RAID10, error) {
	if totalDisks%groupSize != 0 {
		return nil, errors.New("invalid total disks size")
	}
	numGroups := totalDisks / groupSize
	mirrorGroups := make([]*RAID1, numGroups)
	for i := 0; i < numGroups; i++ {
		mirrorGroups[i] = NewRAID1(groupSize).(*RAID1)
	}
	return &RAID10{
		mirrorGroups: mirrorGroups,
		numGroups:    numGroups,
		groupSize:    groupSize,
	}, nil
}

func (r *RAID10) Write(data []byte) error {
	for pos, b := range data {
		groupIndex := pos % r.numGroups
		if err := r.mirrorGroups[groupIndex].WriteChar(b); err != nil {
			return err
		}
	}
	return nil
}

func (r *RAID10) Read(length int) ([]byte, error) {
	ret := make([]byte, 0, length)
	for pos := 0; pos < length; pos++ {
		groupIndex := pos % r.numGroups
		stripeIndex := pos / r.numGroups
		b, err := r.mirrorGroups[groupIndex].ReadChar(stripeIndex)
		if err != nil {
			return nil, err
		}
		ret = append(ret, b)
	}
	return ret, nil
}

func (r *RAID10) FailDisk(index int) error {
	groupIndex := index / r.groupSize
	diskIndex := index % r.groupSize
	if groupIndex < 0 || groupIndex >= len(r.mirrorGroups) {
		return errors.New("invalid mirror group index")
	}
	return r.mirrorGroups[groupIndex].FailDisk(diskIndex)
}

func (r *RAID10) Reconstruct() error {
	for _, group := range r.mirrorGroups {
		if err := group.Reconstruct(); err != nil {
			return err
		}
	}
	return nil
}

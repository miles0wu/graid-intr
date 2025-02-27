package main

import (
	"errors"
	"math/rand"
)

type Raid1 struct {
	disks []*Disk
}

func (r *Raid1) Write(bs []byte) error {
	for _, d := range r.disks {
		_ = d.Write(bs)
	}
	return nil
}

func (r *Raid1) Read(length int) ([]byte, error) {
	numDisks := len(r.disks)
	start := rand.Intn(numDisks)
	i := (start + 1) % numDisks
	for i != start {
		ret, err := r.disks[i].Read(length)
		if err != nil {
			i++
			i %= numDisks
			continue
		}
		return ret, nil
	}
	return nil, errors.New("all disk was broken")
}

func (r *Raid1) ClearDisks(diskIndex int) error {
	if diskIndex < 0 || diskIndex >= len(r.disks) {
		return errors.New("invalid disk index")
	}
	r.disks[diskIndex] = NewDisk()
	return nil
}

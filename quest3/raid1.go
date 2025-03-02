package main

import (
	"errors"
	"math/rand"
)

type RAID1 struct {
	disks    []*Disk
	numDisks int
}

func NewRAID1(numDisks int) RAID {
	disks := make([]*Disk, numDisks)
	for i := range disks {
		disks[i] = NewDisk()
	}
	return &RAID1{
		numDisks: numDisks,
		disks:    disks,
	}
}

func (r *RAID1) Write(bs []byte) error {
	for _, d := range r.disks {
		_ = d.Write(bs)
	}
	return nil
}

func (r *RAID1) Read(length int) ([]byte, error) {
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

func (r *RAID1) FailDisk(diskIndex int) error {
	if diskIndex < 0 || diskIndex >= len(r.disks) {
		return errors.New("invalid disk index")
	}
	r.disks[diskIndex] = nil
	return nil
}

func (r *RAID1) Reconstruct() error {
	var canReconstruct bool
	var source *Disk
	for _, d := range r.disks {
		if d != nil {
			canReconstruct = true
			source = d
			break
		}
	}
	if !canReconstruct {
		return errors.New("no healthy disks to reconstruct")
	}

	for _, d := range r.disks {
		if d != nil {
			continue
		}
		d = NewDisk()
		_ = d.Write(source.data)
	}
	return nil
}

func (r *RAID1) WriteChar(b byte) error {
	for _, d := range r.disks {
		if err := d.WriteChar(b); err != nil {
			return err
		}
	}
	return nil
}

func (r *RAID1) ReadChar(index int) (byte, error) {
	for _, d := range r.disks {
		if d != nil {
			return d.ReadChar(index)
		}
	}
	return 0, errors.New("no healthy disks to read")
}

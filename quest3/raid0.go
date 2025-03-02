package main

import "errors"

type RAID0 struct {
	numDisks int
	disks    []*Disk
}

func NewRAID0(numDisks int) RAID {
	disks := make([]*Disk, numDisks)
	for i := range disks {
		disks[i] = NewDisk()
	}
	return &RAID0{
		numDisks: numDisks,
		disks:    disks,
	}
}

func (r *RAID0) Write(data []byte) error {
	dataLen := len(data)
	if dataLen == 0 {
		return errors.New("no data to write")
	}

	pos := 0
	for pos < dataLen {
		disk := r.disks[pos%r.numDisks]
		err := disk.WriteChar(data[pos])
		if err != nil {
			return err
		}
		pos++
	}
	return nil
}

func (r *RAID0) Read(length int) ([]byte, error) {
	ret := make([]byte, 0, length)
	pos := 0

	// Read in stripe order
	for pos < length {
		disk := r.disks[pos%r.numDisks]
		ch, err := disk.ReadChar(pos / r.numDisks)
		if err != nil {
			return ret, err
		}
		ret = append(ret, ch)
		pos++
	}
	return ret, nil
}

func (r *RAID0) FailDisk(diskIndex int) error {
	if diskIndex < 0 || diskIndex >= len(r.disks) {
		return errors.New("invalid disk index")
	}
	r.disks[diskIndex] = NewDisk()
	return nil
}

func (r *RAID0) Reconstruct() error {
	return errors.New("raid 0 not support reconstruct")
}

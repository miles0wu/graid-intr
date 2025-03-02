package main

import (
	"errors"
)

type Disk struct {
	data []byte
}

func NewDisk() *Disk {
	return &Disk{
		data: []byte{},
	}
}

func (d *Disk) Write(data []byte) error {
	d.data = append(d.data, data...)
	return nil
}

func (d *Disk) WriteChar(ch byte) error {
	d.data = append(d.data, ch)
	return nil
}

func (d *Disk) ReadChar(pos int) (byte, error) {
	if pos >= len(d.data) {
		return 0, errors.New("index out of range")
	}
	return d.data[pos], nil
}

func (d *Disk) Read(length int) ([]byte, error) {
	if length > len(d.data) {
		return nil, errors.New("read error: length exceeds stored data")
	}
	return d.data[:length], nil
}

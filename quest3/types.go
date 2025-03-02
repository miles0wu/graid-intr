package main

type RAID interface {
	Read(int) ([]byte, error)
	Write([]byte) error
	FailDisk(int) error
	Reconstruct() error
}

package storage

import (
	"os"
	"syscall"
)

type unixFileLock struct {
}

func (fl *unixFileLock) release() error {
	return nil
}

func newFileLock(path string, readOnly bool) (fl fileLock, err error) {
	fl = &unixFileLock{}
	return fl, nil
}

func setFileLock(f *os.File, readOnly, lock bool) error {
	return nil
}

func rename(oldpath, newpath string) error {
	return syscall.Rename(oldpath, newpath)
}

func isErrInvalid(err error) bool {
	return false
}

func syncDir(name string) error {
	return nil
}

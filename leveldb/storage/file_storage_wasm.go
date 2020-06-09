package storage

import (
	"os"
	"sync"
	"syscall"
)

// implementation using internal mutex since browserfs doesn't support real locks
type unixFileLock struct {
	f *os.File
	m sync.Mutex
	p string
}

var (
	lockMapMtx sync.RWMutex
	lockMap = make(map[*os.File]*unixFileLock)
	pathMap = make(map[string]*os.File)
)

func (fl *unixFileLock) release() error {
	if fl == nil || fl.f == nil {
		return nil
	}
	if err := setFileLock(fl.f, false, false); err != nil {
		return err
	}
	lockMapMtx.Lock()
	delete(lockMap, fl.f)
	delete(pathMap, fl.p)
	lockMapMtx.Unlock()
	return fl.f.Close()
}

func newFileLock(path string, readOnly bool) (fl fileLock, err error) {
	lockMapMtx.RLock()
	f := pathMap[path]
	lockMapMtx.RUnlock()

	if f == nil {
		var flag int
		if readOnly {
			flag = os.O_RDONLY
		} else {
			flag = os.O_RDWR|os.O_EXCL|os.O_APPEND
		}
		f, err = os.OpenFile(path, flag, 0)
		if os.IsNotExist(err) {
			f, err = os.OpenFile(path, flag|os.O_CREATE, 0644)
		}
		if err != nil {
			return
		}
	}

	err = setFileLock(f, readOnly, true)
	if err != nil {
		f.Close()
		return
	}
	lockMapMtx.RLock()
	fl = lockMap[f]
	lockMapMtx.RUnlock()
	return
}

func setFileLock(f *os.File, readOnly, lock bool) (err error) {
	lockMapMtx.Lock()
	fl, ok := lockMap[f]
	if !ok {
		fl = &unixFileLock{f: f, p: f.Name()}
		pathMap[fl.p] = f
	}
	lockMapMtx.Unlock()

	defer func() {
		if err := recover(); err != nil {
			err = nil
		}
	}()

	if lock {
		fl.m.Lock()
	} else {
		fl.m.Unlock()
	}

	return
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

package filemutex

import (
	"sync"
	"syscall"
)

// FileMutex is similar to sync.RWMutex, but also synchronizes across processes.
// This implementation is based on flock syscall.
type FileMutex struct {
	mux sync.RWMutex
	fd  int
}

func New(filename string) (*FileMutex, error) {
	fd, err := syscall.Open(filename, syscall.O_CREAT|syscall.O_RDONLY, 0750)
	if err != nil {
		return nil, err
	}
	return &FileMutex{fd: fd}, nil
}

func (m *FileMutex) Lock() {
	m.mux.Lock()
	if err := syscall.Flock(m.fd, syscall.LOCK_EX); err != nil {
		panic(err)
	}
}

func (m *FileMutex) Unlock() {
	if err := syscall.Flock(m.fd, syscall.LOCK_UN); err != nil {
		panic(err)
	}
	m.mux.Unlock()
}

func (m *FileMutex) RLock() {
	m.mux.RLock()
	if err := syscall.Flock(m.fd, syscall.LOCK_SH); err != nil {
		panic(err)
	}
}

func (m *FileMutex) RUnlock() {
	if err := syscall.Flock(m.fd, syscall.LOCK_UN); err != nil {
		panic(err)
	}
	m.mux.RUnlock()
}

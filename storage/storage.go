package storage

import (
	"encoding/json"
	"io"
	"sync"
)

type Data struct {
	Meetups map[string]*Meetup `json:"meetups"`
}

// Storage holds all currently active meetings.
type Storage struct {
	data  *Data
	mutex *sync.RWMutex
}

// NewStorage creates a ready to use data store.
func NewStorage() *Storage {
	return &Storage{
		data: &Data{
			make(map[string]*Meetup),
		},
		mutex: &sync.RWMutex{},
	}
}

// Read allows you to read from the storage, giving you a read-lock. While you
// can technically still write, you should refrain from doing so. This ensures
// that no read and write occurrs at the same time, while multiple reads are
// allowed and desirable.
func (s *Storage) Read(accessor func(data *Data)) {
	s.mutex.RLock()
	accessor(s.data)
	defer s.mutex.RUnlock()
}

// ReadWrite allows you to read from or write to the storage, giving you a
// read-write-lock. This ensures that no read or write occurrs during your
// write.
func (s *Storage) ReadWrite(accessor func(data *Data)) {
	s.mutex.Lock()
	accessor(s.data)
	defer s.mutex.Unlock()
}

// Load reads all meetings from the given sourcen and closes it afterwards.
func (s *Storage) Load(source io.ReadCloser) error {
	defer source.Close()
	return json.NewDecoder(source).Decode(&s.data)
}

// Persist writes out the current state and closes the target writer.
func (s *Storage) Persist(target io.WriteCloser) error {
	defer target.Close()
	return json.NewEncoder(target).Encode(s.data)
}

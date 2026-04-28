package store

import "time"


type Store struct{
	data map[string]Entry
}

type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
	HasExpiry bool
}

func New() *Store {
	return &Store{
		data: make(map[string]Entry),
	}
}

func (s *Store) Set(key string, val interface{}) {
	s.data[key] = Entry{Value: val}
}

func (s *Store) SetWithTTL(key string, val interface{}, ttl time.Duration) {
	s.data[key] = Entry{Value: val, ExpiresAt: time.Now().Add(ttl), HasExpiry: true}
}

func (s *Store) Get(key string) (interface{}, bool) {
	entry, ok := s.data[key]

	if (!ok){
		return nil, false
	}

	if entry.HasExpiry && time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Value, ok
}

func (s *Store) Exists(key string) (bool) {
	_ , ok := s.Get(key)
	return ok
}

func (s *Store) Delete(key string) (bool) {
	_, ok := s.data[key]
	if ok {
		delete(s.data, key)
	}
	return ok
}

func (s *Store) TTL(key string) (int) {
	entry, ok := s.data[key]

	if (!ok) {
		return -2
	}

	if (!entry.HasExpiry) {
		return -1
	}
	
	remaining := time.Until(entry.ExpiresAt)
	if remaining < 0 {
		return -2
	}

	return int(remaining.Seconds())
}
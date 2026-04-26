package store


type Store struct{
	data map[string]interface{}
}

func New() *Store {
	return &Store{
		data: make(map[string]interface{}),
	}
}

func (s *Store) Set(key string, val interface{}) {
	s.data[key] = val
}

func (s *Store) Get(key string) (interface{}, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *Store) Exists(key string) (bool) {
	_ , ok := s.data[key]
	return ok
}

func (s *Store) Delete(key string) (bool) {
	_, ok := s.data[key]
	if ok {
		delete(s.data, key)
	}
	return ok
}
package store

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/bbolt"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record version conflict")
)

var bucketNames = [][]byte{[]byte("records"), []byte("audit"), []byte("workflows"), []byte("attachments"), []byte("meta")}

type DB struct {
	db       *bbolt.DB
	mu       sync.RWMutex
	sequence int64
}

func Open(path string) (*DB, error) {
	if path == "" {
		path = filepath.Join(".", "aroma.db")
	}
	bolt, err := bbolt.Open(path, 0o600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	result := &DB{db: bolt}
	err = bolt.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, createErr := tx.CreateBucketIfNotExists(name); createErr != nil {
				return createErr
			}
		}
		return nil
	})
	if err != nil {
		bolt.Close()
		return nil, err
	}
	return result, nil
}

func (s *DB) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *DB) nextSequence() int64 {
	s.sequence++
	return s.sequence
}

func (s *DB) transaction(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(fn)
}

func (s *DB) view(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.View(fn)
}

func putJSON(bucket *bbolt.Bucket, key string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(key), payload)
}

func getJSON(bucket *bbolt.Bucket, key string, target any) error {
	data := bucket.Get([]byte(key))
	if data == nil {
		return ErrNotFound
	}
	return json.Unmarshal(data, target)
}

func (s *DB) Health() bool {
	return s.db != nil
}

package main

import (
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// PersistentCache provides a simple key-value store with TTL using BadgerDB
type PersistentCache struct {
	db *badger.DB
}

// NewCache initializes a new PersistentCache at the given path
func NewCache(path string) *PersistentCache {
	opts := badger.DefaultOptions(path).
		WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("Failed to open BadgerDB at %s: %v", path, err)
	}
	return &PersistentCache{db: db}
}

// Set stores a key-value pair with the specified TTL
func (c *PersistentCache) Set(key string, val []byte, ttl time.Duration) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(
			badger.NewEntry([]byte(key), val).WithTTL(ttl),
		)
	})
}

// Get retrieves the value for a given key. It returns the value, a boolean indicating if the key was found, and an error if any.
func (c *PersistentCache) Get(key string) ([]byte, bool, error) {
	var out []byte

	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		out, err = item.ValueCopy(nil)
		return err
	})

	if err == badger.ErrKeyNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return out, true, nil
}

// Close closes the BadgerDB instance
func (c *PersistentCache) Close() {
	if err := c.db.Close(); err != nil {
		log.Printf("Error closing persistent cache: %v", err)
	}
}

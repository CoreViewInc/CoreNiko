package registry

import (
	"fmt"
	"log"
	bolt "go.etcd.io/bbolt"
	shared "github.com/CoreViewInc/CoreNiko/shared" 
)

// ImageRegistry is a simple abstraction for a Docker image registry using bbolt.
type KanikoImageRegistry struct {
	db *bolt.DB
}

// NewImageRegistry creates a new ImageRegistry and opens the database.
func newImageRegistry(dbPath string) (shared.Registry, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err // Return the error instead of calling log.Fatal, giving the caller the option to handle it.
	}
	return &KanikoImageRegistry{db: db}, nil
}

// Initialize sets up the required buckets in the database.
func (ir *KanikoImageRegistry) Initialize() error {
	return ir.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Images"))
		return err
	})
}

// RecordImage adds a new image record to the registry with the given tag and location.
func (ir *KanikoImageRegistry) RecordImage(tag, location string) error {
	return ir.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Images"))
		return b.Put([]byte(tag), []byte(location))
	})
}

// GetImageLocation retrieves the location of an image with the given tag.
func (ir *KanikoImageRegistry) GetImageLocation(tag string) (string, error) {
	var location string
	err := ir.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Images"))
		v := b.Get([]byte(tag))
		if v == nil {
			return fmt.Errorf("image with tag %s not found", tag)
		}
		location = string(v)
		return nil
	})
	return location, err
}

// Close closes the database.
func (ir *KanikoImageRegistry) Close() error {
	return ir.db.Close()
}

func New(filename string) (shared.Registry,error) {
	registry,err := newImageRegistry(filename)
	if err!=nil{
		return nil,err
	}
	defer registry.Close()
	err = registry.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	return registry,nil
}
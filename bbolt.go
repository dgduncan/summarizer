package summarizer

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

const (
	path = "mydb.db"

	ShowsBucket    = "episodes"
	SummarysBucket = "summaries"
)

func Set(b, k, v string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b))
		return b.Put([]byte(k), []byte(v))
	})
}

func Get(b, k string, db *bolt.DB) (string, error) {
	var val []byte
	err := db.View(func(tx *bolt.Tx) error {
		val = tx.Bucket([]byte(b)).Get([]byte(k))
		return nil
	})

	return string(val), err
}

func Open() (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1})
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ShowsBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(SummarysBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	}); err != nil {
		panic(err) // panic because if db is not valid there is no point in continuing
	}

	return db, err
}

func Close(db *bolt.DB) error {
	return db.Close()
}

package buntdb

import (
	"github.com/qingwave/mossdb"
	"github.com/tidwall/buntdb"

	"log"
	"testing"
)

func TestBuntdb(t *testing.T) {
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("mykey", "myvalue1", nil)
		return err
	})
	if err != nil {
		log.Println(err)
	}
}
func TestMossdb(t *testing.T) {
	// create db instance
	db, err := mossdb.New(&mossdb.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// set, get, delete data
	db.Set("key1", []byte("val1"))
	log.Printf("get key1: %s", db.Get("key1"))
}

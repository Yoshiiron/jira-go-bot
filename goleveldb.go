package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

func main() {
	db, err := leveldb.OpenFile(".\\DB", nil)

	if err != nil {
		log.Fatal(err)
	}
	err = db.Put([]byte("jirakey"), []byte("11111111"), nil)
	//err = db.Put([]byte("jirakey"), []byte("11111111"), nil)

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("key: %s, value: %s\n", key, value)
	}
	iter.Release()
	err = iter.Error()

	defer db.Close()
}

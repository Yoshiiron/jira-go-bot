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
	//err = db.Put([]byte("JIRAUSER15800"), []byte("435902334"), nil)

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		db.Delete([]byte(iter.Key()), nil)
		fmt.Printf("key: %s, value: %s\n", key, value)
	}
	iter.Release()
	err = iter.Error()
	//value, err := db.Get([]byte("JIRAUSER15800"), nil)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//db.Delete([]byte("JIRAUSER15800"), nil)

	//fmt.Println(string(value))
	defer db.Close()
}

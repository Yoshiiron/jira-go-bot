package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	// Открываем базу данных
	db, err := leveldb.OpenFile(".\\DB", nil)
	if err != nil {
		fmt.Println("Ошибка открытия базы данных:", err)
		return
	}
	defer db.Close()

	// Создаем итератор для обхода ключей и значений
	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	// Перебираем ключи и соответствующие им значения
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("Ключ: %s, Значение: %s\n", key, value)
	}

	if err := iter.Error(); err != nil {
		fmt.Println("Ошибка итерации:", err)
		return
	}
}

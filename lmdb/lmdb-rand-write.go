package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/szferi/gomdb"
)

func main() {
	path := "/tmp/foo"
	err := os.Mkdir(path, 0700)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	// open the db
	env, _ := mdb.NewEnv()
	env.Open(path, 0, 0664)
	defer env.Close()
	txn, _ := env.BeginTxn(nil, 0)
	dbi, _ := txn.DBIOpen(nil, 0)

	rand.Seed(time.Now().UnixNano())

	// write some data
	num_entries := 5
	for i := 0; i < num_entries; i++ {
		key := fmt.Sprintf("Key-%d", i)
		val := fmt.Sprintf("Val-%d", rand.Intn(1000))
		txn.Put(dbi, []byte(key), []byte(val), 0)
	}
	txn.Commit()
	env.DBIClose(dbi)
}

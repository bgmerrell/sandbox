package main

import (
	"fmt"

	"github.com/szferi/gomdb"
)

func main() {
	path := "/tmp/foo"

	// open the db
	env, _ := mdb.NewEnv()
	env.SetMapSize(1 << 27)
	env.Open(path, 0, 0664)
	fmt.Println(env.Info())
	defer env.Close()
	txn, err := env.BeginTxn(nil, 0)
	if err != nil {
		fmt.Println("BeginTxn error:", err.Error())
	}
	dbi, _ := txn.DBIOpen(nil, 0)
	defer env.DBIClose(dbi)
	txn.Commit()

	bval, _ := txn.Get(dbi, []byte("Key-3"))
	fmt.Println("Val:", string(bval))
}

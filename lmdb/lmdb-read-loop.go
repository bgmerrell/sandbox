package main

import (
	"fmt"
	"time"

	"github.com/szferi/gomdb"
)

func main() {
	path := "/tmp/foo"

	// open the db
	env, err := mdb.NewEnv()
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	env.Open(path, 0, 0664)
	defer env.Close()

	txn, err := env.BeginTxn(nil, mdb.RDONLY)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	dbi, err := txn.DBIOpen(nil, 0)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	txn.Reset()

	i := 0
	n_entries := 5
	printEnabled := true
	t := time.Time{}
	for {
		// txn, err = env.BeginTxn(nil, mdb.RDONLY)
		txn.Renew()
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
		for j := 0; j < n_entries; j++ {
			bkey := fmt.Sprintf("Key-%d", j)
			bval, err := txn.Get(dbi, []byte(bkey))
			if err == mdb.NotFound {
				fmt.Println("Not found: ", bkey)
				break
			}
			if err != nil {
				panic(err)
			}
			if printEnabled {
				fmt.Printf("%s: %s\n", bkey, bval)
			}
		}
		txn.Reset()
		if printEnabled {
			if !t.IsZero() {
				fmt.Printf("%v\n", time.Since(t))
			}
			t = time.Now()
			fmt.Println("--------------------")
		}
		printEnabled = false
		i += 1
		if i == 0 || i == 1000000 {
			printEnabled = true
			i = 0
		}
	}
}

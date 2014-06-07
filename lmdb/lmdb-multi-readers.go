package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/szferi/gomdb"
)

const nTxns = 126
const nEntries = 5
const nMaxDBs = 9

var wg sync.WaitGroup

func read(txn *mdb.Txn, dbi *mdb.DBI, i int) {
	defer wg.Done()
	bkey := fmt.Sprintf("Key-%d", i)
	bval, err := txn.Get(*dbi, []byte(bkey))
	if err == nil {
		fmt.Printf("%s: %s\n", bkey, bval)
	} else {
		panic(err)
	}
}

func main() {
	path := "named-dbs"
	rand.Seed(time.Now().UnixNano())

	// open the db
	env, err := mdb.NewEnv()
	env.SetMaxDBs(nMaxDBs)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	env.Open(path, mdb.NOTLS, 0664)
	defer env.Close()

	txns := [nTxns]*mdb.Txn{}
	for i := 0; i < nTxns; i++ {
		txn, err := env.BeginTxn(nil, mdb.RDONLY)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		txns[i] = txn
	}
	dbName0 := "ones"
	_, err = txns[0].DBIOpen(&dbName0, 0)
	dbName1 := "tens"
	_, err = txns[0].DBIOpen(&dbName1, 0)
	dbName2 := "hundreds"
	dbi, err := txns[0].DBIOpen(&dbName2, 0)

	for i := 0; i < nTxns; i++ {
		txns[i].Commit()
		txns[i], err = env.BeginTxn(nil, mdb.RDONLY)
	}

	if err != nil {
		fmt.Println("Error opening DBI:", err.Error())
	}

	for i := 0; i < nTxns; i++ {
		wg.Add(1)
		go read(txns[i], &dbi, rand.Intn(nEntries)+1)
	}
	wg.Wait()
}

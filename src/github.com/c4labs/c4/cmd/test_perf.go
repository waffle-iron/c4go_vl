package main

import (
	"log"
	"fmt"
	"flag"
    "github.com/boltdb/bolt"
    "time"
)

var (
		myBucket = []byte("perftest")
		dbLocation string
		totalCnt = 0	
	)

func init (){
		flag.StringVar(&dbLocation, "db", "mybolt.db", "Location of your boltdb file")
}

func getKey(id int) string {
	user := "6a204bd89f3c8348afd5c77c717a097a"
	typeOf := "details"
	value := "2413fb3709b05939f04cf2e92f7d0897fc2596f9ad0b8a9ea855c7bfebaae892"
	return fmt.Sprintf("%s:%s:%s:%d", user, typeOf, value, id) // makes a key of hefty length
}

func handleErr(err error) {
	if err != nil {
		log.Fatalf("Unable to proceed [%s]", err)
	}
}

func main() {
	flag.Parse()

	log.Printf("Starting with dbpath [%s]", dbLocation)
	startInsert := time.Now()

	// insert batches of 10,000 records
	for i := 0; i < 10; i++ {
		insert(i)
	}
	elapsedInsert := time.Since(startInsert)

	startRead := time.Now()
	read()
	elapsedRead := time.Since(startRead)

	log.Printf("TOTAL INSERT took %s for %d items", elapsedInsert, totalCnt)
	log.Printf("TOTAL READ took %s", elapsedRead)
}

func read() {
	start := time.Now()
	db, err := bolt.Open(dbLocation, 0644, nil)
	handleErr(err)
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(myBucket))

		i := 0
		b.ForEach(func(k, v []byte) error {
			if i % 10000 == 0 {
				elapsed := time.Since(start)
				log.Printf("Read [%d] items took %s", i, elapsed)
			}
			i++
			return nil
		})
		return nil
	})
}

func insert(offset int) {
	start := time.Now()
	db, err := bolt.Open(dbLocation, 0644, nil)
	handleErr(err)
	defer db.Close()

	value := []byte(`{"exp":"2016-01-01"}`)

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(myBucket)
		handleErr(err)

		for i := 0; i < 10000; i++ {
			key := getKey(totalCnt)
			err = bucket.Put([]byte(key), value)
			handleErr(err)
			totalCnt++
		}
		return nil
	})
	handleErr(err)

	elapsed := time.Since(start)
	log.Printf("Inserted [%d] items took %s", totalCnt, elapsed)

}
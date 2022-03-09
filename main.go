package main

import (
	"context"
	"math/rand"
	"sort"
	"time"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	pd "github.com/tikv/pd/client"
)

var letters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	client, err := pd.NewClient([]string{"127.0.0.1:2379"}, pd.SecurityOption{})
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	prefix := "usertable:user280"

	var splitKeys [][]byte
	for i := 0; i < 100; i++ {
		key := prefix + randSeq(10)
		splitKeys = append(splitKeys, []byte(key))
	}

	sort.Slice(splitKeys, func(i, j int) bool {
		return string(splitKeys[i]) < string(splitKeys[j])
	})

	resp, err := client.SplitRegions(ctx, splitKeys)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}

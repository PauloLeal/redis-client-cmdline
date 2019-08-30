package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"os"
)

func NewClient(address string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}

func main() {
	var address = flag.String("host", "", "hostname:port")
	var getKey = flag.String("get", "", "get key")
	var deleteKey = flag.String("del", "", "delete key")
	var scan = flag.Bool("scan", false, "get all key")
	flag.Parse()

	if *address == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	client := NewClient(*address)

	if *scan {
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			keys, cursor, err = client.Scan(cursor, "20*", 0).Result()
			if err != nil {
				panic(err)
			}
			n += len(keys)
			for c := 0; c < len(keys); c++ {
				fmt.Println(keys[c])
			}
			if cursor == 0 {
				break
			}
		}

		fmt.Printf("found %d keys\n", n)
	} else if *getKey != "" {
		val, err := client.Get(*getKey).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(val)
	} else if *deleteKey != "" {
		_, err := client.Del(*deleteKey).Result()
		if err != nil {
			panic(err)
		}
	}
}

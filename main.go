package main

import (
	storage "RuCache/Storage"
	"fmt"
	"time"
)

func main() {
	store := storage.NewStore()
	go func() {

	}()
	ttl := 1 * time.Millisecond
	stats, err := store.Set("test1", "value1", &ttl)
	if !stats {
		fmt.Println(err)
	}
	fmt.Println(store.Get("test1"))
	stats, err = store.Del("test1")
	if !stats {
		fmt.Println(err)
	}
	fmt.Println(store.Get("test1"))
}

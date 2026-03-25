package main

import (
	storage "RuCache/Storage"
	"fmt"
)

func main() {
	store := storage.NewStore()
	fmt.Print(store)
}

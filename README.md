# go-memorycache
[![en](https://img.shields.io/badge/lang-en-red.svg)](README.md)
[![ru](https://img.shields.io/badge/lang-ru-green.svg)](README.ru.md)

In-memory cache with expiration and eviction.

## Description
A thread-safe in-memory cache implementation.

Memory is an expensive resource, so an implementation of clearing obsolete records is provided, as well as an implementation of expelling records when the specified limit is reached.<br>
Set `CleanupInterval` to enable obsolete records clearing.<br>
Set `LimitEntries` to enable records expelling when count limit is reached.<br>
Obsolete records are evicted first, then record are expelled by durability and FIFO method.

Special thanks to ks-troyan.

## Install
```
go get github.com/a-projects/go-memorycache@latest
```

## Usage
```golang

import (
	"fmt"
	"time"

	"github.com/a-projects/go-memorycache"
)

func main() {
	// create cache instance
	cache := memorycache.New(memorycache.MemoryCacheOptions{
		// periodic records clearing, every 15 min
		CleanupInterval: time.Minute * 15,
		// cache entries limit, recods
		LimitEntries: 65_536,
		// file name used to restore data from disc when app is restarted
		StoreFile: "cache.bin",
	})

	// retrieve an item from cache by key "foo"
	res, ok := cache.Get("foo")

	// if item with given key are not found
	// reasons:
	//   - entry was never stored in cache
	//   - entry was stored, but expired
	//   - entry was stored, but was removed from cache
	//   - entry was stored, but was expelled
	if !ok {
		// retrieving item from external sources
		res = "bar"

		// add entry to cache as a key value pair
		cache.Set("foo", res, memorycache.MemoryCacheEntryOptions{
			// set expiration timeout
			Expiration: time.Now().Add(time.Minute * 5),
			// set expellence resistence
			Durability: memorycache.Normal,
		})
	}

	// cast result and print to console
	fmt.Printf(res.(string))

	// close cache instance, all chache data will be writen to StoreFile file
	cache.Close()
}
```

package cache

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"

	"github.com/BurntSushi/toml"
)

//Cache is a struct for caching purposes
type Cache struct {
	sync.RWMutex
	items map[int64]*Item
}

//Item is a block with cached data
type Item struct {
	Transactions int
	Total        float64
}

//NewCache returns cache instance
func NewCache() *Cache {
	items := make(map[int64]*Item)

	cache := Cache{
		items: items,
	}

	return &cache
}

//Set sets value for block provided
func (c *Cache) Set(key int64, transactions int, total float64) {
	c.Lock()
	c.items[key] = &Item{
		Transactions: transactions,
		Total:        total,
	}
	c.Unlock()
}

//Get gets value for key provided
func (c *Cache) Get(key int64) (*Item, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	return item, true
}

//Delete deletes value for key provided
func (c *Cache) Delete(key int64) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}
	delete(c.items, key)

	return nil
}

//ReadFile reads cache from file if possible
func (c *Cache) ReadFile() {
	dat, err := ioutil.ReadFile("./internal/configs/cache.toml")
	if err != nil {
		return
	}

	newItems := make(map[string]*Item)

	err = toml.Unmarshal(dat, &newItems) //decode file
	if err != nil {
		fmt.Println(err)
		return
	}
	
	for block, item := range newItems {
		blockNum, err := strconv.ParseInt(block, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		c.items[blockNum] = item
	}
}

//WriteFile writes to file cache received
func (c *Cache) WriteFile() {
	c.Lock()
	defer c.Unlock()

	b := &bytes.Buffer{}

	newItems := make(map[string]*Item) //toml cannot decode map with int64 key
	for block, item := range c.items {
		newItems[strconv.FormatInt(block, 10)] = item
	}
	if err := toml.NewEncoder(b).Encode(&newItems); err != nil { //try to rewrite config toml file
		log.Print("Cannot encode cache to file")
	}

	ioutil.WriteFile(`./internal/configs/cache.toml`, b.Bytes(), 0600)
}

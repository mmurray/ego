package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

var mc = memcache.New("localhost:11211")

type CacheClient interface {
}

func Init() {
}

func Set(key string, obj []byte) {
	
	// err := binary.Write(buf, binary.BigEndian, obj)
	// if err != nil {
	// 	fmt.Println("binary.Write failed:", err)
	// }
	mc.Set(&memcache.Item{Key: key, Value: obj})
}

func Get(key string) (obj *memcache.Item, err error) {
	return mc.Get(key)
}

func Delete(key string) error {
	return mc.Delete(key)
}
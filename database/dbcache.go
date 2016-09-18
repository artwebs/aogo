package database

type DBCache interface {
	Conn() error
	AddCache(table, key, value string) error
	GetCache(key string) (string, error)
	IsExist(key string) bool
	DelCache(table string) error
	Close() error
}

var dbcaches = make(map[string]DBCache)
var dbcachecstr string

func RegisterDBCache(name string, d DBCache) {
	dbcaches[name] = d
}

func OpenDBCache(name, cstr string) DBCache {
	dbcachecstr = cstr
	// var val DBCache
	// switch name {
	// case "memcache":
	// 	val = &Memcache{Cstr: cstr}
	// 	break
	// case "redis":
	// 	val = &RedisCache{Cstr: cstr}
	// 	break
	// default:
	// }

	if drv, ok := dbcaches[name]; ok {
		return drv
	}
	return nil
}

func CloseDBCache(dc DBCache) {
	if dc != nil {
		dc.Close()
		dc = nil
	}
}

package database

type DBCache interface {
	Conn() error
	AddCache(table, key, value string) error
	GetCache(key string) (string, error)
	IsExist(key string) bool
	DelCache(table string) error
	Close() error
}

func OpenDBCache(name, cstr string) DBCache {
	var val DBCache
	switch name {
	case "memcache":
		val = &Memcache{Cstr: cstr}
		break
	case "redis":
		val = &RedisCache{Cstr: cstr}
		break
	default:
	}
	if val != nil {
		return val
	}

	// if drv, ok := dbcaches[name]; ok {
	// 	return drv
	// }
	return nil
}

func CloseDBCache(dc DBCache) {
	if dc != nil {
		dc.Close()
		dc = nil
	}
}

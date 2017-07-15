package web

import "github.com/artwebs/aogo/cache"

var webcacheobj *cache.Cache

func InitAppCache() (*cache.Cache, error) {
	if webcacheobj == nil {
		conf, err := InitAppConfig()
		if err != nil {
			return nil, err
		}
		CobjName := conf.String("Cache::name", "")
		CobjConfig := conf.String("Cache::config", "")
		if CobjName != "" && CobjConfig != "" {
			webcacheobj, err = cache.NewCache(CobjName, CobjConfig)
			if err != nil {
				return nil, err
			}
		}
	}
	return webcacheobj, nil
}

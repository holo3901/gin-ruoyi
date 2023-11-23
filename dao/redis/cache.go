package redis

import (
	"fmt"
	"ruoyi/models"
	"strings"
)

var cursor uint64
var keyss []string

func GetKeyList() ([]string, error) {

	// 使用SCAN命令获取所有键名
	var cursor uint64 = 0
	var keys []string
	for {
		// 执行SCAN命令
		scanCmd := client.Scan(cursor, "*", 10)
		if scanCmd.Err() != nil {
			fmt.Println("Failed to retrieve keys:", scanCmd.Err())
			return nil, scanCmd.Err()
		}
		keys, cursor = scanCmd.Val()

		// 打印缓存列表（即缓存名称）
		for _, key := range keys {
			x := SplitCache(key)
			found := false
			for _, v := range keyss {
				if v == x {
					found = true
				}
			}
			if !found {
				keyss = append(keyss, x)
			}
		}

		// 如果游标为0，表示已经遍历完所有键名
		if cursor == 0 {
			break
		}
	}
	return keyss, nil
}

func GetCacheKeysByname(name string) ([]string, error) {
	keyss, err := client.Keys(name + ":*").Result()
	if err != nil {
		fmt.Println("Failed to get cache keys:", err)
		return nil, err
	}
	return keyss, nil
}

func GetCacheValue(name string, key string) (*models.SysCache, error) {
	result, err := client.Get(name + ":" + key).Result()
	if err != nil {
		return nil, err
	}
	p := new(models.SysCache)
	p.CacheName = name
	p.CacheKey = name + ":" + key
	p.CacheValue = result
	return p, nil
}

func ClearCacheName(name string) error {
	_, err := client.Del(name + ":*").Result()
	if err != nil {
		return err
	}
	return nil
}

func ClearCacheKey(name string, key string) error {
	_, err := client.Set(name+":"+key, "", 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func ClearCacheAll() error {
	_, err := client.FlushDB().Result()
	return err
}

func UnlockByUserName(name string) error {
	_, err := client.Del("sys_logininfor:" + name).Result()
	return err
}

func DeleteDict() error {
	_, err := client.Del("sys_dict:*").Result()
	return err
}

// 获取cache缓存名称
func SplitCache(data string) string {
	var sa = strings.Split(data, ":")

	return sa[0]
}

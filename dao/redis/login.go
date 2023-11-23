package redis

import (
	"strconv"
	"time"
)

func CreateLogin(userid int64, token string) error {
	_, err := client.Set(getRedisKey(keyLoginZsetPF+strconv.Itoa(int(userid))), token, 2*time.Hour).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetLogin(userid int64) (token string, err error) {
	result, err := client.Get(getRedisKey(keyLoginZsetPF + strconv.Itoa(int(userid)))).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

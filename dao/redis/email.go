package redis

import (
	"strconv"
	"time"
)

func CreateEmail(userid int64, yanzengma int, email string) error {
	_, err := client.Set(getRedisKey(KeyEmailZsetPF+strconv.Itoa(int(userid))), yanzengma, 10*time.Minute).Result()
	if err != nil {
		return err
	}
	_, err = client.Set(getRedisKey(KeyYanZengZsetPF+strconv.Itoa(int(userid))), email, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetEmail(userid int64) (yanzengma string, email string, err error) {
	yanzengma, err = client.Get(getRedisKey(KeyEmailZsetPF + strconv.Itoa(int(userid)))).Result()
	if err != nil {
		return "", "", err
	}
	email, err = client.Get(getRedisKey(KeyYanZengZsetPF + strconv.Itoa(int(userid)))).Result()
	if err != nil {
		return "", "", err
	}
	return yanzengma, email, nil
}

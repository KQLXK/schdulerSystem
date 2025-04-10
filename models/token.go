package models

import (
	"context"
	"log"
	"schedule/database"
	"sync"
)

var (
	tokenOnce sync.Once
	tokenDao  TokenDao
)

type TokenDao struct {
}

func NewTokenDao() TokenDao {
	tokenOnce.Do(func() {
		tokenDao = TokenDao{}
	})
	return tokenDao
}

func (TokenDao) Setkey(key string, val string) error {
	err := database.Rdb.Set(context.Background(), key, val, 0).Err()
	if err != nil {
		log.Println("redis set key failed, err:", err)
		return err
	}
	return nil
}

func (TokenDao) Getkey(key string) (string, error) {
	val, err := database.Rdb.Get(context.Background(), key).Result()
	if err != nil {
		log.Println("redis get key failed, err:", err)
		return "", err
	}
	return val, nil
}

func (TokenDao) DelKey(key string) error {
	err := database.Rdb.Del(context.Background(), key).Err()
	if err != nil {
		log.Println("redis Del key failed, err:", err)
		return err
	}
	log.Println("del key sucess")
	return nil
}

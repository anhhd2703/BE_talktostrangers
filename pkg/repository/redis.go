package repository

import (
	"context"
	"fmt"
	db "talktostrangers/pkg/db/model"
	"talktostrangers/pkg/log"
	"talktostrangers/proto"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
)

var TTL = 15 * time.Second

func GetListRoomRedis(redis *db.Redis, body proto.UserInfo) []string {
	keyInfo := proto.UserInfo{
		Id:       "*",
		Language: body.Language,
		Level:    body.Level,
		Status:   proto.UnPaired,
		UserId:   body.UserId,
	}
	ownKey := keyInfo.BuildKey()
	keyInfo.UserId = "*"
	ukey := keyInfo.BuildKey()
	keys := redis.Keys(ukey)
	keys = findAndRemove(ownKey, keys)
	return keys
}
func findAndRemove(key string, arr []string) []string {
	index := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] == key {
			index = i
			break
		}
	}
	return append(arr[:index], arr[index+1:]...)
}

func Matching(redis *db.Redis, body proto.UserInfo) (string, error) {
	listUser := GetListRoomRedis(redis, body)
	if len(listUser) == 0 {
		return "", fmt.Errorf("User not found")
	}
	UpdateKey(redis, listUser[0], proto.Checking)
	// do st
	UpdateKey(redis, listUser[0], proto.Paired)
	log.Infof("[listUser]", listUser[0], body)
	roomId := uuid.New().String()
	redis.Pub(listUser[0], roomId)
	return roomId, nil
}

func PutRoomRedis(ctx context.Context, redis *db.Redis, body proto.UserInfo) (<-chan *redis.Message, error) {
	mkey := body.BuildKey()
	redis.Set(mkey, "*", TTL)
	subscriber := redis.Sub(mkey)
	msg := subscriber.Channel()
	go func() {
		select {
		case <-ctx.Done():
			subscriber.Close()
		}
	}()
	return msg, nil
}

func UpdateKey(redis *db.Redis, key string, status string) error {
	oldInfo := proto.ParseUserInfo(key)
	oldInfo.Status = status
	newKey := oldInfo.BuildKey()
	err := redis.Rename(key, newKey)
	return err
}

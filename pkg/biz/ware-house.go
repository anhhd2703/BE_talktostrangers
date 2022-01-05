package biz

import (
	"fmt"
	"sync"
	"talktostrangers/proto"

	"github.com/google/uuid"
)

type UserId string

type Level string

type Language string

type Candidates map[Language]map[Level][]UserId

type WareHouse struct {
	list  Candidates
	mutex *sync.Mutex
	sub   map[string]chan string
}

func NewWareHouse() *WareHouse {
	return &WareHouse{
		list:  make(Candidates),
		mutex: new(sync.Mutex),
		sub:   make(map[string]chan string),
	}
}

func (wh *WareHouse) Register(info proto.UserInfo) string {
	defer wh.mutex.Unlock()
	wh.mutex.Lock()
	// check user early existed
	userId, _ := wh.Match(info)
	if userId != "" {
		roomId := uuid.New().String()
		wh.Pub(userId, roomId)
		return roomId
	}
	inner, ok := wh.list[Language(info.Language)]
	if ok {
		inner[Level(info.Level)] = append(inner[Level(info.Level)], UserId(info.UserId))
	} else {
		wh.list[Language(info.Language)] = make(map[Level][]UserId)
		wh.list[Language(info.Language)][Level(info.Level)] = append(wh.list[Language(info.Language)][Level(info.Level)], UserId(info.UserId))
	}
	return ""
}

func (wh *WareHouse) Sub(userId string) <-chan string {
	defer wh.mutex.Unlock()
	wh.mutex.Lock()
	wh.sub[userId] = make(chan string)
	return wh.sub[userId]
}
func (wh *WareHouse) UnSub(userId string) {
	defer wh.mutex.Unlock()
	wh.mutex.Lock()
	close(wh.sub[userId])
	delete(wh.sub, userId)
}

func (wh *WareHouse) Pub(userId, key string) error {
	val, ok := wh.sub[userId]
	if ok {
		val <- key
		return nil
	} else {
		return fmt.Errorf("User is not existed")
	}
}

func (wh *WareHouse) Match(info proto.UserInfo) (string, error) {
	language, ok := wh.list[Language(info.Language)]
	if !ok {
		return "", fmt.Errorf("User not found")
	} else {
		if level, isExit := language[Level(info.Level)]; isExit {
			if len(level) > 0 {
				return string(level[0]), nil
			}
		}
	}
	return "", fmt.Errorf("User not found")
}

func (wh *WareHouse) UnRegister(info proto.UserInfo) {
	language, ok := wh.list[Language(info.Language)]
	if ok {
		if level, isExit := language[Level(info.Level)]; isExit {
			language[Level(info.Level)] = remove(level, UserId(info.UserId))
		}
	}
}
func remove(s []UserId, r UserId) []UserId {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func write(out chan string, i string) (err error) {
	defer func() {
		// recover from panic caused by writing to a closed channel
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			fmt.Printf("write: error writing %d on channel: %v\n", i, err)
			return
		}
		fmt.Printf("write: wrote %d on channel\n", i)
	}()
	if len(out) == 1 {
		<-out
	}
	out <- i // write on possibly closed channel
	return err
}

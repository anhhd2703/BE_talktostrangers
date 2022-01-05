package proto

import (
	"strings"
)

const (
	Paired   = "paired"
	UnPaired = "unpaired"
	Checking = "checking"
)

type UserInfo struct {
	Id       string
	UserId   string `json:"user_id"`
	Language string `json:"language"`
	Level    string `json:"level"` //  low  medium  high
	Status   string
}

func (r UserInfo) BuildKey() string {
	if r.Language == "" {
		r.Language = "none"
	}
	if r.Level == "" {
		r.Level = "none"
	}
	if r.Status == "" {
		r.Status = UnPaired
	}
	strs := []string{"lg", r.Language, "lv", r.Level, "status", r.Status, "uid", r.UserId}
	return strings.Join(strs, "/")
}

func ParseUserInfo(key string) UserInfo {
	arr := strings.Split(key, "/")
	return UserInfo{
		Status:   arr[5],
		UserId:   arr[7],
		Language: arr[1],
		Level:    arr[3],
	}
}

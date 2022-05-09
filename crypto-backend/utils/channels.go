package utils

import (
	"crypto-backend/models"
	"time"
)

type MsgType int

const (
	Insert MsgType = iota
	Delete
)

type TimeUpdateMsg struct {
	Type   MsgType
	UserId models.UserId
	Time   time.Time
}

var TimeUpdateChan = make(chan TimeUpdateMsg, 10)

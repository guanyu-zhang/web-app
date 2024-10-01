package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	GroupId  uint
	SenderId uint `gorm:"index"`
	Content  string
	IsActive bool
	UpVotes  int
}

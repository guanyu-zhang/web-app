package daos

import (
	"gorm.io/gorm"
	. "web-app-backend/models"
)

func CreateGroup(db *gorm.DB) (*Group, error) {
	group := &Group{}
	if err := db.Create(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

//func GetGroupByGroupId(db *gorm.DB, groupId int) (*Group, error) {
//	var group Group
//	err := db.Where("id = ?", groupId).First(&group).Error
//	if err != nil {
//		return nil, err
//	}
//	return &group, nil
//}

//func JoinGroup(db *gorm.DB, groupId int, users []User) error {
//	group, err := GetGroupByGroupId(db, groupId)
//	if err != nil {
//		return err
//	}
//	err = db.Model(group).Update("users", users).Error
//	return err
//}

//func GetGroupsByUserId(db *gorm.DB, id int) ([]Group, error){
//	var groups []Group
//	err := db.Preload("users").Where("users = ?", ).Find(&groups).Error
//	if err != nil {
//		return nil, err
//	}
//
//}

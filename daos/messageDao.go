package daos

import (
	"gorm.io/gorm"
	. "web-app-backend/models"
)

func CreateMessage(db *gorm.DB, groupId uint, userId uint, content string, isActive bool) (*Message, error) {
	message := &Message{
		GroupId:  groupId,
		SenderId: userId,
		Content:  content,
		IsActive: isActive,
	}
	if err := db.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func GetGroupIdsByUserId(db *gorm.DB, id uint) ([]uint, error) {
	var groupIds []uint
	if err := db.Model(&Message{}).Where("sender_id = ?", id).Distinct("group_id").
		Pluck("group_id", &groupIds).Error; err != nil {
		return nil, err
	}
	return groupIds, nil
}

func GetAllReceiversByUserId(db *gorm.DB, id uint) (map[uint]map[string]interface{}, error) {
	groupIds, err := GetGroupIdsByUserId(db, id)
	if err != nil {
		return nil, err
	}
	ret := make(map[uint]map[string]interface{})
	for _, gid := range groupIds {
		arrang, err := GetGroupArrangementByGroupId(db, gid)
		if err != nil {
			return nil, err
		}
		for _, u := range arrang {
			if u != id {
				user, err := GetUserByUserId(db, u)
				if err != nil {
					return nil, err
				}
				ret[gid] = map[string]interface{}{"userId": u, "userName": user.Username}
			}
		}
	}
	return ret, nil
}

func GetGroupArrangementByGroupId(db *gorm.DB, groupId uint) ([]uint, error) {
	var userIds []uint
	if err := db.Model(&Message{}).Where("group_id = ?", groupId).Distinct("sender_id").
		Pluck("sender_id", &userIds).Error; err != nil {
		return nil, err
	}
	return userIds, nil
}

func UpVoteMessageByMessageId(db *gorm.DB, messageId uint) error {
	result := db.Model(&Message{}).Where("id = ?", messageId).Update("up_votes", gorm.Expr("up_votes + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DownVoteMessageByMessageId(db *gorm.DB, messageId uint) error {
	result := db.Model(&Message{}).Where("id = ?", messageId).Update("up_votes", gorm.Expr("up_votes - ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetMessagesInOrderByGroupId(db *gorm.DB, groupId uint) ([]Message, error) {
	var messages []Message
	err := db.Model(&Message{}).Where("group_id = ?", groupId).Order("id asc").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

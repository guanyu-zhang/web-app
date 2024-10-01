package daos

import "gorm.io/gorm"
import . "web-app-backend/utils"

func FindPrevGroupForTwo(db *gorm.DB, userId1 uint, userId2 uint) (uint, error) {
	g1, err := GetGroupIdsByUserId(db, userId1)
	if err != nil {
		return 0, err
	}
	g2, err := GetGroupIdsByUserId(db, userId2)
	if err != nil {
		return 0, err
	}

	prevG, err := FindCommonElement(g1, g2)
	if err != nil {
		return 0, err
	}
	return prevG, nil
}

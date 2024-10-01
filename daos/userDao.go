package daos

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
import . "web-app-backend/models"
import "crypto/rand"

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUserId(db *gorm.DB, userId uint) (*User, error) {
	var user User
	err := db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func CreateUser(db *gorm.DB, username string, password string) (*User, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return nil, err
	}
	toHash := append([]byte(password), salt...)
	hashedPassword, err := bcrypt.GenerateFromPassword(toHash, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		Username: username,
		Password: string(hashedPassword),
		Salt:     salt,
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func LoginHelper(db *gorm.DB, username string, password string) (*User, error) {
	user, err := GetUserByUsername(db, username)
	if err != nil {
		return nil, err
	}
	salted := append([]byte(password), user.Salt...)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), salted); err != nil {
		return nil, err
	}
	return user, nil
}

// not suitable for large db
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

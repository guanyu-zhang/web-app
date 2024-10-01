package controllers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	. "web-app-backend/daos"
	"web-app-backend/models"
	. "web-app-backend/utils"
)

var (
	identityKey = "username"
	idKey       = "id"
)

func PayloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*models.User); ok {
			return jwt.MapClaims{
				identityKey: v.Username,
			}
		}
		return jwt.MapClaims{}
	}
}

func IdentityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		tx, err := GetTx(c)
		if err != nil {
			return nil
		}
		user, err := GetUserByUsername(tx, claims[identityKey].(string))
		if err != nil {
			return nil
		}
		c.Set(idKey, user.ID)
		return &models.User{
			Username: claims[identityKey].(string),
		}
	}
}

// used for login 登陆阶段使用
func Authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals models.Login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		username := loginVals.Username
		password := loginVals.Password

		tx, err := GetTx(c)
		if err != nil {
			return nil, err
		}
		_, err = LoginHelper(tx, username, password)
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		return &models.User{
			Username: username,
		}, nil
	}
}

// 看是否有权限操作
func Authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if _, ok := data.(*models.User); ok {
			return true
		}
		return false
	}
}

func Unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func SignUpHandler(c *gin.Context) {

	var signUpVals models.Login
	if err := c.ShouldBind(&signUpVals); err != nil {
		log.Printf("jwt missing sign up values" + jwt.ErrMissingLoginValues.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "missing username or password",
		})
		return
	}
	username := signUpVals.Username
	password := signUpVals.Password

	tx, err := GetTx(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	if _, err := GetUserByUsername(tx, username); err == nil {
		log.Printf("GetUserByUsername return != nil, user exists ")
		c.JSON(http.StatusConflict, gin.H{
			"status":  "fail",
			"message": "user exists",
		})
		return
	}

	user, err := CreateUser(tx, username, password)
	if err != nil {
		log.Printf("CreateUser err " + err.Error())
	}

	c.JSON(200, gin.H{
		"status":   "success",
		"username": user.Username,
	})
}

func GetAllUsersHandler(c *gin.Context) {
	tx, err := GetTx(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	users, err := GetAllUsers(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	var ret []string
	for _, user := range users {
		ret = append(ret, user.Username)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   ret,
	})
}

func UnknownPeopleHandler(c *gin.Context) {
	tx, err := GetTx(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	users, err := GetAllUsers(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	var usernames []string
	for _, u := range users {
		usernames = append(usernames, u.Username)
	}
	userId, _ := c.Get(idKey)
	user, _ := c.Get(identityKey)
	userName := user.(*models.User).Username
	usernameslice := RemoveValue(usernames, userName)
	var unames []string
	groups, err := GetAllReceiversByUserId(tx, userId.(uint))
	for _, g := range groups {
		unames = append(unames, g["userName"].(string))
	}
	ret := SliceDiff(usernameslice, unames)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   ret,
	})
}

func IsUserAllowed(db *gorm.DB, userId uint, channelId uint) (bool, error) {
	gids, err := GetGroupIdsByUserId(db, userId)
	if err != nil {
		return false, err
	}
	for _, gid := range gids {
		if channelId == gid {
			return true, nil
		}
	}
	return false, err
}

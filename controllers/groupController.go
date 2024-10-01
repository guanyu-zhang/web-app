package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "web-app-backend/daos"
	. "web-app-backend/utils"
)

type create_chat_room_for_two struct {
	Receiver string `json:"receiver"`
}

func CreateChatRoomForTwoHandler(c *gin.Context) {
	var ccr_2 create_chat_room_for_two
	if err := c.BindJSON(&ccr_2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx, err := GetTx(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	username2 := ccr_2.Receiver
	userId1, _ := c.Get(idKey)
	user2, err := GetUserByUsername(tx, username2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "the person you search does not exist",
		})
		return
	}
	_, err = FindPrevGroupForTwo(tx, userId1.(uint), user2.ID)
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "you can not create a group that exits",
		})
		return
	}
	group, err := CreateGroup(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	if _, err := CreateMessage(tx, group.ID, userId1.(uint), "dummy init msg", false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}

	if _, err := CreateMessage(tx, group.ID, user2.ID, "dummy init msg", false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"groupId": group.ID,
	})
}

/*
*
return example
"2" is groupid

	{
	    "data": {
	        "2": {
	            "userid": 1,
	            "username": "test"
	        }
	    },
	    "status": "success"
	}
*/
func GetMyGroupsHandler(c *gin.Context) {
	tx, err := GetTx(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	userId, _ := c.Get(idKey)
	groups, err := GetAllReceiversByUserId(tx, userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   groups,
	})
}

package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	. "web-app-backend/daos"
	"web-app-backend/initializers"
	. "web-app-backend/utils"
)

type upvote_message struct {
	MessageId string `json:"messageId"`
}

func UpVoteHandler(c *gin.Context) {
	var upvote upvote_message
	if err := c.BindJSON(&upvote); err != nil {
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
	msgId, err := strconv.ParseUint(upvote.MessageId, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	err = UpVoteMessageByMessageId(tx, uint(msgId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func DownVoteHandler(c *gin.Context) {
	var downVote upvote_message
	if err := c.BindJSON(&downVote); err != nil {
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
	msgId, err := strconv.ParseUint(downVote.MessageId, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	err = DownVoteMessageByMessageId(tx, uint(msgId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

type view_current_messages struct {
	GroupId string `json:"groupId"`
}

func ViewCurrentMessagesHandler(c *gin.Context) {
	var view_current_messages view_current_messages
	if err := c.BindJSON(&view_current_messages); err != nil {
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
	gid, err := strconv.ParseUint(view_current_messages.GroupId, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	msgs, err := GetMessagesInOrderByGroupId(tx, uint(gid))
	var ret []map[string]interface{}
	for _, msg := range msgs {
		u, err := GetUserByUserId(tx, msg.SenderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "unknown error",
			})
			return
		}
		sendername := u.Username
		m := map[string]interface{}{
			"id":          msg.ID,
			"created_at":  msg.CreatedAt,
			"sender_name": sendername,
			"content":     msg.Content,
			"is_active":   msg.IsActive,
			"up_votes":    msg.UpVotes,
		}
		ret = append(ret, m)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	u, exists := c.Get(identityKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   u,
		"data":   ret,
	})
}

func RecordMessage(msg string) {
	var msgJson map[string]interface{}
	err := json.Unmarshal([]byte(msg), &msgJson)
	if err != nil {
		log.Printf("Failed to unmarshal JSON and will not save to db: %v", err)
	}
	name, ok := msgJson["senderName"].(string)
	if ok {
		log.Println("senderName:", name)
	} else {
		log.Println("Key 'senderName' not found or value is not a string")
	}
	u, err := GetUserByUsername(initializers.DB, name)
	if err != nil {
		log.Println("can't find user")
	}
	content, ok := msgJson["content"].(string)
	if ok {
		log.Println("content:", content)
	} else {
		log.Println("Key 'content' not found or value is not a string")
	}
	groupId, ok := msgJson["groupId"].(string)
	if ok {
		log.Println("groupId:", groupId)
	} else {
		log.Println("Key 'groupId' not found or value is not a string")
	}
	gid, err := strconv.ParseUint(groupId, 10, 0)
	_, err = CreateMessage(initializers.DB, uint(gid), u.ID, content, true)
}

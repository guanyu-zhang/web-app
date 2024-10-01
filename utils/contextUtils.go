package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var TxKey = "txKey"
var ErrValueNotFound = errors.New("value not found in context")

func GetTx(c *gin.Context) (*gorm.DB, error) {
	tx, exists := c.Get(TxKey)
	if !exists {
		return nil, ErrValueNotFound
	}
	return tx.(*gorm.DB), nil
}

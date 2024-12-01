package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReqBody(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	return nil
}

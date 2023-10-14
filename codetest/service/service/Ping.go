package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

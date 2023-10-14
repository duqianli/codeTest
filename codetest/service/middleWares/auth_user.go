package middleWares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service/helper"
)

//AuthUserCheck
//check is user
func AuthUserCheck() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			context.Abort()
			context.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized",
			})

			return
		}
		if userClaim == nil {
			context.Abort()
			context.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized user",
			})
			return
		}
		context.Set("userClaim", userClaim)
		context.Next()
	}
}

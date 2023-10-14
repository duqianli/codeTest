package middleWares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service/helper"
)

//AuthAdminCheck
//check is admin
func AuthAdminCheck() gin.HandlerFunc {
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
		if userClaim.IsAdmin != 1 {
			context.Abort()
			context.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "no admin",
			})
			return
		}
		context.Next()
	}
}

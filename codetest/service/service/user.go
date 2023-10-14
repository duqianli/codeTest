package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"service/define"
	helper "service/helper"
	"service/models"
	"strconv"
	"time"
)

var ctx = context.Background()

// GetUserDetail
// @Summary 用户详情
// @Schemes
// @Description do ping
// @Tags 公共方法
// @Accept json
// @Produce json
// @param identity query string false "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user_detail [get]
func GetUserDetail(context *gin.Context) {
	identity := context.Query("identity")
	if identity == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户的唯一标识不能为空",
		})
		return
	}
	data := new(models.UserBasic)
	err := models.DB.Omit("password").Where("identity = ?", identity).First(&data).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Detail By Identity :" + identity + "Error:" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  data,
	})
	return
}

// Login
// @Summary 用户登录
// @Schemes
// @Description do ping
// @Tags 公共方法
// @param username formData string false "username"
// @param password formData string false "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /login [post]
func Login(context *gin.Context) {
	username := context.PostForm("username")
	password := context.PostForm("password")
	if username == "" || password == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}
	password = helper.GetMd5(password)
	fmt.Println(password)
	data := new(models.UserBasic)
	err := models.DB.Where("name = ? AND password = ?", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "密码错误",
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Basic Error:" + err.Error(),
		})
		return
	}
	token, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Generate Token Error:" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// SendCode
// @Summary 发送邮件
// @Schemes
// @Description do ping
// @Tags 公共方法
// @param email formData string true "email"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /sendCode [post]
func SendCode(context *gin.Context) {
	email := context.PostForm("email")
	if email == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	code := helper.GetRand()
	err := helper.SendCode(email, code)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "send code error:" + err.Error(),
		})
		return
	}
	models.RDB.Set(ctx, email, code, time.Second*300)
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "验证码发送成功",
	})
}

// Register
// @Summary 注册
// @Schemes
// @Description do ping
// @Tags 公共方法
// @param code formData string true "code"
// @param name formData string true "name"
// @param password formData string true "password"
// @param phone formData string false "phone"
// @param mail formData string true "mail"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /register [post]
func Register(context *gin.Context) {
	userCode := context.PostForm("code")
	name := context.PostForm("name")
	password := context.PostForm("password")
	phone := context.PostForm("phone")
	mail := context.PostForm("mail")
	if userCode == "" || name == "" || password == "" || mail == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msd":  "参数不正确",
		})
		return
	}

	//验证码是否正确
	sysCode, err := models.RDB.Get(ctx, mail).Result()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Code Error:" + err.Error(),
		})
		return
	}
	if sysCode != userCode {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码错误",
		})
		return
	}
	var cnt int64
	err = models.DB.Model(new(models.UserBasic)).Where("mail = ?", mail).Count(&cnt).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Error:" + err.Error(),
		})
		return
	}
	if cnt > 0 {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已被注册",
		})
		return
	}
	//插入数据
	data := &models.UserBasic{
		Model:    gorm.Model{},
		Identity: helper.GetUUID(),
		Name:     name,
		Password: helper.GetMd5(password),
		Phone:    phone,
		Mail:     mail,
	}
	err = models.DB.Create(&data).Error
	if err != nil {
		log.Printf("Get Code Error:%v\n", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "请重新获取验证码",
		})
		return
	}
	token, err := helper.GenerateToken(data.Identity, name, data.IsAdmin)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Generate Token Error:" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// GetRankList
// @Summary 排行榜
// @Schemes
// @Description do ping
// @Tags 公共方法
// @param page query int false "page"
// @param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /rank_list [get]
func GetRankList(context *gin.Context) {
	size, _ := strconv.Atoi(context.DefaultQuery("size", define.DafaultSize))
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Getproblemlist page strconv error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	list := make([]*models.UserBasic, 0)
	err = models.DB.Model(new(models.UserBasic)).Count(&count).Order("pass_num DESC,submit_num ASC").
		Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Rank List Error:" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

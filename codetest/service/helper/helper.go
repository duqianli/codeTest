package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"isAdmin"`
	jwt.RegisteredClaims
}

var key = []byte("gon-gorm-oj-key")

//GetMd5
//生成md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

//GenerateToken
//生成token
func GenerateToken(identity, name string, isAdmin int) (string, error) {

	userClaim := &UserClaims{
		Identity:         identity,
		Name:             name,
		IsAdmin:          isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	token, err := claim.SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

//AnalyseToken
//解析token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	claim := new(UserClaims)
	claimsss, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !claimsss.Valid {
		fmt.Printf("analyse token err:", err)
		return nil, err
	}
	return claim, nil
}

//SendCode
//发送邮件
func SendCode(toEmail, code string) error {
	e := email.NewEmail()
	e.From = "Jordan Wright <2577595258@qq.com>"
	e.To = []string{toEmail}
	e.Subject = "验证码已发送，请查收。"
	e.HTML = []byte("你的验证码是：<b>" + code + "</b>")
	return e.SendWithTLS("smtp.qq.com:587",
		smtp.PlainAuth("", "2577595258@qq.com", "gczulbflmzvddidi", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})

}

//GetUUID
//获取uuid
func GetUUID() string {
	return uuid.NewV4().String()
}

//GetRand
//获取验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

//CodeSave
//保存代码
func CodeSave(code []byte) (string, error) {
	//代码文件夹路径
	dirName := "/code" + GetUUID()
	//代码路径
	path := dirName + "/main.go"
	//创建路径
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	//创建文件
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	//写入数据
	f.Write(code)
	//关闭流
	defer f.Close()
	return path, nil
}

package service

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"service/define"
	"service/helper"
	"service/models"
	"strconv"
	"sync"
	"time"
)

// GetSubmitList
// @Summary 问题列表
// @Schemes
// @Description do ping
// @Tags 公共方法
// @Accept json
// @Produce json
// @param page query int false "page"
// @param size query int false "size"
// @param problem_identity query int false "problem_identity"
// @param user_identity query string false "user_identity"
// @param status query string false "status"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit_list [get]
func GetSubmitList(context *gin.Context) {
	size, _ := strconv.Atoi(context.DefaultQuery("size", define.DafaultSize))
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Getproblemlist page strconv error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	list := make([]models.SubmitBasic, 0)
	problemIdentity := context.Query("problem_identity")
	userIdentity := context.Query("user_identity")
	status, _ := strconv.Atoi(context.Query("status"))
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get submit list Error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get submit list Error:" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": -1,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
	return
}

// Submit
// @Summary 代码提交
// @Schemes
// @Description do ping
// @Tags 用户私有方法
// @Accept json
// @Produce json
// @param problem_identity query int false "problem_identity"
// @param code body string false "body"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/submit [post]
func Submit(context *gin.Context) {
	problemIdentity := context.Query("problem_identity")
	code, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "read code error:" + err.Error(),
		})
		return
	}
	//保存代码
	path, err := helper.CodeSave(code)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "read save error:" + err.Error(),
		})
		return
	}
	u, _ := context.Get("userClaim")
	userClaim := u.(helper.UserClaims)
	sb := &models.SubmitBasic{
		Identity:        helper.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaim.Identity,
		Path:            path,
	}
	//代码判断
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCase").First(pb).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get problem error:" + err.Error(),
		})
		return
	}
	//答案错误的channal
	WA := make(chan int)
	//超内存的channal
	OOM := make(chan int)
	//编译错误的channal
	CE := make(chan int)
	//通过的个数
	passCount := 0
	//提示信息
	var msg string
	var lock sync.Mutex
	for _, testCase := range pb.TestCases {
		go func() {
			//执行测试
			//创建一个命令
			cmd := exec.Command("go", "run", "runnerCode/runCode.go")
			//配置相关数据
			var stderr, out bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out
			stdinpip, err := cmd.StdinPipe()
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(stdinpip, testCase.Input)
			//运行
			//记录初始内存
			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)
			if err := cmd.Run(); err != nil {
				log.Println(err, stderr.String())
				if err.Error() == "exit status 2" {
					msg = stderr.String()
					CE <- 1
					return
				}
			}
			//记录运行后的内存状态
			var em runtime.MemStats
			runtime.ReadMemStats(&em)
			//答案错误
			if testCase.Output != out.String() {
				msg = "答案错误"
				WA <- 1
				return
			}
			//运行超内存
			if em.Alloc/1024-(bm.Alloc/1024) > uint64(pb.MaxRuntime) {
				msg = "运行超内存"
				OOM <- 1
				return
			}
			lock.Lock()
			passCount++
			lock.Unlock()
		}()
	}
	// 【-1-待判断，1-答案正确，2-答案错误，3-运行超时，4-运行超内存， 5-编译错误，6-非法代码】
	select {
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-CE:
		sb.Status = 5

	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		if passCount == len(pb.TestCases) {
			sb.Status = 1
		} else {
			sb.Status = 3
		}

	}
	//更新相关记录
	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(sb).Error
		if err != nil {
			return errors.New("SubmitBasic Save Error:" + err.Error())
		}
		m := make(map[string]interface{})
		m["submit_num"] = gorm.Expr("submit_num + ?", 1)
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		//更新user_basic
		err = tx.Model(new(models.UserBasic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
		if err != nil {
			return errors.New("UserBasic Modify Error:" + err.Error())
		}
		//更新problem_basic
		err = tx.Model(new(models.ProblemBasic)).Where("identity = ?", problemIdentity).Updates(m).Error
		if err != nil {
			return errors.New("ProblemBasic Modify Error:" + err.Error())
		}
		return nil
	}); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "submit create error:" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"status": sb.Status,
			"msg":    msg,
		},
	})
}

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"service/define"
	"service/helper"
	"service/models"
	"strconv"
)

// GetProblemList
// @Summary 问题列表
// @Schemes
// @Description do ping
// @Tags 公共方法
// @Accept json
// @Produce json
// @param page query int false "page"
// @param size query int false "size"
// @param keyword query int false "keyword"
// @param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem_list [get]
func GetProblemList(context *gin.Context) {
	size, _ := strconv.Atoi(context.DefaultQuery("size", define.DafaultSize))
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	categoryIdentity := context.Query("category_identity")
	keyword := context.Query("keyword")
	if err != nil {
		log.Println("Getproblemlist page strconv error:", err)
	}
	page = (page - 1) * size
	var count int64
	tx := models.GetProblemList(keyword, categoryIdentity)
	list := make([]*models.ProblemBasic, 0)
	err = tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get problem list error:", err)
	}
	context.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// GetProblemDetail
// @Summary 问题详情
// @Schemes
// @Description do ping
// @Tags 公共方法
// @Accept json
// @Produce json
// @param identity query string false "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem_detail [get]
func GetProblemDetail(context *gin.Context) {
	identity := context.Query("identity")
	if identity == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题的唯一标识不能为空",
		})
		return
	}
	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题不存在",
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GET PROBLEM ERROR :" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// PloblemCreate
// @Summary 问题创建
// @Schemes
// @Description do ping
// @Tags 管理员私有方法
// @param authorization header string true "authorization"
// @param title formData string true "title"
// @param content formData string true "content"
// @param max_runtime formData string false "max_runtime"
// @param max_mem formData string false "max_mem"
// @param category_ids formData array false "category_ids"
// @param test_case formData array true "test_case"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
func PloblemCreate(context *gin.Context) {
	title := context.PostForm("title")
	content := context.PostForm("content")
	max_runtime, _ := strconv.Atoi(context.PostForm("max_runtime"))
	max_mem, _ := strconv.Atoi(context.PostForm("max_mem"))
	category_ids := context.PostFormArray("category_ids")
	test_case := context.PostFormArray("test_case")
	if title == "" || content == "" || len(category_ids) == 0 || len(test_case) == 0 {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}

	identity := helper.GetUUID()
	data := &models.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: max_runtime,
		MaxMem:     max_mem,
	}
	//处理分类
	categoryBasics := make([]*models.ProblemCategory, 0)
	for _, id := range category_ids {
		categoryId, _ := strconv.Atoi(id)
		categoryBasics = append(categoryBasics, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(categoryId),
		})
	}
	data.ProblemCategories = categoryBasics

	//处理测试用例
	//{"input":"1 2\n","output":"3\n"}
	testCaseBasics := make([]*models.TestCase, 0)
	for _, testCase := range test_case {
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式化错误",
			})
		}
		if _, ok := caseMap["input"]; !ok {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式化错误",
			})
		}
		if _, ok := caseMap["output"]; !ok {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式化错误",
			})
		}
		testCaseBasic := &models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		}

		testCaseBasics = append(testCaseBasics, testCaseBasic)
	}
	data.TestCases = testCaseBasics
	fmt.Println("测试测试测试")
	//创建
	err := models.DB.Create(data).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Create Error:" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": map[string]interface{}{
			"identity": identity,
		},
	})
}

// ProblemModify
// @Summary 问题修改
// @Schemes
// @Description do ping
// @Tags 管理员私有方法
// @param authorization header string true "authorization"
// @param identity formData string false "identity"
// @param title formData string true "title"
// @param content formData string true "content"
// @param max_runtime formData string false "max_runtime"
// @param max_mem formData string false "max_mem"
// @param category_ids formData []string false "category_ids" collectionFormat[multi]
// @param test_case formData []string true "test_case" collectionFormat[multi]
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-modify [put]
func ProblemModify(context *gin.Context) {
	identity := context.PostForm("identity")
	title := context.PostForm("title")
	content := context.PostForm("content")
	max_runtime, _ := strconv.Atoi(context.PostForm("max_runtime"))
	max_mem, _ := strconv.Atoi(context.PostForm("max_mem"))
	category_ids := context.PostFormArray("category_ids")
	test_case := context.PostFormArray("test_case")
	if identity == "" || title == "" || content == "" || len(category_ids) == 0 || len(test_case) == 0 || max_runtime == 0 || max_mem == 0 {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		//问题基础信息的保存
		problemBasic := &models.ProblemBasic{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxRuntime: max_runtime,
			MaxMem:     max_mem,
		}
		err := tx.Where("identity = ?", identity).Updates(problemBasic).Error
		if err != nil {
			return err
		}
		//查询问题详情
		err = tx.Where("identity = ?", identity).Find(problemBasic).Error
		if err != nil {
			return err
		}
		//关联分类信息的更新
		//1.删除已存在的分类关系
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil {
			return err
		}
		//2.更新新的分类关系
		pcs := make([]*models.ProblemCategory, 0)
		for _, id := range category_ids {
			intId, _ := strconv.Atoi(id)
			pcs = append(pcs, &models.ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(intId),
			})
		}
		err = tx.Create(&pcs).Error
		if err != nil {
			return err
		}
		//关联测试案例的更新
		//1.删除旧的关联关系
		err = tx.Where("problem_identity = ?", problemBasic.Identity).Delete(new(models.TestCase)).Error
		if err != nil {
			return err
		}
		//2。创建新的关联关系
		tcs := make([]*models.TestCase, 0)
		for _, testCase := range test_case {
			caseMap := make(map[string]string)
			err = json.Unmarshal([]byte(testCase), &caseMap)
			if err != nil {
				return err
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("测试案例格式错误")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("测试案例格式错误")
			}
			tcs = append(tcs, &models.TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}

		return nil
	}); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Modify Error:" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题修改成功",
	})

}

package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"service/define"
	"service/helper"
	"service/models"
	"strconv"
)

// CategoryList
// @Summary 分类列表
// @Schemes
// @Description do ping
// @Tags 管理员私有方法
// @Accept json
// @Produce json
// @param authorization header string true "authorization"
// @param page query int false "page"
// @param size query int false "size"
// @param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/Category_list [get]
func CategoryList(context *gin.Context) {
	size, _ := strconv.Atoi(context.DefaultQuery("size", define.DafaultSize))
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	keyword := context.Query("keyword")
	if err != nil {
		log.Println("Getproblemlist page strconv error:", err)
	}
	page = (page - 1) * size
	var count int64
	categoryList := make([]*models.CategoryBasic, 0)
	err = models.DB.Model(new(models.CategoryBasic)).Where("name like ?", "%"+keyword+"%").
		Count(&count).Offset(page).Limit(size).Find(&categoryList).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get category error:" + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": map[string]interface{}{
			"list":  categoryList,
			"count": count,
		},
	})
}

// CategoryCreate
// @Summary 分类创建
// @Schemes
// @Description do ping
// @Tags 管理员私有方法
// @param authorization header string true "authorization"
// @param name formData string false "name"
// @param parentId formData int false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/Category_create [post]
func CategoryCreate(context *gin.Context) {
	name := context.PostForm("name")
	parentId, _ := strconv.Atoi(context.PostForm("parentId"))
	category := &models.CategoryBasic{
		Identity: helper.GetUUID(),
		Name:     name,
		ParentId: parentId,
	}
	err := models.DB.Create(category).Error
	if err != nil {
		log.Println("Category create error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建失败",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "创建成功",
	})
}

//CategoryModify
// @Summary 分类修改
// @Schemes
// @Description do ping
// @Tags 管理员私有方法
// @param authorization header string true "authorization"
// @param identity formData string false "identity"
// @param name formData string false "name"
// @param parentId formData int false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/Category_modify [put]
func CategoryModify(context *gin.Context) {
	name := context.PostForm("name")
	parentId, _ := strconv.Atoi(context.PostForm("parentId"))
	identity := context.PostForm("identity")
	if name == "" || identity == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	category := &models.CategoryBasic{
		Identity: identity,
		Name:     name,
		ParentId: parentId,
	}
	err := models.DB.Model(new(models.CategoryBasic)).Where("identity = ?", identity).Updates(category).Error
	if err != nil {
		log.Println("Category modify error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "修改分类失败",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "修改成功",
	})

}

//CategoryDelete
// @Summary 删除分类
// @Schemes
// @Description do ping
// @Tags 管理员私有方法
// @param authorization header string true "authorization"
// @param identity query string false "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/Category_delete [delete]
func CategoryDelete(context *gin.Context) {
	identity := context.Query("identity")
	if identity == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	var cnt int64
	err := models.DB.Model(new(models.ProblemBasic)).Where("identity =(select id from category_basic where identity = ?)", identity).
		Count(&cnt).Error
	if err != nil {
		log.Println("Get problem_basic error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类关联的问题失败",
		})
		return
	}
	if cnt > 0 {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类下已存在问题",
		})
		return
	}
	err = models.DB.Where("identity = ?", identity).Delete(&models.CategoryBasic{}).Error
	if err != nil {
		log.Println("delete problem_basic error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除分类失败失败",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "删除成功",
	})
}

package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemId     uint           `gorm:"column:problem_id;type:int(11);"json:"problemId"'`  //问题的ID
	CategoryId    uint           `gorm:"column:category_id;type:int(11);"json:"categoryId"` //分类的ID
	CategoryBasic *CategoryBasic `gorm:"foreignKey:id;references:CategoryId"`               //关联分类基础信息表
}

func (table ProblemCategory) TableName() string {
	return "problem_category"
}

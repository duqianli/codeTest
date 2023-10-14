package test

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"service/models"
	"testing"
)

func TestGormTest(t *testing.T) {
	dsn := "root:root@tcp(127.0.0.1:3306)/gin_gron_oj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("openerr:", err)
	}

	data := make([]*models.ProblemBasic, 0)
	dg := db.Find(&data)
	fmt.Println(dg.RowsAffected)
	for _, v := range data {
		fmt.Println("problem===>", v)
	}
}

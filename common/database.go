package common

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"graduationEnd/model"
)

var DB = &gorm.DB{}

func init() {
	db, err := gorm.Open("mysql", "root:123456@/graduation?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{})

	db.AutoMigrate(&model.Board{})

	db.AutoMigrate(&model.Course{})

	db.AutoMigrate(&model.Publication{})
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}

package service

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetMySqlDb(host string, port int, user string, password string, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

package service

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetPgDb(host string, port int, user string, password string, dbName string) (*gorm.DB, error) {
	// if len(timezone) == 0{
	// 	timezone = []string{
	// 		"Asia/Yangon",
	// 	}
	// }
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, url.PathEscape(user), url.PathEscape(password), url.PathEscape(dbName))
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

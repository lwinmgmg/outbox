package service

import "gorm.io/gorm"

var (
	DbMap = map[string]func(string, int, string, string, string) (*gorm.DB, error){
		"postgresql": GetPgDb,
		"mysql":      GetMySqlDb,
	}
)

func GetDb(driver string, host string, port int, user string, password string, dbName string) (*gorm.DB, error) {
	return DbMap[driver](host, port, user, password, dbName)
}

package mock

import (
	"os"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DNSENV = "ARGO_DATABASE_TEST_URL"

func NewMySQL() *gorm.DB {
	dns := os.Getenv(DNSENV)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func ID(size int) string {
	id, _ := gonanoid.New(size)
	return id
}

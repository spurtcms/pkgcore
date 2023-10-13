package cms

import (
	"fmt"

	authority "github.com/spurtcms/spurtcms-core/auth"
	"github.com/spurtcms/spurtcms-core/teams"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// create instance
func NewInstance(a *authority.Option) authority.Authority {

	auth := authority.Authority{
		DB:     a.DB,
		Token:  a.Token,
		Secret: a.Secret,
	}

	authority.MigrationTable(auth.DB)
	teams.MigrateTables(auth.DB)

	return auth
}

// Create Database Instance
func DBInstance(host, username, password, databasename string) *gorm.DB {

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, username, databasename, password) //Build connection string

	// DB, err := gorm.Open("postgres", dbUri)
	DB, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		fmt.Println("Status:", err)
		panic(err)
	}
	return DB
}

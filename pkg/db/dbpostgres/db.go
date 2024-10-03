package dbpostgres

import (
	"fmt"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(c config.Postgres) (*gorm.DB, error) {
	dbSettings := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Username, c.DbName, c.Password)

	db, err := gorm.Open(postgres.Open(dbSettings), &gorm.Config{})
	tx := db.Debug()
	return tx, err
}

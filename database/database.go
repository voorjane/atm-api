package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"nedorez-test/pkg"
	"os"
)

type PgConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func ConnectToDB() (*gorm.DB, error) {
	var pg PgConfig
	setEnvVariables(&pg)

	cfg := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", pg.Host, pg.Username, pg.Password, pg.DBName, pg.Port)
	db, err := gorm.Open(postgres.Open(cfg), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&pkg.Account{})
	return db, nil
}

func setEnvVariables(cfg *PgConfig) {
	cfg.Host = os.Getenv("DATABASE_HOST")
	cfg.Port = os.Getenv("DATABASE_PORT")
	cfg.Username = os.Getenv("POSTGRES_USER")
	cfg.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.DBName = os.Getenv("POSTGRES_DB")
}

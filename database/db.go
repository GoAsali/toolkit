package database

import (
	"fmt"
	"github.com/goasali/toolkit/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Config struct {
	*config.Database
}

func loadConfig() *Config {
	if dbConfig, err := config.LoadDatabase(); err != nil {
		log.Fatalf("Error during parse enviroments for database config: %v", err)
	} else {
		if dbConfig.Host == "" {
			dbConfig.Host = "localhost"
		}
		return &Config{
			Database: dbConfig,
		}
	}

	return nil
}

func defaultDatabaseConfig() *gorm.Config {
	return &gorm.Config{}
}

func Database() (*gorm.DB, error) {
	var err error

	if db == nil {
		log.Info("Database connection instance does not exists, create one")
		log.Info("Connect to db...")

		dbConfig := loadConfig()
		if db, err = dbConfig.loadDatabase(); err != nil {
			return nil, err
		}

		log.Printf("Successfully connected to %s database", dbConfig.Type)
	}

	return db, err
}

func (c *Config) mysql() (*gorm.DB, error) {
	c.Port = "3306"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Name)
	return gorm.Open(mysql.Open(dsn), defaultDatabaseConfig())
}

func (c *Config) sqlite() (*gorm.DB, error) {
	fileName := fmt.Sprintf("%s.sqlite", c.Name)
	return gorm.Open(sqlite.Open(fileName), defaultDatabaseConfig())
}

func (c *Config) postgres() (*gorm.DB, error) {
	if c.Port == "" {
		c.Port = "5432"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", c.Host, c.Username, c.Password, c.Name, c.Port)
	return gorm.Open(postgres.Open(dsn), defaultDatabaseConfig())
}

func (c *Config) loadDatabase() (*gorm.DB, error) {
	switch c.Type {
	case "mysql":
		return c.mysql()
	case "sqlite":
		return c.sqlite()
	case "postgres":
		return c.postgres()
	}

	return nil, UnknownDbTypeError{c.Type}
}

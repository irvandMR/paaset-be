package config

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Config struct {
	Env   *Environment
	DB    *gorm.DB
	Redis *redis.Client
}

func InitConfig() *Config {
	env := LoadEnvironment()
	db := InitDatabaseConfig(env)
	RunMigrations(db)
	rdb := InitRedisConfig(env)

	return &Config{
		Env:   env,
		DB:    db,
		Redis: rdb,
	}
}

func (c *Config) Close() {
	if c.Redis != nil {
		_ = c.Redis.Close()
	}
	if c.DB != nil {
		if sqlDB, err := c.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}

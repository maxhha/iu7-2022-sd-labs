package gorm_repositories

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"iu7-2022-sd-labs/configuration"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GORMRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) GORMRepository {
	return GORMRepository{db}
}

func NewFromConfig(config configuration.GORMRepositoryConfig) (GORMRepository, error) {
	gormConfig := gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	}

	var db *gorm.DB
	var err error

	switch config.Database {
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.DSN), &gormConfig)
		err = Wrap(err, "gorm open")
	default:
		err = fmt.Errorf("unknown config database: %s", config.Database)
	}

	if err != nil {
		return GORMRepository{}, err
	}

	return New(db), nil
}

func (r *GORMRepository) Atomic(fn func(tx repositories.Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		repo := New(tx)
		return fn(&repo)
	})
}

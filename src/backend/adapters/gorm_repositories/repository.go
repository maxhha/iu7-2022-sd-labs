package gorm_repositories

import (
	"iu7-2022-sd-labs/buisness/ports/repositories"

	"gorm.io/gorm"
)

type GORMRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) GORMRepository {
	return GORMRepository{db}
}

func (r *GORMRepository) Consumer() repositories.ConsumerRepository {
	return nil
}

func (r *GORMRepository) Offer() repositories.OfferRepository {
	return nil
}

func (r *GORMRepository) Organizer() repositories.OrganizerRepository {
	return nil
}

func (r *GORMRepository) Product() repositories.ProductRepository {
	return nil
}

func (r *GORMRepository) Room() repositories.RoomRepository {
	return nil
}

func (r *GORMRepository) Atomic(fn func(tx repositories.Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		repo := NewGORMRepository(tx)
		return fn(&repo)
	})
}

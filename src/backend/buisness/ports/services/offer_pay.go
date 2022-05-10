package services

import "iu7-2022-sd-labs/buisness/entities"

type OfferPayService interface {
	PayLink(offer *entities.Offer) (string, error)
}

package offer_pay_service

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
)

type OfferPayService struct{}

func NewOfferPayService() OfferPayService {
	return OfferPayService{}
}

func (*OfferPayService) PayLink(offer *entities.Offer) (string, error) {
	return "", fmt.Errorf("pay service is unimplimented")
}

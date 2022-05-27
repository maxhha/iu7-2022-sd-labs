package ports

import (
	"context"
	"iu7-2022-sd-labs/buisness/entities"
)

type DataLoader interface {
	WithNewLoader(ctx context.Context) context.Context
	LoadAuction(ctx context.Context, id string) (entities.Auction, error)
	LoadManyAuctions(ctx context.Context, ids []string) ([]entities.Auction, error)
	LoadBidStepTable(ctx context.Context, id string) (entities.BidStepTable, error)
	LoadManyBidStepTables(ctx context.Context, ids []string) ([]entities.BidStepTable, error)
	LoadConsumer(ctx context.Context, id string) (entities.Consumer, error)
	LoadManyConsumers(ctx context.Context, ids []string) ([]entities.Consumer, error)
	LoadOffer(ctx context.Context, id string) (entities.Offer, error)
	LoadManyOffers(ctx context.Context, ids []string) ([]entities.Offer, error)
	LoadOrganizer(ctx context.Context, id string) (entities.Organizer, error)
	LoadManyOrganizers(ctx context.Context, ids []string) ([]entities.Organizer, error)
	LoadProduct(ctx context.Context, id string) (entities.Product, error)
	LoadManyProducts(ctx context.Context, ids []string) ([]entities.Product, error)
	LoadRoom(ctx context.Context, id string) (entities.Room, error)
	LoadManyRooms(ctx context.Context, ids []string) ([]entities.Room, error)
}

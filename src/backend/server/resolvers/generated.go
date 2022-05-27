package resolvers

//go:generate go run ../../codegen/graphql_resolvers/main.go --out auction.resolvers_gen.go --entity Auction --resolvers QueryPagination
//go:generate go run ../../codegen/graphql_resolvers/main.go --out bid_step_table.resolvers_gen.go --entity BidStepTable --resolvers QueryPagination
//go:generate go run ../../codegen/graphql_resolvers/main.go --out consumer.resolvers_gen.go --entity Consumer --resolvers QueryPagination
//go:generate go run ../../codegen/graphql_resolvers/main.go --out offer.resolvers_gen.go --entity Offer --resolvers QueryPagination
//go:generate go run ../../codegen/graphql_resolvers/main.go --out organizer.resolvers_gen.go --entity Organizer --resolvers QueryPagination
//go:generate go run ../../codegen/graphql_resolvers/main.go --out product.resolvers_gen.go --entity Product --resolvers QueryPagination
//go:generate go run ../../codegen/graphql_resolvers/main.go --out room.resolvers_gen.go --entity Room --resolvers QueryPagination

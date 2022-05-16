package resolvers

//go:generate go run ../../codegen/graphql_resolvers/main.go --out consumer.resolvers_gen.go --entity Consumer --resolvers QueryPagination
//go:generate go run ../../codegen/graphql_resolvers/main.go --out organizer.resolvers_gen.go --entity Organizer --resolvers QueryPagination
//go:generate go run ../../codegen/graphql_resolvers/main.go --out room.resolvers_gen.go --entity Room --resolvers QueryPagination

module github.com/daniellawrence/cv/backend/search

go 1.22

require (
	github.com/daniellawrence/cv/gen/go v0.0.0
	google.golang.org/protobuf v1.33.0
)

replace github.com/daniellawrence/cv/gen/go => ../../gen/go

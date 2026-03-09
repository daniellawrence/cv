package search

import (
	"testing"

	searchv1 "github.com/daniellawrence/cv/gen/go/search/v1"
	"google.golang.org/protobuf/proto"
)

func TestProto(t *testing.T) {

	req := &searchv1.SearchRequest{
		Query: "golang",
	}

	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	var out searchv1.SearchRequest

	err = proto.Unmarshal(data, &out)
	if err != nil {
		t.Fatal(err)
	}

	if out.Query != req.Query {
		t.Fatalf("expected %s got %s", req.Query, out.Query)
	}
}

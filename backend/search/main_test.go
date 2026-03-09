package search

import (
	"testing"

	pb "github.com/daniellawrence/cv/backend/search/proto"
	"google.golang.org/protobuf/proto"
)

func TestProto(t *testing.T) {

	req := &pb.SearchRequest{
		Query: "hello",
	}

	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	var out pb.SearchRequest

	err = proto.Unmarshal(data, &out)
	if err != nil {
		t.Fatal(err)
	}

	if out.Query != req.Query {
		t.Fatalf("expected %s got %s", req.Query, out.Query)
	}
}

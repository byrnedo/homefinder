package pontuz

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/byrnedo/homefinder/internal/pkg/agents"
)

func TestPontuz(t *testing.T) {
	p := Crawler{}

	l, err := p.GetForSale()
	t.Log(p.listings)
	if err != nil {
		t.Fatal(err)
	}

	if len(l) == 0 {
		t.Fatalf("wrong number of listings %d", len(l))
	}
	for _, listing := range l {
		b, _ := json.Marshal(listing)
		t.Log(string(b))

		if listing.Type == agents.ListingTypeUnknown {
			t.Fatalf("had unknown housing type: %v", listing)
		}
	}
	t.Log(l)
}

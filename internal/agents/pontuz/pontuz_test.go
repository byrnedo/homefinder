package pontuz

import (
	_ "embed"
	"gitlab.com/donalbyrne/homefinder/internal/agents"
	"testing"
)

//go:embed index.html
var testBody string

func TestPontuz(t *testing.T) {
	p := Crawler{body: testBody}

	l, err := p.GetForSale()
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 36 {
		t.Fatalf("wrong number of listings %d", len(l))
	}
	for _, listing := range l {
		if listing.Type == agents.ListingTypeUnknown {
			t.Fatalf("had unknown housing type: %v", listing)
		}
	}
	t.Log(l)
}

package maklarhuset

import (
	_ "embed"
	"testing"

	"github.com/byrnedo/homefinder/internal/pkg/agents"
)

func TestCrawler(t *testing.T) {
	p := Crawler{}

	l, err := p.GetForSale(agents.TargetBjelin)
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 13 {
		t.Fatalf("wrong number of listings: %d", len(l))
	}
	t.Log(l)
}

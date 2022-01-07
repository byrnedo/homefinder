package erikolsson

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
	t.Log(l)
}

package maklarhuset

import (
	_ "embed"
	"testing"
)

func TestCrawler(t *testing.T) {
	p := Crawler{}

	l, err := p.GetForSale()
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 13 {
		t.Fatalf("wrong number of listings: %d", len(l))
	}
	t.Log(l)
}

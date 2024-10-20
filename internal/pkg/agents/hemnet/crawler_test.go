package hemnet

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
	t.Log(l)
}

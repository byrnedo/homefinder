package svenskfast

import (
	_ "embed"
	"testing"
)

//go:embed index.html
var testBody string

func TestCrawler(t *testing.T) {
	p := Crawler{body: testBody}

	l, err := p.GetForSale()
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 10 {
		t.Fatalf("wrong number of listings %d", len(l))
	}
	t.Log(l)
}

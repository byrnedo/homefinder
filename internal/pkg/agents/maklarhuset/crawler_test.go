package maklarhuset

import (
	_ "embed"
	"testing"
)

//go:embed index-fstan.html
var testFstanBody string

//go:embed index-kalmar.html
var testKalmarBody string

func TestCrawler(t *testing.T) {
	p := Crawler{fstanBody: testFstanBody, kalmarBody: testKalmarBody}

	l, err := p.GetForSale()
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 3 {
		t.Fatalf("wrong number of listings: %d", len(l))
	}
	t.Log(l)
}

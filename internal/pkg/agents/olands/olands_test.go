package olands

import (
	_ "embed"
	"testing"
)

//go:embed index.html
var testBody string

func TestOlands(t *testing.T) {
	p := Crawler{body: testBody}

	l, err := p.GetForSale()
	if err != nil {
		t.Fatal(err)
	}

	if len(l) != 71 {
		t.Fatalf("wrong number of listings: %d", len(l))
	}
	t.Log(l)
}

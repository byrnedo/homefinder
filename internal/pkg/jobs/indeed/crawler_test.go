package indeed

import (
	_ "embed"
	"testing"
)

//go:embed index.html
var testBody string

func TestCrawler(t *testing.T) {
	p := Crawler{body: testBody}

	l, err := p.GetJobs()
	if err != nil {
		t.Fatal(err)
	}

	if len(l) == 0 {
		t.Fatalf("wrong number of listings %d", len(l))
	}
	t.Log(l)
}

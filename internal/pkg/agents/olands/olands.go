package olands

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	css "github.com/andybalholm/cascadia"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"github.com/byrnedo/homefinder/internal/pkg/xcss"
)

type Crawler struct {
	body string
}

func (o *Crawler) fetch() error {
	res, err := http.DefaultClient.Get("https://olandsmaklaren.se/till-salu-3")
	if err != nil {
		return nil
	}
	b, _ := ioutil.ReadAll(res.Body)
	o.body = string(b)
	return nil
}

func (o Crawler) Name() string {
	return "Ölandsmäklaren"
}

func (o *Crawler) GetForSale() (listings []agents.Listing, err error) {
	if o.body == "" {
		if err := o.fetch(); err != nil {
			return nil, err
		}
	}

	n, err := html.Parse(strings.NewReader(o.body))
	if err != nil {
		return nil, err
	}
	nodes := css.QueryAll(n, css.MustCompile("#Container>a"))
	if len(nodes) == 0 {
		return nil, xcss.NotFoundErr{}
	}
	var compressSpace = regexp.MustCompile(`\s+`)
	for _, n = range nodes {
		a := n
		img := css.Query(n, css.MustCompile(".oc-img-wrapper>img"))
		listing := agents.Listing{
			Link:  "https://olandsmaklaren.se" + xcss.FindAttr(a, "href"),
			Image: xcss.FindAttr(img, "src"),
		}

		title := xcss.CollectText(css.Query(n, css.MustCompile("h3")))
		listing.Name = title + " " + xcss.FindAttr(a, "href")

		var facts []string
		for _, f := range css.QueryAll(n, css.MustCompile("div.uk-tile>ul>li")) {

			raw := xcss.CollectText(f)
			raw = strings.ReplaceAll(raw, "\n", "")
			raw = strings.TrimSpace(raw)
			raw = compressSpace.ReplaceAllString(raw, " ")
			if strings.HasSuffix(raw, "kvm") {
				sqmStr := strings.TrimSpace(strings.TrimSuffix(raw, "kvm"))
				flVal, _ := strconv.ParseFloat(sqmStr, 32)
				listing.SquareMetres = int(flVal)
			}

			facts = append(facts, raw)
		}
		listing.Facts = facts
		listings = append(listings, listing)
	}

	return
}

package pontuz

import (
	css "github.com/andybalholm/cascadia"
	"github.com/byrnedo/homefinder/internal/agents"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Crawler struct {
	body string
}

func (p Crawler) Name() string {
	return "Pontuz Löfgren"
}

func (p *Crawler) fetch() error {

	res, err := http.DefaultClient.Get("https://www.pontuzlofgren.se/till-salu")
	if err != nil {
		return nil
	}
	b, _ := ioutil.ReadAll(res.Body)
	p.body = string(b)
	return nil

}

func (p *Crawler) GetForSale() (ls []agents.Listing, err error) {
	if p.body == "" {
		if err := p.fetch(); err != nil {
			return nil, err
		}
	}

	n, err := html.Parse(strings.NewReader(p.body))
	if err != nil {
		return nil, err
	}
	nodes := css.QueryAll(n, css.MustCompile("body>div.wrapper>div.ol-wrapper.container>div.col"))
	if len(nodes) == 0 {
		return nil, agents.NotFoundErr{"body>div.wrapper>div.ol-wrapper.container>div.col"}
	}
	var compressSpace = regexp.MustCompile(`\s+`)
	for _, n = range nodes {

		var listingType agents.ListingType
		if agents.HasClass(n, "house") {
			listingType = agents.ListingTypeHouse
		} else if agents.HasClass(n, "project") {
			listingType = agents.ListingTypeProject
		} else if agents.HasClass(n, "housingcooperative") {
			continue
		} else if agents.HasClass(n, "plot") {
			listingType = agents.ListingTypePlot
		} else {
			listingType = agents.ListingTypeUnknown
		}

		a := css.Query(n, css.MustCompile("a"))
		img := css.Query(n, css.MustCompile("img"))
		listing := agents.Listing{
			Upcoming: agents.CollectText(css.Query(n, css.MustCompile("div.oc-badge"))) == "I startblocken",
			Link:     agents.FindAttr(a, "href"),
			Image:    agents.FindAttr(img, "src"),
			Type:     listingType,
		}

		title := agents.CollectText(css.Query(n, css.MustCompile("h3.oc-title")))
		sub := agents.CollectText(css.Query(n, css.MustCompile("h4.oc-sub-title")))
		listing.Name = strings.Join([]string{title, sub}, " ")

		var facts []string
		for _, f := range css.QueryAll(n, css.MustCompile("div.oc-fact")) {
			raw := agents.CollectText(f)
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
		ls = append(ls, listing)
	}

	return
}
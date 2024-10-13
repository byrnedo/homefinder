package bjurfors

import (
	css "github.com/andybalholm/cascadia"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"github.com/byrnedo/homefinder/internal/pkg/xcss"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Crawler struct {
}

func (c Crawler) Name() string {
	return "Bjurfors"
}

func (c Crawler) GetForSale() (listings []agents.Listing, err error) {

	address := "https://www.bjurfors.se/sv/tillsalu/kalmar-lan/morbylanga/farjestaden/?qdata=d%245evjg4l9otbn1ijo&FormId=6362f9d8-928e-4fe3-8abc-a8f157a65244"

	res, err := http.Get(address)
	if err != nil {
		return nil, err
	}

	n, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	nodes := css.QueryAll(n, css.MustCompile("ul.c-search-result__list>li"))
	if len(nodes) == 0 {
		return nil, xcss.NotFoundErr{}
	}
	var compressSpace = regexp.MustCompile(`\s+`)
	for _, n = range nodes {

		a := css.Query(n, css.MustCompile("a"))
		img := css.Query(n, css.MustCompile("picture>img"))
		listing := agents.Listing{
			Link:  "https://bjurfors.se" + xcss.FindAttr(a, "href"),
			Image: "https://bjurfors.se" + xcss.FindAttr(img, "src"),
		}

		title := xcss.CollectText(css.Query(n, css.MustCompile("h3")))
		listing.Name = title + " " + xcss.FindAttr(a, "href")

		var facts []string
		for _, f := range css.QueryAll(n, css.MustCompile("ul.c-object-card__meta-list>li")) {
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

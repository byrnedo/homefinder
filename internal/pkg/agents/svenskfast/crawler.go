package svenskfast

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	css "github.com/andybalholm/cascadia"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"github.com/byrnedo/homefinder/internal/pkg/xcss"
	"golang.org/x/net/html"
)

type Crawler struct {
	body string
}

func (p Crawler) Name() string {
	return "Svenskfastihetsförmedling"
}

func (p *Crawler) fetch(target agents.Target) error {

	u := "https://www.svenskfast.se/hus/kalmar/kalmar/kalmar/?t=Villa,Radhus,Fritidshus,Nyproduktionsprojekt,Lantbruk,Tomt&l=morbylanga/farjestaden"

	res, err := http.DefaultClient.Get(u)
	if err != nil {
		return nil
	}
	b, _ := ioutil.ReadAll(res.Body)
	p.body = string(b)
	return nil

}

func (p *Crawler) GetForSale(target agents.Target) (ls []agents.Listing, err error) {

	if p.body == "" {
		if err := p.fetch(target); err != nil {
			return nil, err
		}
	}

	n, err := html.Parse(strings.NewReader(p.body))
	if err != nil {
		return nil, err
	}
	nodes := css.QueryAll(n, css.MustCompile("body div.search__results div.grid__item"))
	if len(nodes) == 0 {
		return nil, xcss.NotFoundErr{Name: "grid__item"}
	}
	for _, n = range nodes {

		a := css.Query(n, css.MustCompile("a"))
		img := css.Query(n, css.MustCompile("div.search-hit__image"))
		listing := agents.Listing{
			Upcoming: xcss.CollectText(css.Query(n, css.MustCompile("div.oc-badge"))) == "I startblocken",
			Link:     "https://www.svenskfast.se" + xcss.FindAttr(a, "href"),
			Image:    xcss.FindAttr(img, "data-src"),
			Name:     xcss.CleanText(xcss.CollectText(css.Query(n, css.MustCompile("div.search-hit__address")))),
		}

		var facts []string
		for _, f := range css.QueryAll(n, css.MustCompile("div.search-hit__info--text>span")) {

			i := css.Query(f, css.MustCompile("i"))
			if i != nil {
				class := xcss.FindAttr(i, "class")
				classes := strings.Fields(class)
				for _, c := range classes {
					if strings.HasPrefix(c, "icon__") {
						objectType := strings.TrimPrefix(c, "icon__")
						switch objectType {
						case "villa":
							listing.Type = agents.ListingTypeHouse
						case "lagenhet":
							listing.Type = agents.ListingTypeApartment
						case "fritidhus":
							listing.Type = agents.ListingTypeSummerHouse
						case "tomt":
							listing.Type = agents.ListingTypePlot
						case "lantbruk":
							listing.Type = agents.ListingTypeFarm
						case "nyproduktion":
							facts = append(facts, "new-build")
						}
					}
				}
			}

			raw := xcss.CleanText(xcss.CollectText(f))

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

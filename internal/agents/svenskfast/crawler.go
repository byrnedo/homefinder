package svenskfast

import (
	css "github.com/andybalholm/cascadia"
	"gitlab.com/donalbyrne/homefinder/internal/agents"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Crawler struct {
	body string
}

func (p Crawler) Name() string {
	return "Svenskfastihetsförmedling"
}

func (p *Crawler) fetch() error {

	res, err := http.DefaultClient.Get("https://www.svenskfast.se/hus/kalmar/kalmar/kalmar/?t=Radhus,Fritidshus,Nyproduktionsprojekt,Lantbruk,Tomt&l=kalmar/morbylanga/farjestaden")
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
	nodes := css.QueryAll(n, css.MustCompile("body div.search__results div.grid__item"))
	if len(nodes) == 0 {
		return nil, agents.NotFoundErr{}
	}
	for _, n = range nodes {

		a := css.Query(n, css.MustCompile("a"))
		img := css.Query(n, css.MustCompile("div.search-hit__image"))
		listing := agents.Listing{
			Upcoming: agents.CollectText(css.Query(n, css.MustCompile("div.oc-badge"))) == "I startblocken",
			Link:     "https://www.svenskfast.se" + agents.FindAttr(a, "href"),
			Image:    agents.FindAttr(img, "data-src"),
			Name:     agents.CleanText(agents.CollectText(css.Query(n, css.MustCompile("div.search-hit__address")))),
		}

		var facts []string
		for _, f := range css.QueryAll(n, css.MustCompile("div.search-hit__info--text>span")) {

			i := css.Query(f, css.MustCompile("i"))
			if i != nil {
				class := agents.FindAttr(i, "class")
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

			raw := agents.CleanText(agents.CollectText(f))

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

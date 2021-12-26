package maklarhuset

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	css "github.com/andybalholm/cascadia"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"golang.org/x/net/html"
)

type Crawler struct {
	fstanBody  string
	kalmarBody string
}

func (o *Crawler) fetch() error {
	res, err := http.DefaultClient.Get("https://www.maklarhuset.se/bostad/sverige/kalmar-lans-lan/farjestaden-kommun?undefined")
	if err != nil {
		return nil
	}
	b, _ := ioutil.ReadAll(res.Body)
	o.fstanBody = string(b)
	res, err = http.DefaultClient.Get("https://www.maklarhuset.se/bostad/sverige/kalmar-lans-lan/kalmar-kommun?undefined")
	if err != nil {
		return nil
	}
	b, _ = ioutil.ReadAll(res.Body)
	o.kalmarBody = string(b)
	return nil
}

func (o Crawler) Name() string {
	return "Mäklarhuset"
}

func (o *Crawler) GetForSale() (listings []agents.Listing, err error) {
	if o.fstanBody == "" {
		if err := o.fetch(); err != nil {
			return nil, err
		}
	}

	for _, body := range []string{o.fstanBody, o.kalmarBody} {

		n, err := html.Parse(strings.NewReader(body))
		if err != nil {
			return nil, err
		}

		nodes := css.QueryAll(n, css.MustCompile("div.uk-width-large-1-2 a"))

		var compressSpace = regexp.MustCompile(`\s+`)

		for _, n = range nodes {
			n = n.Parent
			dataObj := css.Query(n.Parent, css.MustCompile("div[data-object-id]"))
			if dataObj == nil {
				break
			}
			a := css.Query(n, css.MustCompile("a"))

			listing := agents.Listing{
				Link:  "https://www.maklarhuset.se" + agents.FindAttr(a, "href"),
				Image: agents.FindAttr(dataObj, "data-object-image"),
			}

			title := agents.FindAttr(dataObj, "data-object-address") + " " + agents.FindAttr(dataObj, "data-object-city") + "(" + agents.FindAttr(dataObj, "data-object-id") + ")"

			listing.Name = title

			var facts []string
			for _, f := range css.QueryAll(n, css.MustCompile("figcaption>div.uk-h3>span")) {

				raw := agents.CollectText(f)
				raw = strings.ReplaceAll(raw, "\n", "")
				raw = strings.TrimSpace(raw)
				raw = compressSpace.ReplaceAllString(raw, " ")
				switch raw {
				case "Villa":
					listing.Type = agents.ListingTypeHouse
				case "Fritidshus":
					listing.Type = agents.ListingTypeSummerHouse
				case "Gård":
					listing.Type = agents.ListingTypeFarm
				case "Tomt":
					listing.Type = agents.ListingTypePlot
				case "Lägenhet":
					listing.Type = agents.ListingTypeApartment
				case "Övrigt":
					listing.Type = agents.ListingTypeUnknown
				case "Radhus":
					listing.Type = agents.ListingTypeTerrace

				}
				if listing.SquareMetres == 0 {
					if strings.HasSuffix(raw, "kvm") {
						sqmStr := strings.ReplaceAll(strings.TrimSpace(strings.TrimSuffix(raw, "kvm")), " ", "")
						flVal, _ := strconv.ParseFloat(sqmStr, 32)
						listing.SquareMetres = int(flVal)
					}
				}

				facts = append(facts, raw)
			}
			if listing.Type == agents.ListingTypeApartment {
				continue
			}
			listing.Facts = facts
			listings = append(listings, listing)
		}
	}

	return
}

package maklarhuset

import (
	"io/ioutil"
	"log"
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
	bodies    []string
	addresses []string
}

func (o *Crawler) fetch() error {

	o.addresses = []string{
		"https://www.maklarhuset.se/bostad/sverige/kalmar-lans-lan/farjestaden-kommun?undefined",
	}
	for _, u := range o.addresses {

		res, err := http.DefaultClient.Get(u)
		if err != nil {
			return nil
		}
		func() {
			defer res.Body.Close()
			b, _ := ioutil.ReadAll(res.Body)
			o.bodies = append(o.bodies, string(b))
		}()
	}
	return nil
}

func (o Crawler) Name() string {
	return "Mäklarhuset"
}

func (o *Crawler) GetForSale() (listings []agents.Listing, err error) {
	if err := o.fetch(); err != nil {
		return nil, err
	}

	for i, body := range o.bodies {
		log.Printf("iter %d, url: %s\n", i, o.addresses[i])

		n, err := html.Parse(strings.NewReader(body))
		if err != nil {
			return nil, err
		}

		nodes := css.QueryAll(n, css.MustCompile("div.uk-width-large-1-2 a"))

		var compressSpace = regexp.MustCompile(`\s+`)

		var urlListings []agents.Listing
		for _, n = range nodes {
			n = n.Parent
			dataObj := css.Query(n.Parent, css.MustCompile("div[data-object-id]"))
			if dataObj == nil {
				break
			}
			a := css.Query(n, css.MustCompile("a"))

			listing := agents.Listing{
				Link:  "https://www.maklarhuset.se" + xcss.FindAttr(a, "href"),
				Image: xcss.FindAttr(dataObj, "data-object-image"),
			}

			title := xcss.FindAttr(dataObj, "data-object-address") + " " + xcss.FindAttr(dataObj, "data-object-city") + "(" + xcss.FindAttr(dataObj, "data-object-id") + ")"

			listing.Name = title

			var facts []string
			for _, f := range css.QueryAll(n, css.MustCompile("figcaption>div.uk-h3>span")) {

				raw := xcss.CollectText(f)
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

			if len(facts) > 0 && strings.HasSuffix(facts[0], " kr") {
				priceStr := xcss.RemoveSpace(strings.TrimSuffix(facts[0], " kr"))
				listing.Price, _ = strconv.Atoi(priceStr)
			}
			urlListings = append(urlListings, listing)
		}
		listings = append(listings, urlListings...)
		log.Println(len(urlListings), "listings found")
	}

	return
}

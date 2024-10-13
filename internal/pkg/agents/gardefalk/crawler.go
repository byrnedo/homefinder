package gardefalk

import (
	css "github.com/andybalholm/cascadia"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"github.com/byrnedo/homefinder/internal/pkg/xcss"
	"golang.org/x/net/html"
	"log"
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

	for _, address := range []string{"https://gardefalk.se/vara-objekt/", "https://gardefalk.se/pa-vag/"} {

		res, err := http.Get(address)
		if err != nil {
			return nil, err
		}

		n, err := html.Parse(res.Body)
		if err != nil {
			return nil, err
		}
		nodes := css.QueryAll(n, css.MustCompile("div.project_box"))
		if len(nodes) == 0 {
			log.Println("no projects")
			return listings, nil
		}
		var compressSpace = regexp.MustCompile(`\s+`)
		for _, n = range nodes {

			a := css.Query(n, css.MustCompile("a"))
			img := css.Query(n, css.MustCompile("img"))
			listing := agents.Listing{
				Link:  xcss.FindAttr(a, "href"),
				Image: xcss.FindAttr(img, "src"),
			}

			title := xcss.CollectText(css.Query(n, css.MustCompile("p.addRess")))
			listing.Name = title + " " + xcss.FindAttr(a, "href")

			var facts []string
			for _, f := range css.QueryAll(n, css.MustCompile("ul.propertyDetail>li")) {
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
	}

	return

}

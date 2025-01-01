package hemnet

import (
	"fmt"
	css "github.com/andybalholm/cascadia"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"github.com/byrnedo/homefinder/internal/pkg/xcss"
	"golang.org/x/net/html"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Crawler struct {
}

func (c Crawler) Name() string {
	return "Hemnet"
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (c Crawler) GetForSale() (listings []agents.Listing, err error) {

	for _, address := range []string{"https://www.hemnet.se/bostader?item_types%5B%5D=villa&item_types%5B%5D=tomt&item_types%5B%5D=gard&location_ids%5B%5D=939457"} {

		res, err := http.Get(address)
		if err != nil {
			return nil, err
		}

		body, _ := io.ReadAll(res.Body)
		defer res.Body.Close()

		bodyStr := strings.ReplaceAll(string(body), "noscript", "div")

		n, err := html.Parse(strings.NewReader(bodyStr))
		if err != nil {
			return nil, err
		}
		nodes := css.QueryAll(n, css.MustCompile("div[data-testid=\"result-list\"]>a"))
		if len(nodes) == 0 {
			log.Println("no projects")
			return listings, nil
		}
		var compressSpace = regexp.MustCompile(`\s+`)
		for _, n = range nodes {
			a := n
			img := css.Query(n, css.MustCompile("img[src$='jpg']"))
			listing := agents.Listing{
				Link:  "https://hemnet.se" + xcss.FindAttr(a, "href"),
				Image: xcss.FindAttr(img, "src"),
			}

			title := xcss.CollectText(css.Query(n, css.MustCompile("div[class^='Header_title']")))
			listing.Name = fmt.Sprintf("%s %d", title, hash(listing.Link))

			var facts []string
			for _, f := range css.QueryAll(n, css.MustCompile("span[class^='ForSaleAttributes_primaryAttributes']")) {
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

package rydmanlanga

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
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
	return "Rydman Langå"
}

func (p *Crawler) fetch() error {

	res, err := http.DefaultClient.Get("https://www.rydmanlanga.se/till-salu")
	if err != nil {
		return nil
	}
	b, _ := ioutil.ReadAll(res.Body)
	p.body = string(b)
	return nil

}

func (p *Crawler) GetForSale(target agents.Target) (ls []agents.Listing, err error) {
	if p.body == "" {
		if err := p.fetch(); err != nil {
			return nil, err
		}
	}

	n, err := html.Parse(strings.NewReader(p.body))
	if err != nil {
		return nil, err
	}
	nodes := css.QueryAll(n, css.MustCompile("body div.ol-wrapper div.col"))
	if len(nodes) == 0 {
		return nil, xcss.NotFoundErr{}
	}
	var compressSpace = regexp.MustCompile(`\s+`)
	for _, n = range nodes {
		a := css.Query(n, css.MustCompile("a"))
		img := css.Query(n, css.MustCompile("img"))
		listing := agents.Listing{
			Upcoming: xcss.CollectText(css.Query(n, css.MustCompile("div.oc-badge"))) == "I startblocken",
			Link:     xcss.FindAttr(a, "href"),
			Image:    xcss.FindAttr(img, "src"),
		}

		title := xcss.CollectText(css.Query(n, css.MustCompile("h3.oc-title")))

		hit := false
		for _, term := range []string{"FÄRJESTADEN", "RUNSBÄCK"} {
			if strings.Contains(strings.ToUpper(title), term) {
				hit = true
				break
			}
		}

		if !hit {
			log.Printf("ignoring %s", title)
			continue
		}

		sub := xcss.CollectText(css.Query(n, css.MustCompile("h4.oc-sub-title")))
		listing.Name = strings.Join([]string{title, sub}, " ")

		var facts []string
		for _, f := range css.QueryAll(n, css.MustCompile("div.oc-fact")) {
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
		ls = append(ls, listing)
	}

	return
}

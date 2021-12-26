package agents

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type ListingType string

const (
	ListingTypeHouse       ListingType = "house"
	ListingTypeSummerHouse ListingType = "summerhouse"
	ListingTypeFarm        ListingType = "farm"
	ListingTypeTerrace     ListingType = "terracehouse"
	ListingTypeProject     ListingType = "project"
	ListingTypeUnknown     ListingType = "unknown"
	ListingTypePlot        ListingType = "plot"
	ListingTypeApartment   ListingType = "apartment"
)

type Listing struct {
	Name         string
	Link         string
	Type         ListingType
	Image        string
	Upcoming     bool
	Facts        []string
	SquareMetres int
}

type NotFoundErr struct {
	Name string
}

func (n NotFoundErr) Error() string {
	return "node not found: " + n.Name
}

func CollectText(n *html.Node) (c string) {
	if n == nil {
		return ""
	}
	if n.Type == html.TextNode {
		c += n.Data
	}
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		c += " " + CollectText(ch)
	}
	return strings.TrimSpace(c)
}

func FindAttr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

var compressSpace = regexp.MustCompile(`\s+`)

func CleanText(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = compressSpace.ReplaceAllString(raw, " ")
	return raw
}

func HasClass(n *html.Node, name string) bool {
	for _, a := range n.Attr {
		if a.Key == "class" {
			classes := strings.Fields(a.Val)
			for _, class := range classes {
				if strings.EqualFold(class, name) {
					return true
				}
			}
		}
	}
	return false
}

type Crawler interface {
	Name() string
	GetForSale() ([]Listing, error)
}

package xcss

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

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

var compressSpace = regexp.MustCompile(`\s+`)

func CleanText(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = compressSpace.ReplaceAllString(raw, " ")
	return raw
}
func RemoveSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
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
	if n == nil {
		return ""
	}
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

type NotFoundErr struct {
	Name string
}

func (n NotFoundErr) Error() string {
	return "node not found: " + n.Name
}

package erikolsson

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"github.com/byrnedo/homefinder/internal/pkg/xcss"
)

type Crawler struct {
}

func (c Crawler) Name() string {
	return "Erik Olsson"
}

type property struct {
	VitecObjectId         string      `json:"vitecObjectId"`
	StartPrice            int         `json:"startPrice"`
	MainImageUrl          string      `json:"mainImageUrl"`
	City                  string      `json:"city"`
	AreaName              string      `json:"areaName"`
	Address               string      `json:"address"`
	Price                 string      `json:"price"`
	Rooms                 string      `json:"rooms"`
	ShowPuff              bool        `json:"showPuff"`
	PuffText              string      `json:"puffText"`
	ShowPrice             bool        `json:"showPrice"`
	PrioOrder             int         `json:"prioOrder"`
	PublishedDate         string      `json:"publishedDate"`
	Url                   string      `json:"url"`
	Agency                string      `json:"agency"`
	Fee                   interface{} `json:"fee"`
	Type                  string      `json:"type"`
	IsApartment           bool        `json:"isApartment"`
	PriceText             string      `json:"priceText"`
	ShouldShowOnlyAddress bool        `json:"shouldShowOnlyAddress"`
	NumberOfRooms         float64     `json:"numberOfRooms"`
	HidePrice             bool        `json:"hidePrice"`
}

type response struct {
	ShouldRenderShowMoreButton      bool       `json:"shouldRenderShowMoreButton"`
	Type                            int        `json:"type"`
	Properties                      []property `json:"properties"`
	Hits                            int        `json:"hits"`
	AreaName                        int        `json:"areaName"`
	ShouldListSearchInfo            bool       `json:"shouldListSearchInfo"`
	ShowAsPrioLow                   bool       `json:"showAsPrioLow"`
	ShouldRenderInThreeGrid         bool       `json:"shouldRenderInThreeGrid"`
	ShouldRenderInTwoGrid           bool       `json:"shouldRenderInTwoGrid"`
	ShouldRenderTwoThirdInThreeGrid bool       `json:"shouldRenderTwoThirdInThreeGrid"`
	ShouldRenderOneThirdInThreeGrid bool       `json:"shouldRenderOneThirdInThreeGrid"`
}

func (c Crawler) GetForSale(target agents.Target) (listings []agents.Listing, err error) {

	address := "https://www.erikolsson.se/api/search/?areaids=117205-Kalmar,Kalmar|117224-Färjestaden,Mörbylånga&propertytype=villa,tomt-mark,radhus,parhus,kedjehus&internalOnly=true"
	if target == agents.TargetBjelin {
		address = "https://www.erikolsson.se/api/search/?areaids=201-H%C3%A4rryda%2C%20V%C3%A4stra%20G%C3%B6talands%20l%C3%A4n%7C268-Partille%2C%20V%C3%A4stra%20G%C3%B6talands%20l%C3%A4n%7C208-Kungsbacka%2C%20Hallands%20l%C3%A4n%7C72-Lerum%2C%20V%C3%A4stra%20G%C3%B6talands%20l%C3%A4n%7C259-M%C3%B6lndal%2C%20V%C3%A4stra%20G%C3%B6talands%20l%C3%A4n&propertytype=villa&minarea=75&maxprice=5000000"
	}

	res, err := http.Get(address)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	result := response{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, item := range result.Properties {
		l := agents.Listing{
			Name:         strings.Join([]string{item.Address, item.AreaName + "(" + item.VitecObjectId + ")"}, ","),
			Link:         "https://www.erikolsson.se" + item.Url,
			Type:         c.parseType(item),
			Image:        item.MainImageUrl,
			Upcoming:     strings.EqualFold(item.PuffText, "kommer snart"),
			Facts:        strings.Split(strings.TrimSuffix(item.Rooms, "."), ","),
			SquareMetres: 0,
		}
		l.Price, _ = strconv.Atoi(strings.TrimSuffix(xcss.RemoveSpace(item.Price), "kr"))
		listings = append(listings, l)
	}
	return
}

func (c Crawler) parseType(p property) agents.ListingType {
	switch strings.ToLower(p.Type) {
	case "villa":
		return agents.ListingTypeHouse
	case "tomt-mark":
		return agents.ListingTypePlot
	case "radhus", "parhus", "kedjehus":
		return agents.ListingTypeTerrace
	default:
		return agents.ListingTypeUnknown
	}
}

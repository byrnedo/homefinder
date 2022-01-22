package lanfast

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/byrnedo/homefinder/internal/pkg/agents"
)

type Crawler struct {
}

func (c Crawler) Name() string {
	return "Lasnsförsäkring"
}

type estate struct {
	Url                     string      `json:"url"`
	AllYear                 interface{} `json:"allYear"`
	StreetAddress           string      `json:"streetAddress"`
	StartPrice              float64     `json:"startPrice"`
	FinalPrice              float64     `json:"finalPrice"`
	PriceMin                float64     `json:"priceMin"`
	PriceMax                float64     `json:"priceMax"`
	MonthlyCost             float64     `json:"monthlyCost"`
	LivingSpace             string      `json:"livingSpace"`
	LivingSpaceMax          float64     `json:"livingSpaceMax"`
	LivingSpaceMin          float64     `json:"livingSpaceMin"`
	OtherSpace              string      `json:"otherSpace"`
	PlotSize                float64     `json:"plotSize"`
	NumberOfRooms           string      `json:"numberOfRooms"`
	NumberOfRoomsMax        float64     `json:"numberOfRoomsMax"`
	NumberOfRoomsMin        float64     `json:"numberOfRoomsMin"`
	IsProjectEstate         bool        `json:"isProjectEstate"`
	ProjectName             interface{} `json:"projectName"`
	EstateType              string      `json:"estateType"`
	City                    string      `json:"city"`
	Municipality            string      `json:"municipality"`
	Area                    string      `json:"area"`
	NextViewing             *string     `json:"nextViewing"`
	NextNextViewing         interface{} `json:"nextNextViewing"`
	BiddingText             string      `json:"biddingText"`
	Status                  string      `json:"status"`
	NumberOfEstates         *string     `json:"numberOfEstates"`
	PropertyArea            string      `json:"propertyArea"`
	LdPlus                  bool        `json:"ldPlus"`
	HeadStartIcon           bool        `json:"headStartIcon"`
	HeaderImage             string      `json:"headerImage"`
	HasBids                 bool        `json:"hasBids"`
	PropertyUnitDesignation string      `json:"propertyUnitDesignation"`
}
type response struct {
	Estates     []estate `json:"estates"`
	TotalLength int      `json:"totalLength"`
}

func (c Crawler) GetForSale(target agents.Target) (listings []agents.Listing, err error) {

	address := "https://www.lansfast.se/umbraco/api/findestateapi/loadestates?municipality=Kalmar&areaIds=CMOMRADE4SBR9DBI7G5JTKOQ&areaNames=Färjestaden&estateType=Villa,Bostadsrätt,Fritidshus,Tomt,Nyproduktion&sortOrder=0"

	res, err := http.Get(address)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	result := response{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, item := range result.Estates {
		l := agents.Listing{
			Name:         strings.Join([]string{item.StreetAddress, item.City}, ","),
			Link:         "https://www.lansfast.se" + item.Url,
			Type:         c.parseType(item),
			Image:        "https://www.lansfast.se" + item.HeaderImage,
			Upcoming:     strings.EqualFold(item.Status, "Kommande"),
			Facts:        []string{item.NumberOfRooms + "rum", fmt.Sprintf("%fkvm", item.PlotSize), item.LivingSpace + "kvm"},
			SquareMetres: 0,
		}
		if strings.HasPrefix(item.Url, "http") {
			l.Link = item.Url
		}
		l.Price = int(item.StartPrice)
		listings = append(listings, l)
	}
	return
}

func (c Crawler) parseType(p estate) agents.ListingType {
	switch strings.ToLower(p.EstateType) {
	case "villa":
		return agents.ListingTypeHouse
	case "tomt":
		return agents.ListingTypePlot
	case "båstadsrätt":
		return agents.ListingTypeApartment
	default:
		return agents.ListingTypeUnknown
	}
}

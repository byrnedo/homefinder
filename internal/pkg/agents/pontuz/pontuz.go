package pontuz

import (
	"encoding/json"
	"fmt"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"net/http"
	"strings"
)

type Crawler struct {
	listings []NotarListing
}

func (p *Crawler) Name() string {
	return "Pontuz Löfgren"
}

func (p *Crawler) fetch() error {

	urls := []string{
		"https://data.notar.se/objects?limit=25&sortBy=date&sortOrder=desc&areaId=CMOMRADE2TKJ0RM871EHKICV&areaId=CMOMRADE3512A4TKMP4G0PBS&areaId=CMOMRADE3UJIGL6OQ1L18VHK&areaId=CMOMRADE3H9C0TK9HH31J44L&assignmentStatus=Kommande&assignmentStatus=Till%20salu&assignmentStatus=P%C3%A5g%C3%A5ende&assignmentStatus=Reserverad",
	}

	for _, url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Basic bm90YXI6OWd6ekRmRGZCVXpjMjJleVhqUG53bU5m")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("bad response: %d", res.StatusCode)
		}
		defer res.Body.Close()
		result := []NotarListing{}
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return err
		}
		p.listings = append(p.listings, result...)
	}
	return nil

}

func getListingType(notarType string) agents.ListingType {
	switch strings.ToLower(notarType) {
	case "house":
		return agents.ListingTypeHouse
	case "housingcooperative":
		return agents.ListingTypeApartment
	case "plot":
		return agents.ListingTypePlot
	case "project":
		return agents.ListingTypeProject
	default:
		return agents.ListingTypeUnknown
	}
}

func (p *Crawler) GetForSale() (listings []agents.Listing, err error) {
	if err := p.fetch(); err != nil {
		return nil, err
	}

	for _, notarListing := range p.listings {
		b, _ := json.Marshal(notarListing)
		println(string(b))
		l := agents.Listing{
			Name:     strings.ReplaceAll(fmt.Sprintf("%s (%s)", notarListing.Address, notarListing.ID), "\n", ""),
			Link:     "https://www.notar.se/kopa-bostad/objekt/" + notarListing.ID,
			Type:     getListingType(notarListing.Type),
			Upcoming: false,
			Facts: []string{
				fmt.Sprintf("%dm2", notarListing.LivingSpaceArea),
				fmt.Sprintf("%d rooms", notarListing.NumberOfRooms),
				fmt.Sprintf("%dm2 plot", notarListing.Plot.Area),
			},
			SquareMetres: notarListing.LivingSpaceArea,
			Price:        notarListing.Price.StartingPrice,
		}
		if len(notarListing.Images) > 0 {
			l.Image = fmt.Sprintf("https://notar-media.b-cdn.net/objects/%s/images/%s/JPG/840/%s", notarListing.AgentID, notarListing.CustomerID, notarListing.Images[0].File)
		}
		listings = append(listings, l)
	}
	return
}

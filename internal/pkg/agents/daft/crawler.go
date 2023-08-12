package daft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Crawler struct{}

func (c Crawler) Name() string {
	return "Daft"
}

type Filter struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type GeoFilter struct {
	StoredShapeIds []string `json:"storedShapeIds"`
	GeoSearchType  string   `json:"geoSearchType"`
}

type Paging struct {
	From     string `json:"from"`
	PageSize string `json:"pageSize"`
}

type ListingsRequest struct {
	Section    string        `json:"section"`
	Filters    []Filter      `json:"filters"`
	AndFilters []Filter      `json:"andFilters"`
	Ranges     []interface{} `json:"ranges"`
	Paging     Paging        `json:"paging"`
	GeoFilter  GeoFilter     `json:"geoFilter"`
	Terms      string        `json:"terms"`
}
type ListingsResponse struct {
	Listings []struct {
		Listing struct {
			Id                int      `json:"id"`
			Title             string   `json:"title"`
			SeoTitle          string   `json:"seoTitle"`
			Sections          []string `json:"sections"`
			SaleType          []string `json:"saleType"`
			FeaturedLevel     string   `json:"featuredLevel"`
			FeaturedLevelFull string   `json:"featuredLevelFull"`
			PublishDate       int64    `json:"publishDate"`
			Price             string   `json:"price"`
			AbbreviatedPrice  string   `json:"abbreviatedPrice"`
			NumBedrooms       string   `json:"numBedrooms"`
			PropertyType      string   `json:"propertyType"`
			DaftShortcode     string   `json:"daftShortcode"`
			Seller            struct {
				SellerId             int    `json:"sellerId"`
				Name                 string `json:"name"`
				Phone                string `json:"phone,omitempty"`
				Branch               string `json:"branch,omitempty"`
				BackgroundColour     string `json:"backgroundColour,omitempty"`
				SellerType           string `json:"sellerType"`
				ShowContactForm      bool   `json:"showContactForm"`
				PremierPartnerSeller bool   `json:"premierPartnerSeller"`
				PhoneWhenToCall      string `json:"phoneWhenToCall,omitempty"`
			} `json:"seller"`
			DateOfConstruction string `json:"dateOfConstruction"`
			Media              struct {
				Images []struct {
					Size720X480 string `json:"size720x480"`
					Size600X600 string `json:"size600x600"`
					Size400X300 string `json:"size400x300"`
					Size360X240 string `json:"size360x240"`
					Size300X200 string `json:"size300x200"`
					Size320X280 string `json:"size320x280"`
					Size72X52   string `json:"size72x52"`
					Size680X392 string `json:"size680x392"`
				} `json:"images"`
				TotalImages    int  `json:"totalImages"`
				HasVideo       bool `json:"hasVideo"`
				HasVirtualTour bool `json:"hasVirtualTour"`
				HasBrochure    bool `json:"hasBrochure"`
			} `json:"media"`
			Ber struct {
				Rating string `json:"rating"`
				Code   string `json:"code,omitempty"`
			} `json:"ber"`
			Platform string `json:"platform"`
			Point    struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"point"`
			SeoFriendlyPath string `json:"seoFriendlyPath"`
			Prs             struct {
				TotalUnitTypes int `json:"totalUnitTypes"`
				SubUnits       []struct {
					Id              int    `json:"id"`
					Price           string `json:"price"`
					NumBedrooms     string `json:"numBedrooms"`
					BathroomType    string `json:"bathroomType"`
					PropertyType    string `json:"propertyType"`
					DaftShortcode   string `json:"daftShortcode"`
					SeoFriendlyPath string `json:"seoFriendlyPath"`
					Category        string `json:"category"`
					Image           struct {
						Size1440X960  string `json:"size1440x960"`
						Size1200X1200 string `json:"size1200x1200"`
						Size720X480   string `json:"size720x480"`
						Size600X600   string `json:"size600x600"`
						Size400X300   string `json:"size400x300"`
						Size360X240   string `json:"size360x240"`
						Size300X200   string `json:"size300x200"`
						Size320X280   string `json:"size320x280"`
						Size72X52     string `json:"size72x52"`
						Size680X392   string `json:"size680x392"`
					} `json:"image"`
					Media struct {
						Images []struct {
							Size720X480 string `json:"size720x480"`
							Size600X600 string `json:"size600x600"`
							Size400X300 string `json:"size400x300"`
							Size360X240 string `json:"size360x240"`
							Size300X200 string `json:"size300x200"`
							Size320X280 string `json:"size320x280"`
							Size72X52   string `json:"size72x52"`
							Size680X392 string `json:"size680x392"`
						} `json:"images"`
						TotalImages    int  `json:"totalImages"`
						HasVideo       bool `json:"hasVideo"`
						HasVirtualTour bool `json:"hasVirtualTour"`
						HasBrochure    bool `json:"hasBrochure"`
					} `json:"media"`
					Ber struct {
						Rating string `json:"rating"`
					} `json:"ber"`
				} `json:"subUnits"`
				TagLine string `json:"tagLine"`
			} `json:"prs,omitempty"`
			PageBranding struct {
				BackgroundColour string        `json:"backgroundColour,omitempty"`
				SquareLogos      []interface{} `json:"squareLogos"`
			} `json:"pageBranding,omitempty"`
			Category       string `json:"category"`
			State          string `json:"state"`
			PremierPartner bool   `json:"premierPartner"`
		} `json:"listing"`
		SavedAd bool `json:"savedAd"`
	} `json:"listings"`
	ShowcaseListings []interface{} `json:"showcaseListings"`
	Paging           struct {
		TotalPages     int `json:"totalPages"`
		CurrentPage    int `json:"currentPage"`
		NextFrom       int `json:"nextFrom"`
		PreviousFrom   int `json:"previousFrom"`
		DisplayingFrom int `json:"displayingFrom"`
		DisplayingTo   int `json:"displayingTo"`
		TotalResults   int `json:"totalResults"`
		PageSize       int `json:"pageSize"`
	} `json:"paging"`
	DfpTargetingValues struct {
		PageType       []string `json:"pageType"`
		SearchPageNo   []string `json:"searchPageNo"`
		AreaName       []string `json:"areaName"`
		AdState        []string `json:"adState"`
		DistilledBrand []string `json:"distilledBrand"`
		Section        []string `json:"section"`
		IsPledge       []string `json:"isPledge"`
		CountyName     []string `json:"countyName"`
		IsUserLoggedIn []string `json:"isUserLoggedIn"`
	} `json:"dfpTargetingValues"`
	Breadcrumbs []struct {
		DisplayValue string `json:"displayValue"`
		Url          string `json:"url"`
	} `json:"breadcrumbs"`
	CanonicalUrl string `json:"canonicalUrl"`
	MapView      bool   `json:"mapView"`
	SavedSearch  bool   `json:"savedSearch"`
}

func (c Crawler) GetForSale(target agents.Target) (listings []agents.Listing, err error) {
	return c.getForSalePage(0)
}
func (c Crawler) getForSalePage(from int) (listings []agents.Listing, err error) {
	const daftUrl = "https://gateway.daft.ie/old/v1/listings"

	reqBytes, _ := json.Marshal(ListingsRequest{
		Section: "sharing",
		Filters: []Filter{
			{Name: "adState", Values: []string{"published"}},
		},
		AndFilters: []Filter{},
		Paging: Paging{
			From:     fmt.Sprintf("%d", from),
			PageSize: "50",
		},
		GeoFilter: GeoFilter{
			StoredShapeIds: []string{"19"},
			GeoSearchType:  "STORED_SHAPES",
		},
		Ranges: []interface{}{},
		Terms:  "",
	})
	req, err := http.NewRequest("POST", daftUrl, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Version", "1.0.1691579911")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("Platform", "web")
	req.Header.Set("Brand", "daft")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.daft.ie/")
	req.Header.Set("Origin", "https://www.daft.ie")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.6")

	b, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(b))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		wholeBody, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%s", string(wholeBody))
	}
	result := ListingsResponse{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, daftListing := range result.Listings {
		l := agents.Listing{
			Name:         fmt.Sprintf("%s (%d)", daftListing.Listing.Title, daftListing.Listing.Id),
			Link:         "http://daft.ie" + daftListing.Listing.SeoFriendlyPath,
			Type:         getListingType(daftListing.Listing.PropertyType),
			Upcoming:     false,
			Facts:        []string{daftListing.Listing.Price},
			SquareMetres: 0,
			Price:        0, // TODO
		}
		if len(daftListing.Listing.Media.Images) > 0 {
			l.Image = daftListing.Listing.Media.Images[0].Size360X240
		}
		listings = append(listings, l)
	}
	if result.Paging.NextFrom > 0 {
		sub, err := c.getForSalePage(result.Paging.NextFrom)
		if err != nil {
			return nil, err
		}
		listings = append(listings, sub...)
	}
	return
}

func getListingType(daftType string) agents.ListingType {
	switch strings.ToLower(daftType) {
	case "apartments":
		return agents.ListingTypeApartment
	case "house":
		return agents.ListingTypeHouse
	case "sharing":
		return agents.ListingTypeSharing
	}
	return agents.ListingTypeUnknown
}

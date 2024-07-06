package agents

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
	ListingTypeSharing     ListingType = "sharing"
)

type Listing struct {
	Name         string
	Link         string
	Type         ListingType
	Image        string
	Upcoming     bool
	Facts        []string
	SquareMetres int
	Price        int
}

type Crawler interface {
	Name() string
	GetForSale() ([]Listing, error)
}

type Target string

const (
	TargetSjobyrne Target = "sjobyrne"
)

package jobs

type ListingType string

const (
	ListingTypeContract   ListingType = "contract"
	ListingTypeEmployment ListingType = "employment"
)

type Listing struct {
	ID       string
	Name     string
	Link     string
	Type     ListingType
	Image    string
	Company  string
	Location string
	Facts    []string
}
type Crawler interface {
	Name() string
	GetJobs() ([]Listing, error)
}

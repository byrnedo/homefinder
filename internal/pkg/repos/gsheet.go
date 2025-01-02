package repos

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"strconv"
	"time"
)

type GsheetRepo struct {
	svc           *sheets.Service
	SpreadsheetID string
	SheetID       int
}

func NewGsheetRepo(spreadsheetID string, sheetID int) *GsheetRepo {
	return &GsheetRepo{
		SpreadsheetID: spreadsheetID,
		SheetID:       sheetID,
	}
}

func (c *GsheetRepo) GetHistory(ctx context.Context) (listings []agents.Listing, err error) {
	req := c.svc.Spreadsheets.Values.Get(c.SpreadsheetID, fmt.Sprintf("Sheet%d!A:Z", c.SheetID+1))
	res, err := req.Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	for i, row := range res.Values {
		l, err := rowToListing(row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert row %d to listing: %w", i, err)
		}
		listings = append(listings, *l)

	}
	return
}

func (c *GsheetRepo) SaveHistory(ctx context.Context, listings []agents.Listing) error {
	rows := make([]*sheets.RowData, len(listings))

	now := time.Now()

	for i, s := range listings {

		values := listingToCells(now, s)

		rows[i] = &sheets.RowData{
			Values: values,
		}
	}

	req := &sheets.BatchUpdateSpreadsheetRequest{
		IncludeSpreadsheetInResponse: false,
	}

	req.Requests = append(req.Requests, &sheets.Request{
		AppendCells: &sheets.AppendCellsRequest{
			Fields:  "*",
			Rows:    rows,
			SheetId: int64(c.SheetID),
		},
	})

	call := c.svc.Spreadsheets.BatchUpdate(c.SpreadsheetID, req)
	_, err := call.Context(ctx).Do()
	return err
}

func ref[T any](v T) *T {
	return &v
}

func rowToListing(row []any) (*agents.Listing, error) {
	if len(row) < 8 {
		return nil, fmt.Errorf("row too short")
	}

	listing := agents.Listing{}

	listing.Crawler = row[0].(string)
	listing.Name = row[1].(string)
	listing.Image = row[2].(string)
	listing.Link = row[3].(string)
	listing.Price, _ = strconv.Atoi(row[4].(string))
	listing.SquareMetres, _ = strconv.Atoi(row[5].(string))
	listing.Type = agents.ListingType(row[6].(string))

	return &listing, nil

}
func listingToCells(now time.Time, s agents.Listing) (cells []*sheets.CellData) {

	cells = append(cells,
		&sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: &s.Crawler,
			},
		},
		&sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: &s.Name,
			},
		},
		&sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: &s.Image,
			},
		},
		&sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: &s.Link,
			},
		},
		&sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				NumberValue: ref(float64(s.Price)),
			},
		},
		&sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				NumberValue: ref(float64(s.SquareMetres)),
			},
		},
		&sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: ref(string(s.Type)),
			},
		},
		&sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: ref(now.UTC().Format(time.RFC3339)),
			},
		},
	)

	return cells
}

func (c *GsheetRepo) Authenticate(ctx context.Context, base64Key string) error {
	credBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return err
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return err
	}

	client := config.Client(ctx)

	c.svc, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	return nil
}

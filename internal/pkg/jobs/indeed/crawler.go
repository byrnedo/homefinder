//go:generate go run github.com/Khan/genqlient@latest
package indeed

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/byrnedo/homefinder/internal/pkg/jobs"
)

type Crawler struct {
	body string
}

func (p Crawler) Name() string {
	return "Indeed.com"
}

type authedTransport struct {
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("indeed-api-key", "9e8a2fd82bf5e73386b6c07f12a8a80d422a3b1b9730981cff5f312751ae0554")
	return t.wrapped.RoundTrip(req)
}

func (p *Crawler) GetJobs() (ls []jobs.Listing, err error) {
	ctx := context.Background()
	client := graphql.NewClient("https://apis.indeed.com/graphql?co=SE", &http.Client{Transport: &authedTransport{wrapped: http.DefaultTransport}})

	var cursor *string
	for {
		fmt.Printf("cursor `%v`\n", cursor)
		resp, err := SearchJobs(ctx, client, cursor, nil)
		if err != nil {
			return nil, err
		}

		if len(resp.JobSearch.Results) == 0 {
			break
		}

		for _, j := range resp.JobSearch.Results {
			if j.Job.Expired {
				continue
			}

			facts := []string{}
			for _, attr := range j.Job.Attributes {
				facts = append(facts, attr.Label)
			}

			listing := jobs.Listing{
				ID:       j.Job.Key,
				Name:     j.Job.Title,
				Link:     j.Job.Url,
				Type:     "",
				Company:  j.Job.Employer.Name,
				Location: strings.Join([]string{j.Job.Location.City, j.Job.Location.CountryCode}, ","),
				Facts:    facts,
			}

			ls = append(ls, listing)

		}

		if len(resp.JobSearch.PageInfo.NextPages) == 0 {
			break
		}
		cursor = &resp.JobSearch.PageInfo.NextPages[0].Cursor
	}

	return
}

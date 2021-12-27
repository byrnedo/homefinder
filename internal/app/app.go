package app

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"github.com/byrnedo/homefinder/internal/pkg/agents/erikolsson"
	"github.com/byrnedo/homefinder/internal/pkg/agents/fastighetsbyran"
	"github.com/byrnedo/homefinder/internal/pkg/agents/maklarhuset"
	"github.com/byrnedo/homefinder/internal/pkg/agents/olands"
	"github.com/byrnedo/homefinder/internal/pkg/agents/pontuz"
	"github.com/byrnedo/homefinder/internal/pkg/agents/rydmanlanga"
	"github.com/byrnedo/homefinder/internal/pkg/agents/svenskfast"
	"github.com/byrnedo/homefinder/internal/pkg/repos"

	"github.com/slack-go/slack"
)

const FileName = "/tmp/listings-seen"

func Run(ctx context.Context, historyRepo repos.HistoryRepo) {
	crawlers := []agents.Crawler{
		&pontuz.Crawler{},
		&rydmanlanga.Crawler{},
		&fastighetsbyran.Crawler{},
		&olands.Crawler{},
		&svenskfast.Crawler{},
		&maklarhuset.Crawler{},
		&erikolsson.Crawler{},
	}

	prevListings, err := historyRepo.GetHistory(ctx)
	if err != nil {
		panic("failed to load from disk:" + err.Error())
	}

	curListings := map[string]repos.Void{}

	var newListings []agents.Listing
	for _, c := range crawlers {
		log.Println("checking " + c.Name() + "...")
		listings, err := c.GetForSale()
		if err != nil {
			panic(err)
		}
		log.Printf("found %d listings for %s\n", len(listings), c.Name())

		for _, listing := range listings {

			curListings[c.Name()+":"+listing.Name] = repos.Void{}

			if _, ok := prevListings[c.Name()+":"+listing.Name]; !ok {
				// new listings
				newListings = append(newListings, listing)
				//
			}
		}
	}

	log.Printf("found %d new listings\n", len(newListings))

	if len(newListings) > 0 {
		var blocks []slack.Block
		for i, l := range newListings {
			facts := strings.Join(l.Facts, " - ")
			if facts == "" {
				facts = "-"
			}
			if l.Type == "" {
				l.Type = "-"
			}
			blocks = append(blocks, slack.SectionBlock{
				Type: slack.MBTSection,
				Fields: []*slack.TextBlockObject{
					{
						Type: slack.MarkdownType,
						Text: l.Name,
					},
					{
						Type: slack.MarkdownType,
						Text: facts,
					},
					{
						Type: slack.MarkdownType,
						Text: l.Link,
					},
					{
						Type: slack.MarkdownType,
						Text: string(l.Type),
					},
				},
				Accessory: &slack.Accessory{
					ImageElement: &slack.ImageBlockElement{
						Type:     slack.METImage,
						ImageURL: l.Image,
						AltText:  l.Name,
					},
				},
			})
			if i > 0 && i%50 == 0 {
				msg := &slack.WebhookMessage{
					Blocks: &slack.Blocks{
						BlockSet: blocks,
					},
				}
				b, _ := json.MarshalIndent(msg, "", "  ")
				log.Println(string(b))
				err = slack.PostWebhookContext(ctx, os.Getenv("SLACK_WEBHOOK_URL"), msg)
				blocks = nil
			}
		}
		if blocks != nil {
			msg := &slack.WebhookMessage{
				Blocks: &slack.Blocks{
					BlockSet: blocks,
				},
			}
			b, _ := json.MarshalIndent(msg, "", "  ")
			log.Println(string(b))
			err = slack.PostWebhookContext(ctx, os.Getenv("SLACK_WEBHOOK_URL"), msg)
		}
	}

	if err != nil {
		panic(err)
	}

	if err := historyRepo.SaveHistory(ctx, curListings); err != nil {
		panic("failed to save to disk:" + err.Error())
	}

}

func loadFromDisk() (map[string]bool, error) {
	f, err := os.OpenFile(FileName, os.O_RDONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]bool{}, nil
		}
		return nil, err
	}
	defer f.Close()

	m := map[string]bool{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		m[text] = true

	}
	return m, nil
}

func saveToDisk(ls map[string]bool) error {
	f, err := os.OpenFile(FileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Truncate(0)
	f.Seek(0, 0)
	writer := bufio.NewWriter(f)
	for k, _ := range ls {

		writer.WriteString(k + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

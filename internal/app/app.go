package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/byrnedo/homefinder/internal/pkg/agents"
	"github.com/byrnedo/homefinder/internal/pkg/agents/erikolsson"
	"github.com/byrnedo/homefinder/internal/pkg/agents/fastighetsbyran"
	"github.com/byrnedo/homefinder/internal/pkg/agents/lanfast"
	"github.com/byrnedo/homefinder/internal/pkg/agents/maklarhuset"
	"github.com/byrnedo/homefinder/internal/pkg/agents/olands"
	"github.com/byrnedo/homefinder/internal/pkg/agents/pontuz"
	"github.com/byrnedo/homefinder/internal/pkg/agents/rydmanlanga"
	"github.com/byrnedo/homefinder/internal/pkg/agents/svenskfast"
	"github.com/byrnedo/homefinder/internal/pkg/repos"

	"github.com/slack-go/slack"
)

func dashIfEmpty(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func RunHousefinder(ctx context.Context, historyRepo repos.HistoryRepo) error {
	crawlers := []agents.Crawler{
		&pontuz.Crawler{},
		&rydmanlanga.Crawler{},
		&fastighetsbyran.Crawler{},
		&olands.Crawler{},
		&svenskfast.Crawler{},
		&maklarhuset.Crawler{},
		&erikolsson.Crawler{},
		&lanfast.Crawler{},
	}

	prevListings, err := historyRepo.GetHistory(ctx)
	if err != nil {
		return fmt.Errorf("failed to load from disk: %s", err)
	}

	curListings := map[string]repos.Void{}

	var newListings []agents.Listing
	for _, c := range crawlers {
		log.Println("checking " + c.Name() + "...")
		listings, err := c.GetForSale(agents.TargetSjobyrne)
		if err != nil {
			return err
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

	if err := sendBlocksToSlack(ctx, newListings, ""); err != nil {
		return err
	}

	if err := historyRepo.SaveHistory(ctx, curListings); err != nil {
		return err
	}
	return nil
}

func sendBlocksToSlack(ctx context.Context, newListings []agents.Listing, channel string) error {
	if len(newListings) == 0 {
		return nil
	}
	var blocks []slack.Block
	for i, l := range newListings {
		blocks = append(blocks, homeToBlock(l))
		if i > 0 && i%49 == 0 {
			if err := postToSlack(ctx, blocks, nil, channel); err != nil {
				return err
			}
			blocks = nil
		}
	}
	if blocks != nil {
		if err := postToSlack(ctx, blocks, nil, channel); err != nil {
			return err
		}
	}
	return nil
}

func homeToBlock(l agents.Listing) slack.Block {
	facts := strings.Join(l.Facts, " - ")
	if facts == "" {
		facts = "-"
	}
	return slack.SectionBlock{
		Type: slack.MBTSection,
		Fields: []*slack.TextBlockObject{
			{
				Type: slack.MarkdownType,
				Text: dashIfEmpty(l.Name),
			},
			{
				Type: slack.MarkdownType,
				Text: dashIfEmpty(facts),
			},
			{
				Type: slack.MarkdownType,
				Text: dashIfEmpty(l.Link),
			},
			{
				Type: slack.MarkdownType,
				Text: dashIfEmpty(string(l.Type)),
			},
		},
		Accessory: &slack.Accessory{
			ImageElement: &slack.ImageBlockElement{
				Type:     slack.METImage,
				ImageURL: dashIfEmpty(l.Image),
				AltText:  dashIfEmpty(l.Name),
			},
		},
	}
}

func postToSlack(ctx context.Context, blocks []slack.Block, attachments []slack.Attachment, channel string) error {
	msg := &slack.WebhookMessage{
		Channel: channel,
		Blocks: &slack.Blocks{
			BlockSet: blocks,
		},
		Attachments: attachments,
	}
	return slack.PostWebhookContext(ctx, os.Getenv("SLACK_WEBHOOK_URL"), msg)
}

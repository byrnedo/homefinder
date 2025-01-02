package job

import (
	"context"
	"errors"
	"fmt"
	"github.com/byrnedo/homefinder/internal/pkg/agents/bjurfors"
	"github.com/byrnedo/homefinder/internal/pkg/agents/gardefalk"
	"github.com/byrnedo/homefinder/internal/pkg/agents/hemnet"
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

type crawlerConf struct {
	agents.Crawler
	channel string
}

func RunHousefinder(ctx context.Context, historyRepo repos.HistoryRepo, postToSlack bool) error {
	crawlers := []crawlerConf{
		{Crawler: &pontuz.Crawler{}, channel: "#oland"},
		{Crawler: &fastighetsbyran.Crawler{}, channel: "#oland"},
		{Crawler: &olands.Crawler{}, channel: "#oland"},
		{Crawler: &svenskfast.Crawler{}, channel: "#oland"},
		{Crawler: &maklarhuset.Crawler{}, channel: "#oland"},
		{Crawler: &erikolsson.Crawler{}, channel: "#oland"},
		{Crawler: &lanfast.Crawler{}, channel: "#oland"},
		{Crawler: &bjurfors.Crawler{}, channel: "#oland"},
		{Crawler: &gardefalk.Crawler{}, channel: "#oland"},
		{Crawler: &hemnet.Crawler{}, channel: "#oland"},
		//{Crawler: &daft.Crawler{}, channel: "#daft"},
	}

	prevListings, err := historyRepo.GetHistory(ctx)
	if err != nil {
		return fmt.Errorf("failed to load from disk: %s", err)
	}

	var curListings []agents.Listing

	var crawlerErrs error

	newListings := map[string][]agents.Listing{}
	newListingsFlat := []agents.Listing{}
	newCount := 0
	for _, c := range crawlers {
		log.Println("checking " + c.Name() + "...")
		listings, err := c.GetForSale()
		if err != nil {
			crawlerErrs = errors.Join(crawlerErrs, fmt.Errorf("%s: %w", c.Name(), err))
			log.Printf("ERROR (%s): %s", c.Name(), err)
			continue
		}
		log.Printf("found %d listings for %s\n", len(listings), c.Name())

		for _, listing := range listings {
			listing.Crawler = c.Name()

			curListings = append(curListings, listing)

			// puke
			alreadyRecorded := false
			for _, prevListing := range prevListings {
				if prevListing.Name == listing.Name && prevListing.Crawler == listing.Crawler {
					alreadyRecorded = true
					break
				}
			}

			if !alreadyRecorded {
				newListings[c.channel] = append(newListings[c.channel], listing)
				newListingsFlat = append(newListingsFlat, listing)
				newCount++
			}
		}
	}

	log.Printf("found %d new listings\n", newCount)

	if postToSlack {

		if err := sendBlocksToSlack(ctx, newListings); err != nil {
			return err
		}

	}

	if err := historyRepo.SaveHistory(ctx, newListingsFlat); err != nil {
		return err
	}
	return crawlerErrs
}

func sendBlocksToSlack(ctx context.Context, newListings map[string][]agents.Listing) error {
	if len(newListings) == 0 {
		return nil
	}
	for channel, listings := range newListings {
		var blocks []slack.Block
		for i, l := range listings {
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

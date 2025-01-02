package main

import (
	"context"
	"github.com/byrnedo/homefinder/internal/app/job"
	"os"

	"github.com/byrnedo/homefinder/internal/pkg/repos"
)

func mustEnv(name string) string {
	if v, found := os.LookupEnv(name); !found {
		panic("Missing environment variable: " + name)
	} else {
		return v
	}
}

func gSheetRepo() repos.HistoryRepo {

	cCredentials := mustEnv("CREDENTIALS")
	cSpreadsheetID := mustEnv("SPREADSHEET_ID")

	repo := repos.NewGsheetRepo(cSpreadsheetID, 0)
	err := repo.Authenticate(context.Background(), cCredentials)
	if err != nil {
		panic(err)
	}
	return repo
}

func main() {

	repo := gSheetRepo()

	err := job.RunHousefinder(context.Background(), repo, true)
	if err != nil {
		panic(err)
	}
}

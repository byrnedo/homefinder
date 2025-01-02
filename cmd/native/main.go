package main

import (
	"context"
	"os"

	"github.com/byrnedo/homefinder/internal/app"
	"github.com/byrnedo/homefinder/internal/pkg/repos"
)

func mustEnv(name string) string {
	if v, found := os.LookupEnv(name); !found {
		panic("Missing environment variable: " + name)
	} else {
		return v
	}
}

func main() {

	cCredentials := mustEnv("CREDENTIALS")
	cSpreadsheetID := mustEnv("SPREADSHEET_ID")

	repo := repos.NewGsheetRepo(cSpreadsheetID, 0)
	err := repo.Authenticate(context.Background(), cCredentials)
	if err != nil {
		panic(err)
	}

	err = app.RunHousefinder(context.Background(), repo, false)
	if err != nil {
		panic(err)
	}
}

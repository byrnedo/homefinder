package main

import (
	"context"

	"github.com/byrnedo/homefinder/internal/app"
	"github.com/byrnedo/homefinder/internal/pkg/repos"
)

func main() {
	err := app.RunHousefinder(context.Background(), repos.EmptyHistoryRepo{})
	if err != nil {
		panic(err)
	}
	err = app.RunJobfinder(context.Background(), repos.EmptyHistoryRepo{})
	if err != nil {
		panic(err)
	}
}

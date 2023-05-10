package main

import (
	"context"

	"github.com/byrnedo/homefinder/internal/app"
	"github.com/byrnedo/homefinder/internal/pkg/repos"
)

func main() {
	err := app.RunHousefinder(context.Background(), repos.FileHistoryRepo{"./cache/cache.txt"})
	if err != nil {
		panic(err)
	}
}

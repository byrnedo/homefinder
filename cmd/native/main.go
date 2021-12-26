package main

import (
	"context"

	"github.com/byrnedo/homefinder/internal/app"
	"github.com/byrnedo/homefinder/internal/pkg/repos"
)

func main() {
	app.Run(context.Background(), repos.EmptyHistoryRepo{})
}

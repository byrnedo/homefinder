package main

import (
	"context"
	"gitlab.com/donalbyrne/homefinder/internal/app"
)

func main() {
	app.Run(context.Background())
}

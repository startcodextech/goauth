package application

import (
	"context"
	"github.com/modernice/goes/helper/streams"
	"log"
)

func LogErrors(ctx context.Context, errs ...<-chan error) {
	log.Printf("Logging errors ...")

	in := streams.FanInContext(ctx, errs...)
	for err := range in {
		log.Println(err)
	}
}

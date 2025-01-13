package generator

import (
	"context"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"golang.org/x/exp/rand"
)

func RandomizeUint(ctx context.Context, values []uint) uint {
	segment := logger.StartSegment(ctx, "generator.RandomizeUint")
	defer segment.End()

	if len(values) == 0 {
		return 0
	}

	// Seed the random number generator (convert int64 to uint64)
	rand.Seed(uint64(time.Now().UnixNano()))

	// Generate a random index
	randomIndex := rand.Intn(len(values))

	// Return the random uint value
	return values[randomIndex]
}

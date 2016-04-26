package aws

import (
	"fmt"
	"math"
	"path/filepath"
	"time"
)

func formatTransferRate(value float64) string {
	v := int64(math.Floor(value))
	switch {
	case v > (1000 * 1000 * 1000 * 1024):
		return fmt.Sprintf("%d TB", v/(1000*1000*1000*1024))
	case v > (1000 * 1000 * 1024):
		return fmt.Sprintf("%d GB", v/(1000*1000*1024))
	case v > (1000 * 1024):
		return fmt.Sprintf("%d MB", v/(1000*1024))
	case v > 1024:
		return fmt.Sprintf("%d KB", v/1024)
	default:
		return fmt.Sprintf("%d b")
	}
}

func timeTrack(start time.Time, name string, size int64) {
	elapsed := time.Since(start).Seconds()
	// formatTransferRate(float64(size) / elapsed.Seconds())
	fmt.Printf("%s took %.4fs at %s/s\n", filepath.Base(name), elapsed, formatTransferRate(float64(size)/elapsed))
}

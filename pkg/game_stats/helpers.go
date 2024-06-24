package game_stats

import (
	"fmt"
	"time"
)

func parseTimestamp(timestampUnits []int) (time.Time, error) {
	if len(timestampUnits) != 6 {
		return time.Time{}, fmt.Errorf("invalid number of timestamp units: %d", len(timestampUnits))
	}

	return time.Date(timestampUnits[2], time.Month(timestampUnits[0]), timestampUnits[1], timestampUnits[3], timestampUnits[4], timestampUnits[5], 0, time.Local), nil
}

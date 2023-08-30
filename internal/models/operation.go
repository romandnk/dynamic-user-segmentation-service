package models

import "time"

type Operation struct {
	UserID      int
	SegmentSlug string
	Date        time.Time
	Action      string
}

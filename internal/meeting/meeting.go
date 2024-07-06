package meeting

import (
	"time"
)

type Meeting struct {
	MeetingID uint       `json:"meeting_id"`
	HostID    uint       `json:"host_id"`
	Code      string     `json:"code"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Active    bool       `json:"active"`
}

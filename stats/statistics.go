package stats

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

var stats StatsType

type StatsType struct {
	Tm_start            time.Time
	HttpRequests        int64
	HttpFailedRequests  int64
	HttpUnknownRequests int64
}

func Init() {
	stats.Tm_start = time.Now()
}

func Stats() (string, error) {
	local_stats := StatsType{
		Tm_start:            stats.Tm_start,
		HttpRequests:        atomic.LoadInt64(&stats.HttpRequests),
		HttpFailedRequests:  atomic.LoadInt64(&stats.HttpFailedRequests),
		HttpUnknownRequests: atomic.LoadInt64(&stats.HttpUnknownRequests),
	}

	if bin, err := json.Marshal(local_stats); err != nil {
		return "", err
	} else {
		return string(bin), nil
	}
}

func HttpRequestsInc() {
	atomic.AddInt64(&stats.HttpRequests, 1)
}

func HttpFailedRequestsInc() {
	atomic.AddInt64(&stats.HttpFailedRequests, 1)
}

func HttpUnknownRequestsInc() {
	atomic.AddInt64(&stats.HttpUnknownRequests, 1)
}

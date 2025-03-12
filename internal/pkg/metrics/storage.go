package metrics

import (
	"sync"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type Storage struct {
	mu                        sync.RWMutex
	totalEvents               int64
	eventTimestamps           []time.Time
	playerBets                map[int]int
	playerWins                map[int]int
	playerDeposits            map[int]int64
	// lastWindowCount           int
	// lastWindowStart           time.Time
	// movingAverageDataPoints   []float64
}

type Metrics struct {
    TotalEvents                 int64            `json:"events_total"`
    EventsPerMinute             float64          `json:"events_per_minute"`
    EventsPerSecondMovingAverage float64         `json:"events_per_second_moving_average"`
    TopPlayerBets               PlayerStat       `json:"top_player_bets"`
    TopPlayerWins               PlayerStat       `json:"top_player_wins"`
    TopPlayerDeposits           PlayerStatDeposit `json:"top_player_deposits"`
}

type PlayerStat struct {
    ID    int `json:"id"`
    Count int `json:"count"`
}

type PlayerStatDeposit struct {
    ID    int   `json:"id"`
    Count int64 `json:"count"`
}

func NewStorage() *Storage {
	return &Storage{
		eventTimestamps: make([]time.Time, 0),
		playerBets:      make(map[int]int),
		playerWins:      make(map[int]int),
		playerDeposits:  make(map[int]int64),
	}
}

func (s *Storage) RecordEvent(event casino.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	
	s.totalEvents++
	
	s.eventTimestamps = append(s.eventTimestamps, now)
	
	switch event.Type {
	case "bet":
		s.playerBets[event.PlayerID]++
		if event.HasWon {
			s.playerWins[event.PlayerID]++
		}
	case "deposit":
		s.playerDeposits[event.PlayerID] += int64(event.AmountEUR)
	}
	
	if len(s.eventTimestamps) > 0 && now.Sub(s.eventTimestamps[0]) > 2*time.Minute {
		cutoff := now.Add(-2 * time.Minute)
		var i int
		for i = 0; i < len(s.eventTimestamps) && s.eventTimestamps[i].Before(cutoff); i++ {}
		s.eventTimestamps = s.eventTimestamps[i:]
	}
}

func (s *Storage) GetMetrics() Metrics {
	s.mu.RLock()
    defer s.mu.RUnlock()

	now := time.Now()
    metrics := Metrics{
        TotalEvents: s.totalEvents,
    }

	if len(s.eventTimestamps) > 0 {
		oldest := s.eventTimestamps[0]
		minutes := now.Sub(oldest).Minutes()
		if minutes > 0 {
			metrics.EventsPerMinute = float64(len(s.eventTimestamps)) / minutes
		}
	}

	var count int
	cutoff := now.Add(-1 * time.Minute)
	for _, ts := range s.eventTimestamps {
		if ts.After(cutoff) {
			count++
		}
	}
	metrics.EventsPerSecondMovingAverage = float64(count) / 60.0

	metrics.TopPlayerBets = s.findMaxKeyValue(s.playerBets)
    metrics.TopPlayerWins = s.findMaxKeyValue(s.playerWins)
    metrics.TopPlayerDeposits = s.findMaxKeyValue64(s.playerDeposits)

    return metrics
}

func (s *Storage) findMaxKeyValue(m map[int]int) PlayerStat {
    result := PlayerStat{}
    for id, count := range m {
        if count > result.Count {
            result.ID = id
            result.Count = count
        }
    }
    return result
}

func (s *Storage) findMaxKeyValue64(m map[int]int64) PlayerStatDeposit {
    result := PlayerStatDeposit{}
    for id, count := range m {
        if count > result.Count {
            result.ID = id
            result.Count = count
        }
    }
    return result
}


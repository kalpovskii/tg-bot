package req

import (
	"crypto-prices-bot/database"
	"errors"
	"fmt"
	"time"
)

type UserStats struct {
	ID int64

	FirstReq  time.Time
	ReqAmount int64
	LastPair  string
}

func ShowStats(id int64) (UserStats, error) {
	stats, err := database.UserStats(id)
	if err != nil {
		return UserStats{}, errors.New(
			fmt.Sprintf("failed to get statistics for user %d: %s", id, err))
	}

	return UserStats{
		ID:        stats.ID,
		FirstReq:  stats.FirstReq,
		ReqAmount: stats.ReqAmount,
		LastPair:  stats.LastPair,
	}, nil
}

func RecordStats(id int64, pair string) error {
	if !database.StatsExist(id) {
		err := database.New(id)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create stats for user: %d: %s", id, err))
		}

		err = database.FirstReq(id)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to record request for user: %d: %s", id, err))
		}

		if !database.StatsExist(id) {
			return errors.New(fmt.Sprintf("failed to create stats for user: %d", id))
		}
	}

	err := database.IncreaseReq(id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to increase request amount for user %d: %s", id, err))
	}

	err = database.LastPair(id, pair)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to record last pair for user %d: %s", id, err))
	}

	return nil
}

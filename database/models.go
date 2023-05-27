package database

import "time"

type UserModel struct {
	ID        int64
	FirstReq  time.Time
	ReqAmount int64
	LastPair  string
}

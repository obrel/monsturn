package pdb

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Stat struct {
	Realm      string `json:"realm"`
	Username   string `json:"username"`
	Allocation string `json:"allocation"`
	SentP      int64  `json:"sentp"`
	SentB      int64  `json:"sentb"`
	RecvP      int64  `json:"recvp"`
	RecvB      int64  `json:"recvb"`
}

func NewStat(realm, username, allocation string, sentp, sentb, recvp, recvb int64) *Stat {
	return &Stat{
		realm,
		username,
		allocation,
		sentp,
		sentb,
		recvp,
		recvb,
	}
}

func (s *Stat) Insert() error {
	query := fmt.Sprintf(
		"INSERT INTO statistics (realm, username, allocation, sentp, sentb, recvp, recvb) VALUES ('%v', '%v', '%v', %v, %v, %v, %v) RETURNING id",
		s.Realm, s.Username, s.Allocation, s.SentP, s.SentB, s.RecvP, s.RecvB)

	err := tx(func(db *sqlx.Tx) error {
		_ = db.MustExec(query)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

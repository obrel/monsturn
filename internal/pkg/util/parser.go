package util

import (
	"strconv"
	"strings"

	"github.com/obrel/monsturn/internal/data"
)

func MessageParser(msg string) (*data.TurnStat, error) {
	stat := &data.TurnStat{}
	data := strings.Split(msg, ",")

	for _, dt := range data {
		str := strings.Split(strings.Trim(dt, " "), "=")
		val, err := strconv.ParseInt(str[1], 10, 64)
		if err != nil {
			return nil, err
		}

		switch str[0] {
		case "sentp":
			stat.SentP = val
		case "recvp":
			stat.RecvP = val
		case "sentb":
			stat.SentB = val
		case "recvb":
			stat.RecvB = val
		}
	}

	return stat, nil
}

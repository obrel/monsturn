package util

import (
	"errors"
	"strconv"
	"strings"

	"github.com/obrel/monsturn/internal/data"
)

func ExtractData(ch string) (*data.ChannelData, error) {
	str := strings.Split(ch, "/")

	if str[2] == "" || str[4] == "" || str[6] == "" {
		return nil, errors.New("Invalid channel data.")
	}

	channel := &data.ChannelData{
		Realm:      str[2],
		User:       str[4],
		Allocation: str[6],
	}

	return channel, nil
}

func MessageParser(msg string) (*data.TurnData, error) {
	stat := &data.TurnData{}
	str := strings.Split(msg, ",")

	for _, dt := range str {
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

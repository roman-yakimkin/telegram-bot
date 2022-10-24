package utils

import (
	"regexp"
	"strconv"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
)

func ParseMsgText(msgText string, regExp *regexp.Regexp) (int, int, error) {
	var err error
	matches := regExp.FindAllStringSubmatch(msgText, -1)
	if len(matches) != 1 {
		return 0, 0, localerr.ErrIncorrectAmountValue
	}
	intPartStr, fracPartStr := matches[0][1], matches[0][3]
	if intPartStr == "" && fracPartStr == "" {
		return 0, 0, localerr.ErrIncorrectAmountValue
	}
	var intPart, fracPart int
	if intPartStr != "" {
		intPart, err = strconv.Atoi(intPartStr)
		if err != nil {
			return 0, 0, err
		}
	}
	if fracPartStr != "" {
		fracPart, err = strconv.Atoi(fracPartStr)
		if err != nil {
			return 0, 0, err
		}
	}
	return intPart, fracPart, nil
}

package ean

import (
	"strconv"
	"strings"
)

type Ean int64

func (ean Ean) ToString() string {
	eanStr := strconv.FormatInt(int64(ean), 10)

	// Pad with leading zeros if the length is less than 13
	if len(eanStr) < 13 {
		eanStr = strings.Repeat("0", 13-len(eanStr)) + eanStr
	}

	return eanStr
}

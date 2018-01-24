package seal_rtmp_conn

import "strings"

func handleParseStreamName(s string) (stream string, token string) {

	const TOKEN_STR = "?token="

	loc := strings.Index(s, TOKEN_STR)
	if loc < 0 {
		stream = s
	} else {
		stream = s[0:loc]
		token = s[loc+len(TOKEN_STR):]
	}

	return
}

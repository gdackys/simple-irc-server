package nick

import (
	"fmt"
	"regexp"
	s "simple-irc-server/server"
)

var pattern = regexp.MustCompile(`^[a-zA-Z\[\]\\` + "`" + `_^{|}][a-zA-Z0-9\[\]\\` + "`" + `_^{|}-]{0,8}$`)

func paramsParser(rawParams string) (s.CommandParams, error) {
	matches := pattern.FindStringSubmatch(rawParams)

	if matches == nil {
		return nil, fmt.Errorf("invalid params: %v", rawParams)
	}

	params := make(s.CommandParams)
	params["nickname"] = matches[0]

	return params, nil
}

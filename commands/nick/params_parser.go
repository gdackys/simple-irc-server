package nick

import (
	"fmt"
	"regexp"
)

var pattern = regexp.MustCompile(`^[a-zA-Z\[\]\\` + "`" + `_^{|}][a-zA-Z0-9\[\]\\` + "`" + `_^{|}-]{0,8}$`)

func ParamsParser(raw string) (*Params, error) {
	matches := pattern.FindStringSubmatch(raw)

	if matches == nil {
		return nil, fmt.Errorf("invalid params: %v", raw)
	}

	return &Params{matches[0]}, nil
}

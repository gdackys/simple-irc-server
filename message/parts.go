package message

import (
	"fmt"
	"strings"
)

type Parts struct {
	prefix  string
	user    string
	host    string
	command string
	params  string
}

func (p *Parts) String() string {
	parts := []string{
		fmt.Sprintf("prefix: \"%s\"", p.prefix),
		fmt.Sprintf("user: \"%s\"", p.user),
		fmt.Sprintf("host: \"%s\"", p.host),
		fmt.Sprintf("command: \"%s\"", p.command),
		fmt.Sprintf("params: \"%s\"", p.params),
	}

	return strings.Join(parts, ", ")
}

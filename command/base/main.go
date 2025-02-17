package base

type Base struct {
	name string
}

func New(name string) *Base {
	return &Base{name}
}

func (cmd *Base) GetName() string {
	return cmd.name
}

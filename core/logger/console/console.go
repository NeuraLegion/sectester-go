package console

import "os"

type Console struct {
	isColored bool
}

func New() *Console {
	return &Console{isColored: len(os.Getenv("NO_COLOR")) == 0}
}

func (c *Console) IsColored() bool {
	return c.isColored
}

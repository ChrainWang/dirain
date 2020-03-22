package pathmgr

import "fmt"

type ErrIllegalPath struct {
	r byte
	p int
}

func (self *ErrIllegalPath) Error() string {
	return fmt.Sprintf("[ILLEGAL PATH]Met illegal byte '%c' on position: %d", self.r, self.p)
}

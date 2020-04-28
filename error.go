package dirain

import "fmt"

type IllegalDirectoryErr struct {
	col int
}

func (self *IllegalDirectoryErr) Error() string {
	return fmt.Sprintf("Illegal rune at column %d.", self.col)
}

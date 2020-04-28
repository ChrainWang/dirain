package dirain

import "os"

type DirXpander struct {
	input []rune

	start int
	end   int

	output []rune
}

func NewDirXpander(input string) *DirXpander {
	return &DirXpander{
		input: []rune(input),
	}
}

func (self *DirXpander) getRune() (*rune, bool) {
	self.end++
	if self.end >= len(self.input) {
		return nil, false
	} else {
		return &self.input[self.end], true
	}
}

func (self *DirXpander) writeNode() {
	self.output = append(self.output, '/')
	self.output = append(self.output, self.input[self.start:self.end]...)
	self.start = self.end
}

func (self *DirXpander) rollback() error {

	for ptr := len(self.output) - 2; ptr >= 0; ptr-- {
		if self.output[ptr] == '/' {
			self.output = self.output[:ptr]
			return nil
		}
	}
	return self.illegalDirectory()
}

func (self *DirXpander) illegalDirectory() *IllegalDirectoryErr {
	return &IllegalDirectoryErr{self.end}
}

// when this function is called, it must meet the following conditions:
// 1. self.start == self.end
// 2. self.input[self.start] == '/' || self.input[self.start] == '\\'
func (self *DirXpander) xtractNode() error {
	r, more := self.getRune()
	if !more {
		return nil
	}

	switch *r {
	case '/', '\\':
		return self.xtractNode()
	case '.':
		self.start = self.end
		return self.dotPrefixed()
	default:
		self.start = self.end
		self.inNodeXtract()
		self.writeNode()
		return self.xtractNode()
	}

}

func (self *DirXpander) inNodeXtract() {
	if r, more := self.getRune(); more {
		switch *r {
		case '/', '\\':
			return
		default:
			self.inNodeXtract()
		}
	}
}

func (self *DirXpander) dotPrefixed() error {
	r, more := self.getRune()
	if !more {
		return nil
	}

	switch *r {
	case '/', '\\':
		return self.xtractNode()
	case '.':
		return self.doubleDotPrefixed()
	default:
		self.inNodeXtract()
		self.writeNode()
		return self.xtractNode()
	}

}

func (self *DirXpander) doubleDotPrefixed() error {
	r, more := self.getRune()
	if !more {
		return self.rollback()
	}

	switch *r {
	case '/', '\\':
		if err := self.rollback(); err == nil {
			return self.xtractNode()
		} else {
			return err
		}
	default:
		self.inNodeXtract()
		self.writeNode()
		return self.xtractNode()
	}
}

func (self *DirXpander) SetInput(p string) {
	self.start = 0
	self.end = 0
	self.input = []rune(p)
}

func (self *DirXpander) Expand() (string, error) {
	err := self.expand()
	if err == nil && len(self.output) == 0 {
		self.output = append(self.output, '/')
	}
	return string(self.output), err
}

func (self *DirXpander) expand() error {

	var err error
	if self.output == nil || len(self.output) == 0 {

		var parent string

		if self.input == nil || len(self.input) == 0 {
			if wd, err := os.Getwd(); err == nil {
				self.output = []rune(wd)
				return nil
			} else {
				return err
			}
		}

		switch {
		case self.input[0] == '/', self.input[0] == '\\':
			self.output = make([]rune, 0, len(self.input))
		case self.input[0] == '~' && (self.input[1] == '/' || self.input[1] == '\\'):
			parent, err = os.UserHomeDir()
			self.end = 1
		default:
			parent, err = os.Getwd()
			self.end = -1
		}

		if err != nil {
			return err
		}

		self.output = make([]rune, len(parent), len(parent)+len(self.input))
		copy(self.output, []rune(parent))

	} else if self.input[0] != '/' && self.input[0] != '\\' {
		self.end = -1
	}

	return self.xtractNode()
}

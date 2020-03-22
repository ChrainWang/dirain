package pathmgr

import (
	"unicode"
)

const (
	ILLEGAL_PATH_ERR_MSG string = "[ILLEGAL PATH]Illegal character %"
)

type PathExpander struct {
	rawPath  []byte
	start    int
	end      int
	rpLength int

	resultPath    []byte
	resultPointer int
}

func NewPathExpander(rawPath string) *PathExpander {
	return &PathExpander{rawPath: []byte(rawPath), rpLength: len(rawPath)}
}

func (self *PathExpander) nextRune() (rune, bool) {
	self.end++
	if self.end == self.rpLength {
		return rune(0), false
	} else {
		return rune(self.rawPath[self.end]), true
	}
}

func (self *PathExpander) writeNode() {
	node := self.rawPath[self.start:self.end]
	self.resultPath = append(self.resultPath, node...)
	self.resultPointer += (self.end - self.start)
	self.start = self.end
}

func (self *PathExpander) newIllegalPathErr() *ErrIllegalPath {
	return &ErrIllegalPath{r: self.rawPath[self.end], p: self.end}
}

func (self *PathExpander) checkPrefix() error {
	if r := rune(self.rawPath[self.start]); r == '/' {

		self.resultPath = make([]byte, 0, self.rpLength)
		return self.extractNodeBegin()

	} else if r == '~' {
		if homeDir, err := HomeDir(); err == nil {

			self.resultPath = make([]byte, 0, self.rpLength+len(homeDir))
			self.resultPath = append(self.resultPath, []byte(homeDir)...)
			self.resultPointer = len(self.resultPath)

			if r, hasNext := self.nextRune(); hasNext {

				if r == '/' {
					self.start = self.end
					return self.extractNodeBegin()
				} else {
					return self.newIllegalPathErr()
				}

			} else {
				// The given path is exactly '~'
				return nil
			}
		} else {
			return err
		}
	} else if r == '.' {
		// RELATIVE PATH begin with '.'
		if cwd, err := CWD(); err == nil {
			self.resultPath = make([]byte, 0, self.rpLength+len(cwd)+1)
			self.resultPath = append(self.resultPath, []byte(cwd)...)
			self.resultPath = append(self.resultPath, '/')
			self.resultPointer = len(self.resultPath)
		} else {
			return nil
		}
		return self.extractNodeBeginWithDot()
	} else if unicode.IsDigit(r) || unicode.IsLetter(r) {
		// RELATIVE PATH
		if cwd, err := CWD(); err == nil {
			self.resultPath = make([]byte, 0, self.rpLength+len(cwd)+1)
			self.resultPath = append(self.resultPath, []byte(cwd)...)
			self.resultPath = append(self.resultPath, '/')
			self.resultPointer = len(self.resultPath)
		} else {
			return nil
		}

		return self.extractNode()
	} else {
		return self.newIllegalPathErr()
	}
}

// When this function is called
// both of the poihters START and END are pointing to the same rune which is a '/'
func (self *PathExpander) extractNodeBegin() error {
	if r, hasNext := self.nextRune(); !hasNext {
		return nil
	} else if r == '/' {
		self.start = self.end
		return self.extractNodeBegin()
	} else if r == '.' {
		return self.extractNodeBeginWithDot()
	} else if unicode.IsDigit(r) || unicode.IsLetter(r) {
		return self.extractNode()
	} else {
		return self.newIllegalPathErr()
	}
}

func (self *PathExpander) extractNode() error {
	if r, haxNext := self.nextRune(); !haxNext {
		self.writeNode()
		return nil
	} else if r == '.' || r == '_' || unicode.IsLetter(r) || unicode.IsLetter(r) {
		return self.extractNode()
	} else if r == '/' {
		self.writeNode()
		return self.extractNodeBegin()
	} else {
		return self.newIllegalPathErr()
	}
}

func (self *PathExpander) extractNodeBeginWithDot() error {
	if r, hasNext := self.nextRune(); !hasNext {
		// The path is like: *****/.
		return nil
	} else if r == '.' {
		if r, hasNext = self.nextRune(); !hasNext {
			// The path is like: *****/..
			return self.rollback()
		} else if r == '/' {
			// The path is like: *****/../*****
			if err := self.rollback(); err == nil {
				self.start = self.end
				return self.extractNodeBegin()
			} else {
				return err
			}
		} else {
			return self.newIllegalPathErr()
		}
	} else if unicode.IsDigit(r) || unicode.IsLetter(r) {
		return self.extractNode()
	} else {
		return self.newIllegalPathErr()
	}
}

func (self *PathExpander) rollback() error {
	if self.resultPointer == 0 {
		return self.newIllegalPathErr()
	} else {
		self.resultPointer--
		if self.resultPath[self.resultPointer] == '/' {
			self.resultPath = self.resultPath[:self.resultPointer]
			return nil
		} else {
			return self.rollback()
		}
	}
}

func (self *PathExpander) Expand() (string, error) {
	if self.rpLength == 0 {
		return CWD()
	}

	if err := self.checkPrefix(); err != nil {
		return "", err
	}

	return string(self.resultPath[0:self.resultPointer]), nil
}

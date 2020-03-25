package pathmgr

import "testing"

func TestExpandPath(t *testing.T) {
	p := NewPathExpander("../fdsfs./f")
	if path, err := p.Expand(); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(path)
	}
}

package pathmgr

import "testing"

func TestExpandPath(t *testing.T) {
	p := NewPathExpander("a/../f")
	path, _ := p.Expand()
	t.Log(path)
}

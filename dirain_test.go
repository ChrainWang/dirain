package dirain

import "testing"

func TestAbsDirExpand(t *testing.T) {
	xpander := NewDirXpander("/abcdefg/../abcdd")

	if res, err := xpander.Expand(); err == nil {
		t.Log(res)
	} else {
		t.Error(err)
	}
}

func TestWDRelativeExpand(t *testing.T) {
	xpander := NewDirXpander("../..")

	if res, err := xpander.Expand(); err == nil {
		t.Log(res)
	} else {
		t.Error(err)
	}
}

func TestHomeRelativeExpand(t *testing.T) {
	xpander := &DirXpander{
		input: []rune("~/abcdefg/..abc/../abc"),
	}

	if res, err := xpander.Expand(); err == nil {
		t.Log(res)
	} else {
		t.Error(err)
	}
}

// func TestJoin(t *testing.T) {
// 	xpander := &DirXpander{
// 		input: []rune("~/abcdefg"),
// 	}
// 	xpander.expand()
// 	t.Log(string(xpander.output))
// 	xpander.SetInput("~/../abc")
// 	xpander.expand()
// 	t.Log(string(xpander.output))
// }

func TestAbs(t *testing.T) {
	pats := []string{
		"/a/b/c/d/e",
		"a/b/c/d/e",
		"./a/./b/./c/./d/./e/.",
		"a/../a/b/../b/c/../c/d/../d/e",
	}
	for _, p := range pats {
		if res, err := Abs(p); err == nil {
			t.Log(res)
		} else {
			t.Error(err)
		}
	}
}

func TestJoin(t *testing.T) {
	pats := []string{
		"~a/b/c/d/e",
		"a/b/c/d/e",
		"./a/./b/./c/./d/./e/.",
		"../a/../a/b/../b/c/../c/d/../d/e",
	}
	if res, err := Join(pats...); err == nil {
		t.Log(res)
	} else {
		t.Error(err)
	}
}

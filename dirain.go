package dirain

func Abs(path string) (string, error) {
	xpander := &DirXpander{
		input: []rune(path),
	}

	return xpander.Expand()
}

func Join(paths ...string) (string, error) {
	xpander := &DirXpander{}
	for len(paths) != 0 {
		xpander.SetInput(paths[0])
		if err := xpander.expand(); err != nil {
			return "", err
		}
		paths = paths[1:]
	}
	return string(xpander.output), nil
}

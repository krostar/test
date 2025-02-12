package ok

import "errors"

type firework struct{}

func launch(f *firework) error {
	return errors.New("boom")
}

func Awesome() error {
	f := new(firework)
	return launch(f)
}

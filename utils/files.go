package utils

import (
	"os"

	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/clipsync/errortypes"
)

func Exists(pth string) (exists bool, err error) {
	_, err = os.Stat(pth)
	if err == nil {
		exists = true
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	err = &errortypes.ReadError{
		errors.Wrapf(err, "utils: Failed to stat %s", pth),
	}
	return
}

func ExistsDir(pth string) (exists bool, err error) {
	stat, err := os.Stat(pth)
	if err == nil {
		exists = stat.IsDir()
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	err = &errortypes.ReadError{
		errors.Wrapf(err, "utils: Failed to stat %s", pth),
	}
	return
}

func ExistsFile(pth string) (exists bool, err error) {
	stat, err := os.Stat(pth)
	if err == nil {
		exists = !stat.IsDir()
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	err = &errortypes.ReadError{
		errors.Wrapf(err, "utils: Failed to stat %s", pth),
	}
	return
}

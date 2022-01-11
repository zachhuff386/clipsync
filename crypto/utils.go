package crypto

import (
	"crypto/rand"

	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/clipsync/errortypes"
)

func RandNonce() (nonce [24]byte, err error) {
	_, err = rand.Read(nonce[:])
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "crypto: Random nonce read error"),
		}
		return
	}

	return
}

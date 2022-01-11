package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/clipsync/config"
	"github.com/zachhuff386/youtube/errortypes"
	"golang.org/x/crypto/nacl/box"
)

var (
	ClientId         string
	PublicKey        [32]byte
	PrivateKey       [32]byte
	ClientPublicKeys = map[string][32]byte{}
)

func GenerateKey() (err error) {
	senderPubKey, senderPrivKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "crypto: Failed to generate nacl key"),
		}
		return
	}

	senderPubKeyStr := base64.RawStdEncoding.EncodeToString(senderPubKey[:])
	senderPrivKeyStr := base64.RawStdEncoding.EncodeToString(senderPrivKey[:])

	config.Config.PublicKey = senderPubKeyStr
	config.Config.PrivateKey = senderPrivKeyStr

	fmt.Printf("public_key=%s\n", senderPubKeyStr)

	return
}

func LoadKeys() (err error) {
	if config.Config.PublicKey == "" || config.Config.PrivateKey == "" {
		err = &errortypes.ReadError{
			errors.New("crypto: Missing public and private key"),
		}
		return
	}

	senderPubKeyByt, err := base64.RawStdEncoding.DecodeString(
		config.Config.PublicKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to parse public key"),
		}
		return
	}

	senderPrivKeyByt, err := base64.RawStdEncoding.DecodeString(
		config.Config.PrivateKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to parse private key"),
		}
		return
	}

	ClientId = config.Config.PublicKey
	copy(PublicKey[:], senderPubKeyByt)
	copy(PrivateKey[:], senderPrivKeyByt)

	for _, client := range config.Config.Clients {
		clientPubKeyByt, e := base64.RawStdEncoding.DecodeString(
			client.PublicKey)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrapf(
					e,
					"crypto: Failed to parse client public key '%s'",
					client.PublicKey,
				),
			}
			return
		}

		var clientPubKey [32]byte
		copy(clientPubKey[:], clientPubKeyByt)
		ClientPublicKeys[client.PublicKey] = clientPubKey
	}

	return
}

func Encrypt(clientId string, data []byte) (encData []byte, err error) {
	clientPubKey, ok := ClientPublicKeys[clientId]
	if !ok {
		err = &errortypes.ReadError{
			errors.Newf("crypto: Client key not found '%s'", clientId),
		}
		return
	}

	nonce, err := RandNonce()
	if err != nil {
		return
	}

	encDataByt := box.Seal(nonce[:], data, &nonce, &clientPubKey, &PrivateKey)

	encData = make([]byte, base64.RawStdEncoding.EncodedLen(len(encDataByt)))
	base64.RawStdEncoding.Encode(encData, encDataByt)

	return
}

func Decrypt(clientId string, encData []byte) (data []byte, err error) {
	clientPubKey, ok := ClientPublicKeys[clientId]
	if !ok {
		err = &errortypes.ReadError{
			errors.Newf("crypto: Client key not found '%s'", clientId),
		}
		return
	}

	encDataBox := make([]byte, base64.RawStdEncoding.DecodedLen(len(encData)))
	_, err = base64.RawStdEncoding.Decode(encDataBox, encData)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode client data"),
		}
		return
	}

	if len(encDataBox) < 30 {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Invalid client data"),
		}
		return
	}

	var nonce [24]byte
	copy(nonce[:], encDataBox[:24])
	encDataByt := encDataBox[24:]

	data, ok = box.Open([]byte{}, encDataByt, &nonce, &clientPubKey,
		&PrivateKey)
	if !ok {
		data = nil
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decrypt client data"),
		}
		return
	}

	return
}

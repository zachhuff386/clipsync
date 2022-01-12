package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/clipsync/errortypes"
	"github.com/zachhuff386/clipsync/utils"
)

var (
	Config = &ConfigData{}
	Path   = "./clipsync.conf"
)

type Client struct {
	Address   string `json:"address"`
	PublicKey string `json:"public_key"`
}

type ConfigData struct {
	Bind       string    `json:"bind"`
	PrivateKey string    `json:"private_key"`
	PublicKey  string    `json:"public_key"`
	Clients    []*Client `json:"clients"`
	loaded     bool      `json:"-"`
}

func (c *ConfigData) Save() (err error) {
	if !c.loaded {
		err = &errortypes.WriteError{
			errors.New("config: Config file has not been loaded"),
		}
		return
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File marshal error"),
		}
		return
	}

	err = ioutil.WriteFile(Path, data, 0600)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File write error"),
		}
		return
	}

	return
}

func Load() (err error) {
	data := &ConfigData{
		Bind:    ":9774",
		Clients: []*Client{},
	}

	exists, err := utils.Exists(Path)
	if err != nil {
		return
	}

	if !exists {
		data.loaded = true
		Config = data
		return
	}

	file, err := ioutil.ReadFile(Path)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File read error"),
		}
		return
	}

	err = json.Unmarshal(file, data)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File unmarshal error"),
		}
		return
	}

	data.loaded = true

	Config = data

	return
}

func Save() (err error) {
	err = Config.Save()
	if err != nil {
		return
	}

	return
}

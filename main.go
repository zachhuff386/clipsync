package main

import (
	"flag"
	"fmt"

	"github.com/zachhuff386/clipsync/config"
	"github.com/zachhuff386/clipsync/constants"
	"github.com/zachhuff386/clipsync/crypto"
	"github.com/zachhuff386/clipsync/server"
)

const help = `
Usage: clipsync COMMAND [CONFIG-FILE]

Commands:
  start        Start clipbord server
  generate-key Generate key for server
`

func GenerateKey() (err error) {
	err = config.Load()
	if err != nil {
		return
	}

	err = crypto.GenerateKey()
	if err != nil {
		return
	}

	err = config.Save()
	if err != nil {
		return
	}

	return
}

func Start() (err error) {
	err = config.Load()
	if err != nil {
		return
	}

	err = crypto.LoadKeys()
	if err != nil {
		return
	}

	server.Init()

	return
}

func main() {
	flag.Usage = func() {
		fmt.Println(help)
	}

	flag.Parse()

	switch flag.Arg(0) {
	case "version":
		fmt.Printf("pritunl-link v%s\n", constants.Version)
		break
	case "start":
		err := Start()
		if err != nil {
			panic(err)
		}
		break
	case "generate-key":
		err := GenerateKey()
		if err != nil {
			panic(err)
		}
		break
	default:
		fmt.Println(help)
	}
}

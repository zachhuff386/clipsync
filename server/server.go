package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/clipsync/clipboard"
	"github.com/zachhuff386/clipsync/config"
	"github.com/zachhuff386/clipsync/crypto"
	"github.com/zachhuff386/clipsync/errortypes"
	"github.com/zachhuff386/clipsync/utils"
)

var (
	lastSet    = time.Now()
	httpClient = &http.Client{
		Timeout: 15 * time.Second,
	}
)

func initHttp() {
	server := &http.Server{
		Addr:           config.Config.Bind,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 2048,
		Handler: http.HandlerFunc(func(
			w http.ResponseWriter, req *http.Request) {

			if req.Header.Get("User-Agent") != "clipsync" ||
				req.Method != "POST" ||
				!strings.HasPrefix(req.URL.Path, "/v1/") {

				utils.WriteText(w, 404, "Not Found")
				return
			}

			data, err := ioutil.ReadAll(req.Body)
			if err != nil {
				utils.WriteText(w, 500, "Read Error")
				return
			}
			_ = req.Body.Close()

			err = handleClipboardReq(req.URL.Path[4:], data)
			if err != nil {
				utils.LogError(err)
				utils.WriteText(w, 500, "Clipboard Error")
				return
			}

			utils.WriteText(w, 200, "Ok")
		}),
	}

	err := server.ListenAndServe()
	if err != nil {
		err = &errortypes.WriteError{
			errors.Newf("server: Failed to start http server"),
		}
		panic(err)
	}

	return
}

func handleClipboardReq(clientId string, encData []byte) (err error) {
	data, err := crypto.Decrypt(clientId, encData)
	if err != nil {
		return
	}

	clipboard.Set(string(data))

	return
}

func initWatch() {
	for {
		time.Sleep(10 * time.Millisecond)

		clipboard.Wait()
		if time.Since(lastSet) < 100*time.Millisecond {
			continue
		}
		lastSet = time.Now()

		data := clipboard.Get()
		if data == "" {
			continue
		}

		err := handleClipboardChange(data)
		if err != nil {
			utils.LogError(err)
			continue
		}
	}
}

func handleClipboardChange(dataStr string) (err error) {
	waiter := sync.WaitGroup{}
	data := []byte(dataStr)

	for _, client := range config.Config.Clients {
		waiter.Add(1)
		go func(clnt *config.Client) {
			defer waiter.Done()

			encData, e := crypto.Encrypt(clnt.PublicKey, data)
			if e != nil {
				err = e
				return
			}

			u := &url.URL{
				Scheme: "http",
				Host:   clnt.Address,
				Path:   fmt.Sprintf("/v1/%s", config.Config.PublicKey),
			}

			req, e := http.NewRequest(
				"POST",
				u.String(),
				bytes.NewBuffer(encData),
			)
			if e != nil {
				err = &errortypes.RequestError{
					errors.Wrap(e, "server: Client request error"),
				}
				return
			}

			req.Header.Set("User-Agent", "clipsync")

			res, e := httpClient.Do(req)
			if e != nil {
				err = &errortypes.RequestError{
					errors.Wrap(e, "server: Client response error"),
				}
				return
			}

			if res.StatusCode != 200 {
				err = &errortypes.RequestError{
					errors.Newf(
						"server: Client response status '%d'",
						res.StatusCode,
					),
				}
				return
			}
		}(client)
	}

	return
}

func Init() {
	go initHttp()
	initWatch()
}

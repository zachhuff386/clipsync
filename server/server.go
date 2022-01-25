package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/clipsync/clipboard"
	"github.com/zachhuff386/clipsync/config"
	"github.com/zachhuff386/clipsync/crypto"
	"github.com/zachhuff386/clipsync/errortypes"
	"github.com/zachhuff386/clipsync/utils"
)

var (
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

		data := clipboard.Get()
		if data == "" {
			continue
		}

		handleClipboardChange(data)
	}
}

func handleClipboardChange(dataStr string) {
	data := []byte(dataStr)

	for _, client := range config.Config.Clients {
		go func(clnt *config.Client) {
			err := sendClipboardChange(clnt, data)
			if err != nil {
				utils.LogError(err)
				return
			}
		}(client)
	}

	return
}

func sendClipboardChange(clnt *config.Client, data []byte) (err error) {
	encData, err := crypto.Encrypt(clnt.PublicKey, data)
	if err != nil {
		return
	}

	u := &url.URL{
		Scheme: "http",
		Host:   clnt.Address,
		Path:   fmt.Sprintf("/v1/%s", config.Config.PublicKey),
	}

	req, err := http.NewRequest(
		"POST",
		u.String(),
		bytes.NewBuffer(encData),
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "server: Client request error"),
		}
		return
	}

	req.Header.Set("User-Agent", "clipsync")

	res, err := httpClient.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "server: Client response error"),
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

	return
}

func Init() {
	go initHttp()
	initWatch()
}

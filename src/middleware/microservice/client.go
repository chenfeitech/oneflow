package microservice

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"middleware/jsonrpc_client"
	"math/rand"
	"net/http"
	"time"
	"utils/crypto_helper"
)

const (
	appID     = "data001"
	appSecret = "bc743a9a-0507-4c1b-ba47-7befccb260df"
	salt      = "!&&@@%#$!*^"
)

func Request(host string, method string, args interface{}, reply interface{}) error {
	timestamp := time.Now().Unix()
	random := fmt.Sprint(rand.Int())
	accessToken := md5.Sum([]byte(fmt.Sprintf("%v|%v|%v|%v|%v", appID, timestamp, appSecret, salt, random)))

	body, err := jsonrpc_client.EncodeClientRequest(method, args)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%v:46621/msc/api", host), &crypto_helper.DecryptoReader{ioutil.NopCloser(bytes.NewReader(body))})
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("App-Id", appID)
	req.Header.Add("Timestamp", fmt.Sprint(timestamp))
	req.Header.Add("Random", random)
	req.Header.Add("Access-Token", hex.EncodeToString(accessToken[:]))
	//req.Header.Add("X-Encrypto-Transport", "false")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		resp.Body = &crypto_helper.DecryptoReader{resp.Body}
		defer resp.Body.Close()
	}

	return jsonrpc_client.DecodeClientResponse(resp.Body, reply)
}

func SafeRequest(host string, method string, args interface{}, reply interface{}) error {
	timestamp := time.Now().Unix()
	random := fmt.Sprint(rand.Int())
	accessToken := md5.Sum([]byte(fmt.Sprintf("%v|%v|%v|%v|%v", appID, timestamp, appSecret, salt, random)))

	body, err := jsonrpc_client.EncodeClientRequest(method, args)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%v:46621/msc/sapi", host), &crypto_helper.DecryptoReader{ioutil.NopCloser(bytes.NewReader(body))})
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("App-Id", appID)
	req.Header.Add("Timestamp", fmt.Sprint(timestamp))
	req.Header.Add("Random", random)
	req.Header.Add("Access-Token", hex.EncodeToString(accessToken[:]))
	//req.Header.Add("X-Encrypto-Transport", "false")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		resp.Body = &crypto_helper.DecryptoReader{resp.Body}
		defer resp.Body.Close()
	}

	return jsonrpc_client.DecodeClientResponse(resp.Body, reply)
}

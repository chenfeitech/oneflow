package http_proxy

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type UrlInfo struct {
	HttpMethed int // 0:get 1:post
	ReqUrl     string
}

var PlantUrls map[int]UrlInfo

type ProxyRequest struct {
	PlantType int             `json:"plant_type"`
	Params    json.RawMessage `json:"params"`
}

type StringResult []byte

func (res StringResult) ToJson(w io.Writer) {
	w.Write([]byte(res))
}

func HttpGet(url string, param string) (res StringResult, err error) {
	get_url := url
	if len(param) > 0 {
		get_url = fmt.Sprintf("%s?%s", url, param)
	}
	resp, err := http.Get(get_url)
	if err != nil {
		return res, log.Error("http send error ", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, log.Error("http read error ", err)
	}

	log.Info("req url: ", get_url)
	log.Info("get ret: ", string(body))

	return StringResult(body), nil
}

func NewHttpGet(url string, param string) (StringResult, error) {
	var res StringResult
	var err error = nil
	var trys int = 0
	for trys = 0; trys < 3; trys++ {
		res, err = HttpGet(url, param)
		if err != nil {
			log.Error("request failed, trys=", trys)
			time.Sleep(1 * time.Second)
			continue
		} else {
			break
		}
	}
	log.Info("trys=", trys)
	return res, err
}

func HttpPost(url string, param string) (res StringResult, err error) {
	param_body := strings.NewReader(param)
	log.Info("req url: ", url, ", req: ", param)
	resp, err := http.Post(url, "application/json", param_body)
	if err != nil {
		return res, log.Error("http send error ", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, log.Error("http read error ", err)
	}
	log.Info("post ret: ", string(body))

	return StringResult(body), nil
}

func NewHttpPost(url string, param string) (StringResult, error) {
	var res StringResult
	var err error = nil
	var trys int = 0
	for trys = 0; trys < 3; trys++ {
		res, err = HttpPost(url, param)
		if err != nil {
			log.Error("request failed, trys=", trys)
			time.Sleep(1 * time.Second)
			continue
		} else {
			break
		}
	}
	log.Info("trys=", trys)
	return res, err
}

func WarphttpPost(url string, param string) (res StringResult, err error) {
	//param_body := strings.NewReader(strconv.Quote(param))
	//log.Info("eacape string=", param_body)
	param_body := strings.NewReader(param)
	log.Info("req url: ", url, " params=", param_body)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", param_body)
	if err != nil {
		return res, log.Error("http send error ", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, log.Error("http read error ", err)
	}
	log.Info("post ret: ", string(body))

	return StringResult(body), nil
}

func NewWarphttpPost(url string, param string) (StringResult, error) {
	//	var res StringResult
	//	var err error = nil
	//	var trys int = 0
	//	for trys = 0; trys < 3; trys++ {
	//		res, err = WarphttpPost(url, param)
	//		if err != nil {
	//			time.Sleep(1 * time.Second)
	//			log.Error("request failed, trys=", trys)
	//			continue
	//		} else {
	//			break
	//		}
	//	}
	//	log.Info("trys=", trys)
	//	return res, err
	return NewWarphttpPost4(url, param, 3, 1*time.Second)
}

func NewWarphttpPost4(url string, param string, retry int, interval time.Duration) (StringResult, error) {
	var res StringResult
	var err error = nil
	var trys int = 0
	for trys = 0; trys < retry; trys++ {
		res, err = WarphttpPost(url, param)
		if err != nil {
			time.Sleep(interval)
			log.Error("request failed, trys=", trys)
			continue
		} else {
			break
		}
	}
	log.Info("trys=", trys)
	return res, err
}

type HttpStatus struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
}

func WarphttpPost_V2(url string, param string) (res StringResult, http_status HttpStatus, err error) {
	param_body := strings.NewReader(param)
	log.Info("req url: ", url, " params=", param_body)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", param_body)
	if err != nil {
		return res, http_status, log.Error("http send error, err=", err)
	}
	http_status.Status = resp.Status
	http_status.StatusCode = resp.StatusCode
	http_status.Proto = resp.Proto

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, http_status, log.Error("http read error ", err)
	}
	log.Info("post ret: ", string(body))

	return StringResult(body), http_status, nil
}

func NewWarphttpPost_V2(url string, param string, retry int, interval time.Duration) (StringResult, error) {
	var res StringResult
	var err error = nil
	var trys int = 0
	http_status := HttpStatus{}
	for trys = 0; trys < retry; trys++ {
		res, http_status, err = WarphttpPost_V2(url, param)
		if err != nil {
			time.Sleep(interval)
			log.Error("request failed, trys=", trys)
			continue
		} else {
			if http_status.StatusCode != 200 {
				time.Sleep(interval)
				log.Error("request failed, trys=", trys, ", http_code=", http_status.StatusCode)
				continue
			} else {
				break
			}
		}
	}
	log.Info("trys=", trys)
	return res, err
}

package asr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	//AsrURL 语音识别接口的服务器地址
	AsrURL = "http://vop.baidu.com/server_api"
)

type asrResult struct {
	ErrNo  int      `json:"err_no"`
	ErrMsg string   `json:"err_msg"`
	SN     string   `json:"sn"`
	Result []string `json:"result"`
}

type asrRequest struct {
	Format  string `json:"format"`
	Rate    int    `json:"rate"`
	Channel int    `json:"channel"`
	Token   string `json:"token"`
	Cuid    string `json:"cuid"`
	Len     int    `json:"len"`
	Speech  string `json:"speech"`
}

//ToText 调用接口获取语音对应的文字
func ToText(token string, reader io.Reader) (string, error) {
	requestData := asrRequest{}
	requestData.Format = "wav"
	requestData.Rate = 8000
	requestData.Channel = 1
	requestData.Cuid = "baidu-yuyin-for-go"
	requestData.Token = token
	speechBody, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	speech := base64.StdEncoding.EncodeToString(speechBody)
	requestData.Speech = speech
	requestData.Len = len(speechBody)

	requestJSON, _ := json.Marshal(requestData)

	request, err := http.NewRequest("POST", AsrURL, bytes.NewReader(requestJSON))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Length", fmt.Sprintf("%d", len(requestJSON)))
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resultBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result asrResult
	err = json.Unmarshal(resultBody, &result)
	if err != nil {
		return "", err
	}
	if result.ErrNo > 0 {
		return "", fmt.Errorf("error - %d - %s", result.ErrNo, result.ErrMsg)
	}
	return result.Result[0], nil
}

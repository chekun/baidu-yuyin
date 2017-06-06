package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	baiduTokenURL = "https://openapi.baidu.com/oauth/2.0/token"

	grantType = "client_credentials"
)

//Oauth Oauth对象
type Oauth struct {
	clientID     string
	clientSecret string
	cacheMan     CacheMan
}

//New 创建Oauth请求对象
func New(clientID, clientSecret string, cache CacheMan) *Oauth {
	return &Oauth{
		clientID:     clientID,
		clientSecret: clientSecret,
		cacheMan:     cache,
	}
}

//GetToken 获取AccessToken值
func (oauth *Oauth) GetToken() (string, error) {
	if oauth.cacheMan != nil && oauth.cacheMan.IsValid() {
		return oauth.cacheMan.Get()
	}

	//fetch token from baidu
	params := url.Values{}
	params.Add("grant_type", grantType)
	params.Add("client_id", oauth.clientID)
	params.Add("client_secret", oauth.clientSecret)

	resp, err := http.PostForm(baiduTokenURL, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if bytes.Contains(body, []byte("error")) {
		return "", errors.New("failed to retrive access_token, wrong client")
	}
	var tokenResult struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	err = json.Unmarshal(body, &tokenResult)
	if err != nil {
		return "", err
	}
	if oauth.cacheMan != nil {
		err = oauth.cacheMan.Set(tokenResult.AccessToken, tokenResult.ExpiresIn)
		if err != nil {
			return "", err
		}
	}
	return tokenResult.AccessToken, nil
}

百度语音识别SDK for Go
===============================

第一个可用版本，随心使用，注意安全 😄

## 安装

```
go get -u github.com/chekun/baidu-yuyin
```

## 使用方法：

```go
import (
	"fmt"
	"os"

	"github.com/chekun/baidu-yuyin/asr"
	"github.com/chekun/baidu-yuyin/oauth"
)

clientID := "your-client-id"
clientSecret := "your-client-secret"

auth := oauth.New(clientID, clientSecret, oauth.NewMemoryCacheMan())
//一次性使用也可以不缓存token, 如下
//auth := oauth.New(clientID, clientSecret, nil)
//也可以实现自己的缓存，往下看⬇️ ️ ️
token, err := auth.GetToken()
if err != nil {
    panic(err)
}

file, err := os.Open("speech.wav")
if err != nil {
    panic(err)
}
defer file.Close()
fmt.Println(asr.ToText(token, file))
```

你也可以用实现自己的`token缓存`, 实现这个`oauth.CacheMan`接口即可

```go
type CacheMan interface {
	Get() (string, error)
	Set(string, int) error
	IsValid() bool
}
```

最后在 `oauth.New` 的时候传入。


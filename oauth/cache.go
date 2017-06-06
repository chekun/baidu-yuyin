package oauth

import "time"

//CacheMan 缓存控制接口
type CacheMan interface {
	Get() (string, error)
	Set(string, int) error
	IsValid() bool
}

//MemoryCacheMan 内存缓存
type MemoryCacheMan struct {
	token     string
	expiresAt time.Time
}

//NewMemoryCacheMan 创建内存缓存对象
func NewMemoryCacheMan() *MemoryCacheMan {
	return &MemoryCacheMan{
		token:     "",
		expiresAt: time.Now(),
	}
}

//Get 获取token
func (cache *MemoryCacheMan) Get() (string, error) {
	return cache.token, nil
}

//Set 设置token
func (cache *MemoryCacheMan) Set(token string, expiresIn int) error {
	cache.token = token
	cache.expiresAt = time.Now().Add(time.Second * time.Duration(expiresIn))
	return nil
}

//IsValid 检测token是否有效
func (cache *MemoryCacheMan) IsValid() bool {
	sub := time.Now().Sub(cache.expiresAt)
	return sub.Seconds() > 10
}

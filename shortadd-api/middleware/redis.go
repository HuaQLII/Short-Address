package middleware

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pilu/go-base62"
	"github/im-lauson/Short-Address/global"
	"time"
)

// ===============================================================================================
// = redis key值的定义
const (
	// UrlIdKey 全局自增器，每次自增一，保证永远不会重复
	UrlIdKey = "next.url.id"
	// ShortLinkKey 映射短地址和长地址之间的关系
	ShortLinkKey = "shortlink:%s:url"
	// URLHashKey 映射长地址的hash值到段地址之间的关系
	UrlHashKey = "urlhash:%s:url"
	// ShortLinkDetailKey 段地址详细信息的Key，映射长地址和短地址地址的详细信息
	ShortLinkDetailKey = "shortlink:%s:detail"
)

// ===============================================================================================
// = Redis客户端的数据结构
type RedisCli struct {
	Cli *redis.Client
}

// ===============================================================================================
// = env.go最后r报错点击自动生成的
//func (r *RedisCli) UnShorten(eid string) (string, error) {
//	//TODO implement me
//	panic("implement me")
//}

// ===============================================================================================
// = Url详细信息的结构体

type URLDetail struct {
	URL                 string        `json:"url"`
	CreateAt            string        `json:"createAt"`
	ExpirationInMinutes time.Duration `json:"expirationInMinutes"`
}

// NewRedisCli ===============================================================================================
// = 初始化redis客户端的结构体
func NewRedisCli(addr string, passwd string, db int) *RedisCli {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db})
	if _, err := c.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	return &RedisCli{Cli: c}
}

// ===============================================================================================
// = 把长地址转换为一个短地址

func (r *RedisCli) Shorten(url string, exp int64) (string, error) {

	// ! 把url转化成啥1的哈希值
	h := toSha1(url)
	// ! 到缓存中看一下这个url值是否已经转化过了
	d, err := r.Cli.Get(context.Background(), fmt.Sprintf(UrlHashKey, h)).Result()
	if err == redis.Nil {
		//说明这个值不存在
	} else {
		if d == "{}" {

		} else {
			return d, nil
		}
	}
	err = r.Cli.Incr(context.Background(), UrlIdKey).Err()
	if err != nil {
		return "", err
	}

	// ! 全局变量转为base62的短地址
	id, err := r.Cli.Get(context.Background(), UrlIdKey).Uint64()
	if err != nil {
		return "", err
	}
	eid := base62.Encode(int(id))

	err = r.Cli.Set(context.Background(), fmt.Sprintf(ShortLinkKey, eid),
		url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}
	err = r.Cli.Set(context.Background(), fmt.Sprintf(UrlHashKey, h), eid,
		time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}
	detail, err := json.Marshal(
		&URLDetail{
			URL:                 url,
			CreateAt:            time.Now().String(),
			ExpirationInMinutes: time.Duration(exp),
		},
	)
	if err != nil {
		return "", err
	}
	// ! 把短地址格式化到段地址的Key中
	err = r.Cli.Set(context.Background(), fmt.Sprintf(ShortLinkDetailKey, eid),
		detail, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}
	return eid, nil
}

// ===============================================================================================
// = 短地址返回的信息
func (r *RedisCli) ShortLinkInfo(eid string) (interface{}, error) {
	d, err := r.Cli.Get(context.Background(), fmt.Sprintf(ShortLinkDetailKey, eid)).Result()
	if err == redis.Nil {
		return "", global.StatusError{Code: 404, Err: fmt.Errorf("Unknow short URL ")}
	} else if err != nil {
		return "", err
	} else {
		return d, nil
	}
}

// ===============================================================================================
// = 短地址转换为一个长地址
func (r *RedisCli) UnShorten(eid string) (string, error) {

	url, err := r.Cli.Get(context.Background(), fmt.Sprintf(ShortLinkKey, eid)).Result()
	if err == redis.Nil {
		return "", global.StatusError{Code: 404, Err: fmt.Errorf("Unknow short URL ")}
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}

func toSha1(str string) string {
	var (
		sha = sha1.New()
	)
	return string(sha.Sum([]byte(str)))
}

package services

import (
	"goprizm/log"
	"goprizm/sysutils"
	"net/url"
	"strings"

	redis "github.com/go-redis/redis"
)

func Redis() redis.UniversalClient {
	return RedisEnv("")
}

func RedisEnv(env string) redis.UniversalClient {
	if env == "" {
		env = "REDIS"
	}

	rUrls := strings.Split(sysutils.Getenv(env, "redis://localhost:6379/0"), ",")
	var (
		hosts    []string
		password string
	)
	for _, rUrl := range rUrls {
		u, err := url.Parse(rUrl)
		if err != nil {
			log.Fatalf("Failed to parse redis urls:%s from env:%s err:%v", rUrls, env, err)
		}

		if u.User != nil {
			if p, _ := u.User.Password(); p != "" {
				password = p
			}
		}
		hosts = append(hosts, u.Host)
	}

	return redis.NewUniversalClient(&redis.UniversalOptions{Addrs: hosts, Password: password})
}

package utils

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"encoding/json"
	"goprizm/log"
	"goprizm/sysutils"

	redis "gopkg.in/redis.v5"
)

const (
	redisprefix = "ADMIN:"
	// ActiveFilterKey used to store Active filters
	ActiveFilterKey = "ActiveFilters"
	// TimeRangeFilterKey used to store selected Time Range
	TimeRangeFilterKey = "TimeRangeFilter"
	// UserPrefCacheInfoKey used to save users pref
	UserPrefCacheInfoKey = "UserPrefCacheInfo"
	// TimeRangeFilterColumnNameKey - UserPrefCacheInfo ->
	// TimeRangeFilterColumnName for added at/updated at
	TimeRangeFilterColumnNameKey = "TimeRangeFilterColumnName"
)

var (
	client *redis.Client
)

func init() {
	client = Redis()
}

//Put value into Redis cache
func Put(s *model.SessionContext, key string, value interface{}) {
	content, err := json.Marshal(value)
	if err != nil {
		logutil.Errorf(s, "Redis Put Error %s", err.Error())
		return
	}
	client.Set(getRedisKey(s, key), string(content), 0)
}

//Get value from Redis cache
func Get(s *model.SessionContext, key string) string {
	redisKey := getRedisKey(s, key)
	val, err := client.Get(redisKey).Result()
	if err == redis.Nil {
		logutil.Errorf(s, "Redis Key Not Found %s", redisKey)
		return ""
	}
	return val
}

func getRedisKey(s *model.SessionContext, key string) string {
	return redisprefix + s.User.TenantId + ":" + s.User.UserName + ":" + key
}

//Redis client
func Redis() *redis.Client {
	return RedisEnv("")
}

//RedisEnv client
func RedisEnv(env string) *redis.Client {
	opts := RedisOptions(env)
	c := redis.NewClient(opts)
	logutil.Printf(nil, "Redis Client connected on %s ", opts.Addr)
	pingTest(c)
	return c
}

//RedisOptions option
func RedisOptions(env string) *redis.Options {
	if env == "" {
		env = "REDIS"
	}
	opts, err := redis.ParseURL(sysutils.Getenv(env, "redis://localhost:6379/0"))
	if err != nil {
		log.Fatalf("Failed to parse redis urls from env:%s err:%v", env, err)
	}
	return opts
}

func pingTest(c *redis.Client) {
	pong, err := c.Ping().Result()
	logutil.Debugf(nil, "Redis Client ping test ping:%s, error:%s", pong, err)
}

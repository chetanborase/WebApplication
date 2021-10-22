package cache

import (
	"WebApplication/env"
	"fmt"
)

func redisUrl() string {
	return fmt.Sprintf("%s:%s",
		env.Get(env.CacheRedisAddress, "nil"),
		env.Get(env.CacheRedisPort, "nil"))
}
func redisPass() string {
	return fmt.Sprintf("%s",
		env.Get(env.CacheRedisPassword, "nil"))
}

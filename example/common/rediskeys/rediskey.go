package rediskeys

type RedisKey string

func (key RedisKey) WithParams(params ...string) string {
	if len(params) == 0 {
		return string(key)
	}
	k := string(key)
	for _, v := range params {
		k += ":" + v
	}
	return k
}

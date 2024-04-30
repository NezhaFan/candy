package redis

type hyperloglog struct {
	base
}

func NewHyperLogLog(key string) hyperloglog {
	return hyperloglog{base{key}}
}

// 添加，至少有一个添加成功返回1，否则返回0
func (h hyperloglog) PFAdd(vals ...any) (bool, error) {
	n, err := rdb.PFAdd(ctx, h.key, vals...).Result()
	return n == 1, err
}

// 统计数量，是存在0.81%误差的
func (h hyperloglog) PFCount() int64 {
	return rdb.PFCount(ctx, h.key).Val()
}

// 合并其他的
func (h hyperloglog) PFMerge(mergekeys ...string) (bool, error) {
	r, err := rdb.PFMerge(ctx, h.key, mergekeys...).Result()
	return r == OK, err
}

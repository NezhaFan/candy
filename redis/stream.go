package redis

import (
	"github.com/redis/go-redis/v9"
)

type XMessage = redis.XMessage

// stream 和 group 和 consumer 的关系为  m:1:n
// 但是实际中1个group对应多个stream处理比较麻烦，所以这里封装简化为 1:1:n

// 频道，生产者
type stream struct {
	base
}

func NewStream(key string) stream {
	return stream{base{key}}
}

// values推荐map[string]any形式，也可以[]any
func (s stream) XAdd(id string, values any) (string, error) {
	if id == "" {
		id = "*"
	}
	args := &redis.XAddArgs{
		Stream: s.key,
		ID:     id,
		Values: values,
	}
	return rdb.XAdd(ctx, args).Result()
}

// 长度
func (s stream) XLen() int64 {
	return rdb.XLen(ctx, s.key).Val()
}

// 删除
func (s stream) XDel(ids ...string) int64 {
	return rdb.XDel(ctx, s.key, ids...).Val()
}

// 分组
type group struct {
	name   string
	stream string
}

// 一个groupname内消息具有唯一性。 两个groupname的组可以读到同一条消息。
// 0表示从前向后读，$表示从后向前读。 MKSTREAM表示如果stream不存在则创建。  重复创建group会报错，忽略。
func NewGroup(groupname string, stream stream) (g group) {
	g.name = groupname
	g.stream = stream.key
	rdb.XGroupCreateMkStream(ctx, stream.key, groupname, "0")
	return
}

func (g group) Destroy() error {
	return rdb.XGroupDestroy(ctx, g.stream, g.name).Err()
}

type consumer struct {
	args *redis.XReadGroupArgs
}

// count表示从每个stream每次读取数量
func NewConsumer(name string, group group, count int64) (c consumer) {
	c.args = &redis.XReadGroupArgs{
		Consumer: name,
		Group:    group.name,
		Streams:  []string{group.stream, ">"}, // 只给未被读取的消息
		Count:    count,
		Block:    0,
		NoAck:    false, // false表示需要主动Ack确认
	}
	return
}

// 读取
func (c consumer) Read() ([]XMessage, error) {
	arr, err := rdb.XReadGroup(ctx, c.args).Result()
	if len(arr) == 0 || err != nil {
		return nil, err
	}
	return arr[0].Messages, nil
}

// 确认
func (c consumer) Ack(ids ...string) error {
	return rdb.XAck(ctx, c.args.Streams[0], c.args.Group, ids...).Err()
}

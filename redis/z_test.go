package redis

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

type User struct {
	Id   int    `redis:"id"`
	Name string `redis:"name"`
}

func init() {
	if err := Connect("127.0.0.1:6379", "123456"); err != nil {
		panic(err)
	}
}

func TestString(t *testing.T) {
	str := NewString("test-string")
	err := str.Set("123", time.Hour)
	fmt.Println("set 123:", err)
	b, err := str.SetNX("234", time.Hour)
	fmt.Println("setnx 234:", b, err)
	v, err := str.Get()
	fmt.Println("get:", v, err)
	ttl := str.TTL()
	fmt.Println("ttl:", ttl)
	ok, err := str.Expire(time.Second)
	fmt.Println("expire:", ok, err)
	ttl = str.TTL()
	fmt.Println("ttl:", ttl)
	ok = str.Del()
	fmt.Println("del:", ok)
	v, err = str.Get()
	fmt.Println("get:", `"`+v+`"`, err)
}

func BenchmarkString(b *testing.B) {
	str := NewString("test-string")
	defer str.Del()

	wg := &sync.WaitGroup{}
	wg.Add(b.N)
	fmt.Println("测试次数：", b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			ok, err := str.SetNX("1", time.Minute)
			if err != nil {
				panic(err)
			}
			if ok {
				fmt.Println("设置成功")
			}
		}()
	}
	wg.Wait()
}

func TestBitmap(t *testing.T) {
	bitmap := NewBitmap("test-bitmap")
	defer bitmap.Del()

	const year = 366
	// 假设用户id分别为0、1、10000
	var id0, id1, id10000 int64 = 0, 1, 10000

	// 用户1三天
	bitmap.SetBit(id0*year+1, 1)
	bitmap.SetBit(id0*year+365, 1)
	bitmap.SetBit(id0*year+366, 1)
	// 用户1二天
	bitmap.SetBit(id1*year+1, 1)
	bitmap.SetBit(id1*year+365, 1)
	// 用户10000二天
	bitmap.SetBit(id10000*year+1, 1)
	bitmap.SetBit(id10000*year+365, 1)

	n := bitmap.GetBit(id10000*year + 365)
	fmt.Println("getbit", id10000*year+365, ":", n)

	n = bitmap.BitCount(id0*year+1, id0*year+year)
	fmt.Println("用户1登录天数:", n)

	n = bitmap.BitCount(id1*year+1, id1*year+year)
	fmt.Println("用户2登录天数:", n)

	n = bitmap.BitCount(id10000*year+1, id10000*year+year)
	fmt.Println("用户10000登录天数:", n)

	fmt.Println("总计数：", bitmap.BitCount(0, -1))
}

func TestHash(t *testing.T) {
	hash := NewHash("test-hash")

	hash1 := hash.SubKey("1")
	defer hash1.Del()

	err := hash1.HSet("name", "Alice")
	fmt.Println("hset name Alice:", err)
	v, err := hash1.HGet("name")
	fmt.Println("hget name:", v, err)

	hash2 := hash.SubKey("2")
	defer hash2.Del()

	var u1, u2 User
	u1.Id = 2
	u1.Name = "Bob"
	err = hash2.HMSet(&u1)
	fmt.Println("hmset id:", 2, "name:", "Bob")

	err = hash2.HMGet(&u2, "name")
	fmt.Println("hmget name:", u2, err)
	err = hash2.HGetAll(&u2)
	fmt.Println("hgetall:", u2, err)

	i, err := hash2.HIncrBy("id", 10)
	fmt.Println("hincrby id 10:", i, err)
}

func TestHyperLogLog(t *testing.T) {
	hyperloglog1 := NewHyperLogLog("test-hyperloglog2")
	defer hyperloglog1.Del()
	hyperloglog2 := NewHyperLogLog("test-hyperloglog2")
	defer hyperloglog2.Del()

	n := 1000
	es := make([]any, 0, n)
	for i := 0; i < n; i++ {
		es = append(es, i)
	}

	b, err := hyperloglog1.PFAdd(es...)
	fmt.Println("pfadd:", b, err)

	size := hyperloglog1.PFCount()
	fmt.Println("pfcount:", size)

	hyperloglog2.PFAdd(1, 2, 3, 10002, 10003)
	b, err = hyperloglog1.PFMerge("test-hyperloglog2")
	fmt.Println("pfmerge:", b, err)

	size = hyperloglog1.PFCount()
	fmt.Println("pfcount:", size)
}

func TestList(t *testing.T) {
	list := NewList("test-list")
	defer list.Del()

	_, err := list.LPush(2, 1)
	fmt.Println("lpush 2 1:", err)
	_, err = list.RPush(3, 3, 4)
	fmt.Println("rpush 3 3 4:", err)

	n := list.LLen()
	fmt.Println("llen:", n)

	r := list.LRange(0, -1)
	fmt.Println("lrange:", r)

	v, err := list.LPop()
	fmt.Println("lpop:", v, err)

	v, err = list.RPop()
	fmt.Println("rpop:", v, err)

	n, err = list.LRem("2", 1)
	fmt.Println("lrem 4:", n, err)

	n, err = list.Rem("3")
	fmt.Println("rem 3:", n, err)

}

func TestSet(t *testing.T) {
	set := NewSet("test-set")
	defer set.Del()

	n, err := set.SAdd(1, 2, 3, 1)
	fmt.Println("sadd 1 2 3 1:", n, err)

	n = set.SCard()
	fmt.Println("scard:", n)

	b := set.SIsMember("2")
	fmt.Println("sismember 2:", b)

	n, err = set.SRem(2, 3)
	fmt.Println("srem 2 3:", n, err)

	r := set.SMembers()
	fmt.Println("smsmbers:", r)
}

func TestZSet(t *testing.T) {
	zset := NewZSet("test-zset")
	zset.Expire(time.Minute)

	n, err := zset.ZAdd(
		Z{Score: 0.5, Member: "A"},
		Z{Score: 1, Member: "C"},
		Z{Score: 0.5, Member: "B"},
		Z{Score: 1, Member: "D"},
	)
	fmt.Println("添加ABCD:", n, err)

	members := zset.ZRange(0, -1)
	fmt.Println("正序获取全部成员:", members)

	r := zset.ZRangeWithScores(0, -1)
	fmt.Println("正序获取全部成员和积分:", r)

	n = zset.ZCard()
	fmt.Println("成员个数:", n)

	score := zset.ZScore("B")
	fmt.Println("B的积分:", score)

	n, err = zset.ZRem("B")
	fmt.Println("移除B:", n, err)

	n = zset.ZLexCount("(A", "[C")
	fmt.Println("获取(A,C]之间成员个数:", n)

	members = zset.ZRevRange(0, -1)
	fmt.Println("逆序获取全部成员:", members)

	r = zset.ZRevRangeWithScores(0, -1)
	fmt.Println("逆序获取全部成员和积分:", r)

	r = zset.ZRangeByScore("1", MaxInf, 1)
	fmt.Println("获取积分>=1的1个成员和积分:", r)

	n, err = zset.ZRemRangeByRank(-2, -2)
	fmt.Println("删除排名倒数第2的成员:", n, err)

	n = zset.ZRank("D")
	fmt.Println("获取D从小到大的排名:", n)

	n = zset.ZRevRank("D")
	fmt.Println("获取D从大到小的排名:", n)

	// 全部删除
	n, err = zset.ZRemRangeByScore(MinInf, MaxInf)
	fmt.Println("删除积分为-inf到+inf之间的成员:", n, err)

}

func BenchmarkMutex(b *testing.B) {
	wg := &sync.WaitGroup{}
	wg.Add(b.N)
	start := time.Now()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			defer wg.Done()
			mutex := NewMutex(context.Background(), "mutex")
			if err := mutex.Lock(); err != nil {
				panic(err)
			}
			defer mutex.UnLock()
			time.Sleep(time.Millisecond * 2)
		}(i)
	}

	wg.Wait()
	fmt.Println("尝试次数", b.N, "耗时", time.Since(start).Milliseconds(), "ms")
}

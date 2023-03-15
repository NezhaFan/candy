package candy

import (
	"sync"
	"time"
)

type Map[KT comparable, VT any] struct {
	// 键值过期时间(秒)
	expireSecond uint32
	// 下次清理时间时间戳(在惰性删除时判断)
	// nextCleanTime uint32
	// 数据
	data sync.Map
}

type mapItem[VT any] struct {
	expireAt int64
	value    VT
}

// 并发安全、可设置过期时间的map 。  (惰性删除+定期删除)
func (m *Map[KT, VT]) Expires(t time.Duration) {
	m.expireSecond = uint32(t.Seconds())
}

func (m *Map[KT, VT]) Set(key KT, value VT) {
	if m.expireSecond == 0 {
		m.data.Store(key, value)
		return
	}
	item := mapItem[VT]{
		expireAt: time.Now().Unix() + int64(m.expireSecond),
		value:    value,
	}
	m.data.Store(key, item)
}

func (m *Map[KT, VT]) Get(key KT) (VT, bool) {
	v, ok := m.data.Load(key)
	// 不存在或者无过期
	if !ok || m.expireSecond == 0 {
		return v.(VT), false
	}

	now := time.Now().Unix()
	item := v.(mapItem[VT])

	// 未过期
	if now < item.expireAt {
		return item.value, true
	}

	// 已过期
	m.data.Delete(key)

	// 被动定期清理
	// if now > int64(m.nextCleanTime) {
	// 	go m.clean(now)
	// }

	var result VT
	return result, false
}

func (m *Map[KT, VT]) Has(key KT) bool {
	_, ok := m.Get(key)
	return ok
}

func (m *Map[KT, VT]) Del(key KT) {
	m.data.Delete(key)
}

// func (m *Map[KT, VT]) clean(now int64) {
// 	if m.expireSecond >= 600 {
// 		m.nextCleanTime = uint32(now) + m.expireSecond
// 	} else {
// 		m.nextCleanTime = uint32(now) + 600
// 	}
// 	m.data.Range(func(key, value any) bool {
// 		item := value.(mapItem[VT])
// 		if now >= item.expireAt {
// 			m.data.Delete(key)
// 		}
// 		return false
// 	})
// }

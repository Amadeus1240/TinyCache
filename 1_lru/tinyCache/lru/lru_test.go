package lru

import (
	"reflect"
	"testing"
)

type myString string

func (m myString) Len() int {
	return len(m)
}

// 测试 Get 方法
func Test_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", myString("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(myString)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

// 测试，当使用内存超过了设定值时，是否会触发“无用”节点的移除
func Test_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	Cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(Cap), nil)
	lru.Add(k1, myString(v1))
	lru.Add(k2, myString(v2))
	lru.Add(k3, myString(v3))
	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("RemoveOldest key1 failed")
	}
}

func Test_OnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("k1", myString("v1"))
	lru.Add("k2", myString("v2"))
	lru.Add("k3", myString("v3"))
	lru.Add("k4", myString("v4"))

	// callback 是某条记录被移除时的回调函数
	// 由于记录 k1 和 k2 被移除了，所以 k1 和 k2 被加入 keys 中
	expect := []string{"k1", "k2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

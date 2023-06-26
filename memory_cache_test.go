package memorycache

import (
	"testing"
	"time"
)

const (
	key1   = "test_Key1"
	key2   = "test_Key2"
	key3   = "test_Key3"
	value1 = "test_Value1"
	value2 = "test_Value2"
	value3 = "test_Value3"
)

func TestGet(t *testing.T) {
	cache := New(MemoryCacheOptions{})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
	})

	if _, ok := cache.Get(key1); !ok {
		t.Fatalf("incorrect result: expected true, got %t", ok)
	}

	if result, _ := cache.Get(key1); result.(string) != value1 {
		t.Fatalf("incorrect result: expected %s, got %s", value1, result.(string))
	}
}

func TestSet(t *testing.T) {
	cache := New(MemoryCacheOptions{})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
	})

	if _, ok := cache.Get(key1); !ok {
		t.Fatalf("incorrect result: expected true, got %t", ok)
	}

	cache.Set(key1, value2, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
	})

	if result, _ := cache.Get(key1); result.(string) != value2 {
		t.Fatalf("incorrect result: expected %s, got %s", value2, result.(string))
	}
}

func TestDel(t *testing.T) {
	cache := New(MemoryCacheOptions{})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
	})

	cache.Del(key1)

	if _, ok := cache.Get(key1); ok {
		t.Errorf("incorrect result: expected false, got %t", ok)
	}
}

func TestReset(t *testing.T) {
	cache := New(MemoryCacheOptions{})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Second * 1),
	})

	cache.Reset()

	if count := cache.Count(); count != 0 {
		t.Fatalf("incorrect result: expected 0, got %d", count)
	}
}

func TestClose(t *testing.T) {
	cache := New(MemoryCacheOptions{
		CleanupInterval: time.Second * 2,
	})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Second * 1),
	})

	time.Sleep(time.Second * 1)
	cache.Close()
	time.Sleep(time.Second * 2)

	if count := cache.Count(); count == 0 {
		t.Fatalf("incorrect result: expected 1, got %d", count)
	}
}

func Test_cleanup(t *testing.T) {
	cache := New(MemoryCacheOptions{
		CleanupInterval: time.Second * 2,
	})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Second * 1),
	})

	time.Sleep(time.Second * 3)

	if count := cache.Count(); count > 0 {
		t.Fatalf("incorrect result: expected 0, got %d", count)
	}
}

func Test_expiration(t *testing.T) {
	cache := New(MemoryCacheOptions{})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Second * 1),
	})

	time.Sleep(time.Second * 2)

	if _, ok := cache.Get(key1); ok {
		t.Fatalf("incorrect result: expected false, got %t", ok)
	}
}

func Test_limit_1(t *testing.T) {
	cache := New(MemoryCacheOptions{
		LimitEntries: 2,
	})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
		Durability: Normal,
	})

	cache.Set(key2, value2, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
		Durability: Weak,
	})

	cache.Set(key3, value3, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
		Durability: Strong,
	})

	if _, ok := cache.Get(key1); !ok {
		t.Fatalf("incorrect result: expected true, got %t", ok)
	}

	if _, ok := cache.Get(key2); ok {
		t.Fatalf("incorrect result: expected false, got %t", ok)
	}
}

func Test_limit_2(t *testing.T) {
	cache := New(MemoryCacheOptions{
		LimitEntries: 1,
	})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
		Durability: Strong,
	})

	cache.Set(key2, value2, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
		Durability: Strong,
	})

	cache.Set(key3, value3, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
		Durability: Strong,
	})

	if _, ok := cache.Get(key3); !ok {
		t.Fatalf("incorrect result: expected true, got %t", ok)
	}
}

func Test_limit_3(t *testing.T) {
	cache := New(MemoryCacheOptions{
		LimitEntries: 2,
	})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * (-5)),
		Durability: Normal,
	})

	cache.Set(key2, value2, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
		Durability: Weak,
	})

	cache.Set(key3, value3, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
		Durability: Strong,
	})

	if _, ok := cache.Get(key1); ok {
		t.Fatalf("incorrect result: expected false, got %t", ok)
	}

	if _, ok := cache.Get(key2); !ok {
		t.Fatalf("incorrect result: expected true, got %t", ok)
	}
}

func Test_limit_4(t *testing.T) {
	cache := New(MemoryCacheOptions{
		LimitEntries: 2,
	})

	cache.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 2),
		Durability: Normal,
	})

	cache.Set(key2, value2, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 1),
		Durability: Normal,
	})

	cache.Set(key3, value3, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
		Durability: Normal,
	})

	if _, ok := cache.Get(key1); !ok {
		t.Fatalf("incorrect result: expected true, got %t", ok)
	}

	if _, ok := cache.Get(key2); ok {
		t.Fatalf("incorrect result: expected false, got %t", ok)
	}
}

func Test_StoreFile(t *testing.T) {
	cache1 := New(MemoryCacheOptions{
		StoreFile: "cache.bin",
	})

	cache1.Set(key1, value1, MemoryCacheEntryOptions{
		Expiration: time.Now().Add(time.Minute * 5),
	})

	cache1.Close()

	cache2 := New(MemoryCacheOptions{
		StoreFile: "cache.bin",
	})

	if _, ok := cache2.Get(key1); !ok {
		t.Fatalf("incorrect result: expected true, got %t", ok)
	}

	if result, _ := cache2.Get(key1); result.(string) != value1 {
		t.Fatalf("incorrect result: expected %s, got %s", value1, result.(string))
	}
}

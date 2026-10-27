package pokecache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			defer cache.Close()

			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Fatalf("expected to find key %q", c.key)
			}
			if string(val) != string(c.val) {
				t.Errorf("Get(%q) = %q, expected %q", c.key, val, c.val)
			}
		})
	}
}

func TestGetMissing(t *testing.T) {
	cache := NewCache(5 * time.Second)
	defer cache.Close()

	if val, ok := cache.Get("missing"); ok || val != nil {
		t.Errorf("Get(missing) = %q, %v; expected nil, false", val, ok)
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond

	cache := NewCache(baseTime)
	defer cache.Close()

	cache.Add("https://example.com", []byte("testdata"))
	if _, ok := cache.Get("https://example.com"); !ok {
		t.Fatal("expected to find key before reaping")
	}

	time.Sleep(waitTime)

	if _, ok := cache.Get("https://example.com"); ok {
		t.Error("expected key to be reaped")
	}
}

func TestAddAndGetCopyValues(t *testing.T) {
	cache := NewCache(5 * time.Second)
	defer cache.Close()

	original := []byte("testdata")
	cache.Add("key", original)
	original[0] = 'X'

	got, _ := cache.Get("key")
	if string(got) != "testdata" {
		t.Fatalf("mutating the added slice changed the cache: got %q", got)
	}

	got[0] = 'Y'
	again, _ := cache.Get("key")
	if string(again) != "testdata" {
		t.Errorf("mutating a returned slice changed the cache: got %q", again)
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewCache(time.Millisecond)
	defer cache.Close()

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() {
			key := fmt.Sprintf("key-%d", i%5)
			for range 100 {
				cache.Add(key, []byte("value"))
				cache.Get(key)
			}
		})
	}
	wg.Wait()
}

func TestCloseIsIdempotent(t *testing.T) {
	cache := NewCache(time.Millisecond)
	cache.Close()
	cache.Close()
}

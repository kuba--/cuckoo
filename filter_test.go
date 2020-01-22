package cuckoo

import (
	"math/rand"
	"testing"
)

func TestIndex(t *testing.T) {
	f := NewFilter(0)
	data := []byte("ABCDEFGHIJKLMNOPRSTUWXYZ")

	x := sum32(data)
	fp := f.fingerprint(x)
	i1 := f.index1(x)
	i2 := f.index2(i1, fp)

	i11 := f.index2(i2, fp)
	i22 := f.index2(i1, fp)

	if i1 != i11 {
		t.Error(i1, i11)
	}
	if i2 != i22 {
		t.Error(i2, i22)
	}
}

func TestFilter(t *testing.T) {
	const maxItems = 32 * 1024
	const capacity = 1024 * maxItems

	f := NewFilter(capacity)
	cnt := f.Count()
	if cnt != 0 {
		t.Error(cnt)
	}

	// Generate test cache
	cache := make(map[string]int)
	for i := 0; i < maxItems; i++ {
		n := 1 + rand.Intn(0xff)
		var item []byte
		for j := 0; j < n; j++ {
			item = append(item, byte(rand.Intn(0xff)))
		}
		key := string(item)
		if _, exists := cache[key]; exists {
			continue
		}

		cache[string(item)] = i
		if err := f.Insert(item); err != nil {
			t.Errorf("[%d]: count: %v: %v\n", i, f.Count(), err)
		}
	}

	// Count number of false positives
	nerrors := 0
	for i := 0; i < maxItems; i++ {
		n := 1 + rand.Intn(0xff)
		var item []byte
		for j := 0; j < n; j++ {
			item = append(item, byte(rand.Intn(0xff)))
		}

		if _, exists := cache[string(item)]; exists != f.Lookup(item) {
			nerrors++
		}
	}

	t.Logf("err rate (%.1f / %.1f): %.5f\n", float32(nerrors), float32(f.Count()), (float32(nerrors) / float32(f.Count())))

	// Cleaning up
	for k := range cache {
		if !f.Delete([]byte(k)) {
			t.Error(k)
		}
	}

	cnt = f.Count()
	if cnt != 0 {
		t.Error(cnt)
	}
}

func BenchmarkInsert(b *testing.B) {
	f := NewFilter(1024 * 1024)
	data := "0123456789ABCDEFGHIJKLMNOPRSTUWXYZ"
	datalen := len(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := []byte(data[rand.Intn(datalen):])
		if f.Lookup(item) {
			continue
		}

		if err := f.Insert(item); err != nil {
			b.Error(err, item)
		}
	}
}

func BenchmarkLookup(b *testing.B) {
	f := NewFilter(1024 * 1024)
	data := "0123456789ABCDEFGHIJKLMNOPRSTUWXYZ"
	datalen := len(data)
	for i := 0; i < datalen; i++ {
		item := []byte(data[i:])
		f.Insert(item)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := []byte(data[rand.Intn(datalen):])
		f.Lookup(item)
	}
}

func BenchmarkDelete(b *testing.B) {
	f := NewFilter(1024 * 1024)
	data := "0123456789ABCDEFGHIJKLMNOPRSTUWXYZ"
	datalen := len(data)
	for i := 0; i < datalen; i++ {
		item := []byte(data[i:])
		f.Insert(item)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := []byte(data[rand.Intn(datalen):])
		f.Delete(item)
	}
}

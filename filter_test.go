package cuckoo

import (
	"math/rand"
	"testing"
)

func TestFilter(t *testing.T) {
	f := NewFilter(0)
	cnt := f.Count()
	if cnt != 0 {
		t.Error(cnt)
	}

	if err := f.Insert([]byte("foo")); err != nil {
		t.Error(err)
	}

	t.Log(f.Lookup([]byte("foo")))
	t.Log(f.Lookup([]byte("bar")))

	cnt = f.Count()
	if cnt != 1 {
		t.Error(cnt)
	}

	t.Log(f.Delete([]byte("foo")))
	t.Log(f.Delete([]byte("bar")))

	cnt = f.Count()
	if cnt != 0 {
		t.Error(cnt)
	}
}

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

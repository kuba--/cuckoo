package internal

import "testing"

func TestBucket(t *testing.T) {
	data := "abcdefghijklmnoprstuwxyz"
	b := &Bucket{}

	for i := 0; i < BucketSize; i += 1 {
		c := data[i]
		item := byte(c)
		if b.Lookup(item) {
			t.Errorf("Lookup(%v)\n", c)
		}

		if !b.Insert(item) {
			t.Errorf("Insert(%v)\n", c)
		}
	}

	for i := 0; i < BucketSize; i += 1 {
		c := data[i]
		item := byte(c)
		if !b.Lookup(item) {
			t.Errorf("Lookup(%v)\n", c)
		}

		if !b.Delete(item) {
			t.Errorf("Delete(%v)\n", c)
		}
	}
}

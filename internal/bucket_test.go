package internal

import "testing"

func TestBucket(t *testing.T) {
	data := "abcdefghijklmnoprstuwxyz"
	b := &Bucket{}

	for i := 0; i < BucketSize; i += 1 {
		c := data[i]
		item := byte(c)
		if b.Lookup(item) {
			t.Errorf("Lookup(%v == %v): true\t(expected false)\n", c, item)
		}

		if !b.Insert(item) {
			t.Errorf("Insert(%v == %v): %v(expected true)\n", c, item)
		}
	}

	for i := 0; i < BucketSize; i += 1 {
		c := data[i]
		item := byte(c)
		if !b.Lookup(item) {
			t.Errorf("Lookup(%v == %v): false\t(expected true)\n", c, item)
		}

		if !b.Delete(item) {
			t.Errorf("Delete(%v == %v): false\t(expected true)\n", c, item)
		}
	}
}

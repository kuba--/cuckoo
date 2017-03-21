package cuckoo

import "github.com/kuba--/cuckoo"

var f = cuckoo.NewFilter(32 * 1024 * 1024)

// Fuzz test function
func Fuzz(data []byte) int {
	if f.Lookup(data) {
		if !f.Delete(data) {
			return 0
		}
	} else {
		if err := f.Insert(data); err != nil {
			panic(err)
		}
	}
	return 1
}

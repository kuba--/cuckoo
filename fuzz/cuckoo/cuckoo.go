package cuckoo

import (
	"log"

	"github.com/kuba--/cuckoo"
)

func Fuzz(data []byte) int {
	f := cuckoo.NewFilter(32 * 1024 * 1024)

	if f.Lookup(data) {
		if !f.Delete(data) {
			return 0
		}
	} else {
		if err := f.Insert(data); err != nil {
			log.Println(err)
			return 0
		}
	}
	return 1
}

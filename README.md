[![GoDoc](https://godoc.org/github.com/kuba--/cuckoo?status.svg)](http://godoc.org/github.com/kuba--/cuckoo)
[![Go Report Card](https://goreportcard.com/badge/github.com/kuba--/cuckoo)](https://goreportcard.com/report/github.com/kuba--/cuckoo)
[![Build Status](https://travis-ci.org/kuba--/cuckoo.svg?branch=master)](https://travis-ci.org/kuba--/cuckoo)

# Cuckoo Filter: Practically Better Than Bloom
Package _cuckoo_ provides a Cuckoo Filter (Practically Better Than Bloom).
Cuckoo filters provide the ﬂexibility to add and remove items dynamically.
A cuckoo filter is based on cuckoo hashing (and therefore named as cuckoo filter).
It is essentially a cuckoo hash table storing each key's fingerprint.

Implementation is based heavily on whitepaper: "Cuckoo Filter: Better Than Bloom" by Bin Fan, Dave Andersen and Michael Kaminsky
(https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf).

### Example
```go
    var item []byte

    // Create a new filter with 32MB capacity.
    f := cuckoo.NewFilter(32 * 1024 * 1024)

    if f.Lookup(item) {
        if !f.Delete(item) {
            // Cannot delete the item
            // ...
        }
    } else {
        if err := f.Insert(item); err != nil {
            // Hashtable is considered full
            // ...
        }
    }
```

### Fuzzing - randomized testing
- get go-fuzz
```
    $ go get github.com/dvyukov/go-fuzz/go-fuzz
    $ go get github.com/dvyukov/go-fuzz/go-fuzz-build
```

- build the test program with necessary instrumentation
```
    $ go-fuzz-build github.com/kuba--/cuckoo/fuzz/cuckoo
```

- ready to go
```
    $ go-fuzz -bin=./cuckoo-fuzz.zip -workdir=/tmp/cuckoo
```


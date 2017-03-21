/*
Package cuckoo provides a Cuckoo Filter (Practically Better Than Bloom).
Cuckoo filters provide the ﬂexibility to add and remove items dynamically.
A cuckoo filter is based on cuckoo hashing (and therefore named as cuckoo filter).
It is essentially a cuckoo hash table storing each key's fingerprint.

Implementation is based heavily on whitepaper: "Cuckoo Filter: Better Than Bloom" by Bin Fan, Dave Andersen and Michael Kaminsky
(https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf).
*/
package cuckoo

import (
	"fmt"
	"hash/fnv"
	"math/rand"

	"github.com/kuba--/cuckoo/internal"
)

// MaxNumKicks stands for maximum number of cuckoo kicks before insert failure.
var MaxNumKicks uint32 = 512

// DefaultHash32 is a default hash function.
var DefaultHash32 = fnv.New32a()

// Filter keeps hashtable with fingerprints.
type Filter struct {
	buckets []internal.Bucket
	count   uint32
}

// NewFilter creates a new Cuckoo Filter with given capacity (rounding up to the nearest power of 2).
func NewFilter(capacity uint32) *Filter {
	var len uint32 = 1
	if capacity > internal.BucketSize {
		len = power2(capacity) / internal.BucketSize
	}
	return &Filter{buckets: make([]internal.Bucket, len)}
}

/*
Insert inserts an item to the filter.
Pseudo code:
	f = fingerprint(x);
	i1 = hash(x);
	i2 = i1 ⊕ hash(f);
	if bucket[i1] or bucket[i2] has an empty entry then
		add f to that bucket;
		return Done;

	// must relocate existing items;
	i = randomly pick i1 or i2;
	for n = 0; n < MaxNumKicks; n++ do
		randomly select an entry e from bucket[i];
		swap f and the fingerprint stored in entry e;
		i = i ⊕ hash( f );

	if bucket[i] has an empty entry then
		add f to bucket[i];
		return Done;

	// Hashtable is considered full;
	return Failure;
*/
func (f *Filter) Insert(item []byte) error {
	x := sum32(item)

	fp := f.fingerprint(x)
	i := f.index1(x)
	ok := f.buckets[i].Insert(fp)
	if !ok {
		i = f.index2(i, fp)
		ok = f.buckets[i].Insert(fp)
	}

	for n := uint32(0); !ok && n < MaxNumKicks; n++ {
		j := rand.Intn(internal.BucketSize)
		e := f.buckets[i][j]
		f.buckets[i][j] = fp

		fp = e
		i = f.index2(i, fp)
		ok = f.buckets[i].Insert(fp)
	}

	if !ok {
		return fmt.Errorf("Hashtable is considered full (MaxNumKicks: %v, Count: %v)", MaxNumKicks, f.Count())
	}

	f.count++
	return nil
}

/*
Lookup reports if the item is inserted, with false positive rate.
Pseudo code:
	f = fingerprint(x);
	i1 = hash(x);
	i2 = i1 ⊕ hash(f);
	if bucket[i1] or bucket[i2] has f then
		return True;
	return False;
*/
func (f *Filter) Lookup(item []byte) bool {
	x := sum32(item)

	fp := f.fingerprint(x)
	i := f.index1(x)
	ok := f.buckets[i].Lookup(fp)
	if !ok {
		i = f.index2(i, fp)
		ok = f.buckets[i].Lookup(fp)
	}

	return ok
}

/*
Delete deletes an item from teh filter.
Pseudo code:
	f = fingerprint(x);
	i1 = hash(x);
	i2 = i1 ⊕ hash(f);
	if bucket[i1] or bucket[i2] has f then
		remove a copy of f from this bucket;
		return True;
	return False;
*/
func (f *Filter) Delete(item []byte) bool {
	x := sum32(item)

	fp := f.fingerprint(x)
	i := f.index1(x)
	ok := f.buckets[i].Delete(fp)
	if !ok {
		i = f.index2(i, fp)
		ok = f.buckets[i].Delete(fp)
	}

	if ok {
		f.count--
	}
	return ok
}

/*
Count returns how many items are in the filter.
*/
func (f *Filter) Count() uint32 {
	return f.count
}

// fingerprint generates a byte tag for given Hash32 value.
func (*Filter) fingerprint(x uint32) byte {
	fp := x & 0xff
	if fp == 0 {
		fp++
	}
	return byte(fp)
}

// index1 generates an index for given Hash32 value
func (f *Filter) index1(x uint32) uint32 {
	return x & uint32(len(f.buckets)-1)
}

// index2 generates alternative index.
func (f *Filter) index2(i1 uint32, fp byte) uint32 {
	return i1 ^ f.index1(sum32([]byte{fp}))
}

// sum32 generates a Hash32 value for given item.
func sum32(item []byte) uint32 {
	DefaultHash32.Reset()
	DefaultHash32.Write(item)
	return DefaultHash32.Sum32()
}

// power2 returns nearest power of 2 for given number.
func power2(n uint32) uint32 {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

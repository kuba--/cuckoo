package internal

// BucketSize defines number of bytes per Bucket
const BucketSize = 4

// Bucket represents a row in Cuckoo Filter's hashtable.
type Bucket [BucketSize]byte

/*
Insert tries to insert a fingerprint into the bucket.
If all slots are occupied (!= 0) return false, otherwise true.
*/
func (b *Bucket) Insert(fp byte) bool {
	for i, fpi := range b {
		if fpi == 0 {
			b[i] = fp
			return true
		}
	}
	return false
}

// Lookup lineary checks if given fingerprint exists in the bucket.
func (b *Bucket) Lookup(fp byte) bool {
	for _, fpi := range b {
		if fpi == fp {
			return true
		}
	}
	return false
}

// Delete tries to find and reset a slot for given fingerprint.
func (b *Bucket) Delete(fp byte) bool {
	for i, fpi := range b {
		if fpi == fp {
			b[i] = 0
			return true
		}
	}
	return false
}

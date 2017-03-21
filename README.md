# Cuckoo Filter: Practically Better Than Bloom
Package _cuckoo_ provides a Cuckoo Filter (Practically Better Than Bloom).
Cuckoo filters provide the ﬂexibility to add and remove items dynamically.
A cuckoo filter is based on cuckoo hashing (and therefore named as cuckoo filter).
It is essentially a cuckoo hash table storing each key's fingerprint.

Implementation is based heavily on whitepaper: "Cuckoo Filter: Better Than Bloom" by Bin Fan, Dave Andersen and Michael Kaminsky
(https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf).


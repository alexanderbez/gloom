# gloom


[![GoDoc](https://godoc.org/github.com/alexanderbez/gopq?status.svg)](https://godoc.org/github.com/alexanderbez/gopq)
[![Build Status](https://travis-ci.org/alexanderbez/gloom.svg?branch=master)](https://travis-ci.org/alexanderbez/gloom)

A simple and intuitive implementation of a Bloom filter using enhanced double
hashing.

Within gloom, enhanced double hashing is used to set bit positions. The choice
for double hashing was shown to be effective without any loss in the asymptotic
false positive probability, leading to less computation and potentially less
need for randomness in practice by Adam Kirsch and Michael Mitzenmacher in
[Less Hashing, Same Performance: Building a Better Bloom Filter](http://citeseerx.ist.psu.edu/viewdoc/download;jsessionid=AF0A7F109B5F97C758DD773942A1054F?doi=10.1.1.152.579&rep=rep1&type=pdf).

The enhanced double hash is of the form:

g<sub>i</sub>(x) = H<sub>1</sub>(x) + iH<sub>2</sub>(x) + f(i), where

H<sub>1</sub>
is FNV-1a 64-bit, H<sub>2</sub> is Murmur3 64-bit, and f(i) = i<sup>3</sup>


## What is a Bloom filter?

A Bloom filter is a space-efficient probabilistic data structure used for set
inclusion queries. False positive matches are possible, but false negatives are
not – in other words, a query returns either "possibly in set" or "definitely
not in set".

Essentially, a Bloom filter contains a single bit vector of size `m` and `k`
independent and uniform hash functions to insert `n` set items. Upon inserting
a set item into the filter, the `k` hash functions return all bit vector positions
to set for said item.

To test inclusion of an item, all mapped `k` bit positions must contain a set
bit. It is possible to get a false positive for an item, but under a desired
and given probability. Although Bloom filters allow false positives, the
space savings often outweigh this drawback.

## API

To initialize a Bloom filter, a given set size and desired false positive
probability is needed. The size of the bit vector and the number of hash functions
to use is determined by these parameters. 

```golang
import (
   	"github.com/alexanderbez/gloom"
)

bf, err := gloom.NewBloomFilter(n, gloom.DefaultFalsePosProb)

item := []byte("foo")
bf.Set(item)

ok, err := bf.Includes(item)
```

## Tests

```shell
$ go test -v ./...
```

## Contributing

1. [Fork it](https://github.com/alexanderbez/gloom/fork)
2. Create your feature branch (`git checkout -b feature/my-new-feature`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to the branch (`git push origin feature/my-new-feature`)
5. Create a new Pull Request
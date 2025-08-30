# wavelet-tree

> A golang implementation of the [wavelet tree](https://en.wikipedia.org/wiki/Wavelet_Tree) datastructure.

## Features

- Efficient storage with bit-packing.
- Fast calculation of `rank(i)` with `RRR`.

## Why?

The compression ratio for the wavelet tree isn't the greatest when
compared to other string compression algorithms but it is much
faster for doing lookups.

So if you need to store multiple strings in memory and have random
access while still enjoying the memory efficiency improvements of
compression, you may want to consider a wavelet tree.

## Explainers

Explainers for the algorithms can be found under [docs](./docs).


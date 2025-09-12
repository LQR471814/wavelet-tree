# wavelet-tree

> A golang implementation of the [wavelet tree](https://en.wikipedia.org/wiki/Wavelet_Tree) datastructure.

## Features

- Efficient storage with bit-packing.
- Fast calculation of $\text{rank}(b, i)$ with RRR.
- Well-tested.

## What does this even do?

You may not have heard of the wavelet tree data structure before,
so here's a quick sales pitch.

The wavelet tree data structure is a string compression
data structure that is very fast at doing *lookups* with a decent
compression ratio.

Here's an example of a scenario where it would be useful:

> Let's suppose you have to store (some large number) of names in
> a database. You want to save on storage costs by compressing
> these strings, but cannot compromise on lookup times when you
> need to lookup a name associated with a customer ID.
>
> If you use a traditional string compression algorithm like LZ4,
> you will need to decompress and recompress all of the strings at
> once per transaction (which will increase latency), compress
> each string individually (which will increase overhead), or
> compress in chunks (which will only have dubious effects on
> latency).
>
> You also cannot just keep all the strings uncompressed in
> memory! Remember, this is a *very* large number of strings. You
> may be able to offload uncompressed strings into a large
> persistent storage medium, but you will suffer latency costs for
> it.
>
> If only there was a data structure capable of compressing many
> strings while retaining fast and efficient random access of such
> strings.

Behold, the wavelet tree.

## Explainers

Explainers for the algorithms can be found under [docs](./docs).


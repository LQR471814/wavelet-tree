# RRR

The RRR datastructure primarily exists to speed up the bitwise
operations $\text{rank}(b, i)$ and $\text{select}(b, i)$.

- $\text{rank}(b,i)$ - Returns the number of bits of value $b$ (0
  or 1) up to a bit index $i$.
- $\text{select}(b, i)$ - Returns the $i$'th bit $b$ (0 or 1) in
  the bit vector.

## Blocks

RRR divides up the bitvector into **blocks** of bits, each of
fixed-size $b$.

Each block stores:

- **Class:** The number of `1`s in the block.
- **Offset:** A number representing the unique combination of the
  `1`s positions in the block, given its **class**.

As such, the block effectively functions as a cache for
calculations critical for `rank` and `select`.

Usually, the theoretical optimal value for the block size $b$ is
calculated with the following formula:

$$
b = \frac{\log_{2}(n)}{2}
$$

Where $n$ is the number of bits the bitvector is storing.

I will not attempt to prove this, you can likely find a proof for
it in the [original paper](https://arxiv.org/abs/0705.0552).

> [!NOTE]
> Current CPUs can only handle up to 64-bits in a single
> instruction, so technically that means that $b$ (and by
> extension, $n$) has a "maximum value" for which performance hits
> may come if exceeded.
>
> This is a sound line of thinking, however if we set $b= 64$ and
> solve for the maximum $n$ corresponding with the optimal usage
> of 64 bits for one block, we would get $2^{128}$.
>
> Needless to say, the optimal value for $b$ will probably never
> even come close to 64 bits, so using the optimal value for $b$
> will be fine.

## Superblocks

If we have many blocks, we may run into slowdowns when querying
for rank. Thankfully, we can add **superblocks** to accelerate
rank queries.

A superblock is effectively a block of blocks of a fixed-size $k$.

Each superblock stores the **cumulative rank** up to the start of
the superblock.

That way querying rank does not involve summing up the class of
all the blocks up to the target bit index $i$, but only requires
one to skip to the superblock that contains $i$ and look up the
cumulative rank.

The fact that this lookup is only stored per $k$ blocks also
reduces the memory overhead of storing this lookup.

Like $b$, $k$ also has an optimal value that is usually given by
the formula:

$$
k = \log_{2}(n)
$$


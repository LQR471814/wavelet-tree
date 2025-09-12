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

## Finding offsets

The [combinatorial number system](https://en.wikipedia.org/wiki/Combinatorial_number_system)
gives the relationship between a number and the possible
combinations of unique subsets of a size $r$ for a given set $S$.

> [!NOTE]
> The letter $k$ is more commonly used for the size of the subset
> rather than $r$. I used $r$ because it is more commonly used
> when talking about counting.

The total number of these subsets is more commonly computed with

$$
\binom{n}{r}
$$

Where $n$ is the number of elements in set $S$ and $r$ is the size
of each unique subset.

Let's suppose we have a 5 element set $S$ which we want to find
all combinations of 3 element subsets.

If we were to list them out in a tabular format, it would look
something like this:

![mapping of subsets into bitvectors](https://upload.wikimedia.org/wikipedia/commons/8/85/Combinatorial_number_system%3B_5_choose_3.svg)

> WatchduckYou can name the author as "T. Piesk", "Tilman Piesk"
> or "Watchduck"., CC BY 4.0
> <https://creativecommons.org/licenses/by/4.0>, via Wikimedia
> Commons

You can ignore the numbers inside the red boxes. But doesn't this
look just like all the possible combinations of 5-bit block
where 3 bits are set to 1?

In other words, this provides the **offset** for all the possible
combinations of blocks of **class** 3.

Well then the question becomes, how do you compute the offset for
a given combination?

The offset $N$ of a given subset is given by the following
relationship.

$$
N = \binom{c_{r}}{r} + \dots + \binom{c_{2}}{2} + \binom{c_1}{1}
  = \sum_{i = 1}^{r} \binom{c_{i}}{i}
$$

The values of $c_{i}$ follow a strictly decreasing relationship:

$$
c_{i} > ... > c_{2} > c_{1} \geq 0
$$

> Using the row at offset 2 in the table. (the third row)
>
> The subset would be $c=\{0, 2, 3\}$.
> ($c_{1}= 0$, $c_{2} = 2$, $c_{3} = 3$)
>
> $$R = \binom{0}{1} + \binom{2}{2} + \binom{3}{3} = 2$$
>
> The offset is the expected value 2.

This operation of finding a number for a particular subset is
commonly called "ranking". (though it is different from the
$\text{rank}(b, i)$ operation of RRR)

The opposite process, "unranking" derives a subset from a given
offset.

Unranking involves an algorithm better expressed with pseudo-code:

```
while N > 0:
    find the largest value of 'c' such that nCr(c, r) <= N
    N = N - nCr(c, r)
    r = r - 1
```


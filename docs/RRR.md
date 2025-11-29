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

## Finding offsets for a block

The [combinatorial number system](https://en.wikipedia.org/wiki/Combinatorial_number_system)
gives the relationship between a number and the possible
combinations of unique subsets of a size $k$ for a given set $S$.

The total number of these subsets is more commonly computed with

$$
\binom{n}{k}
$$

Where $n$ is the number of elements in set $S$ and $k$ is the size
of each unique subset.

Let's suppose we have a 5 element set $S$ which we want to find
all combinations of 3 element subsets.

If we were to list them out in a tabular format, it would look
something like this:

![subsets table](https://upload.wikimedia.org/wikipedia/commons/8/85/Combinatorial_number_system%3B_5_choose_3.svg)

> WatchduckYou can name the author as "T. Piesk", "Tilman Piesk"
> or "Watchduck"., CC BY 4.0
> <https://creativecommons.org/licenses/by/4.0>, via Wikimedia
> Commons

You can ignore the numbers inside the red boxes. But doesn't this
look just like all the possible combinations of a 5-bit block
where 3 bits are set to 1? (where everything in red is a 1,
everything white is a 0)

In other words, the row number provides the **offset** for all the
possible combinations of blocks of **class** 3. A unique number
associated with every possible combination of block with class 3.

Then the question becomes, how do you compute the offset for a
given combination?

Let's say the set consists of the 1-bit positions (starting from
0), say (ordered from least to greatest):

$B = [0, 1, 3]$

The offset of the subset $B$, is given by:

$$
\sum_{i=0}^{|B|-1} \binom{B[i]}{i+1}
$$

> Where $B[i]$ means the element at index $i$ in the list $B$, and
> $|B|$ gives the length of the list.

$$
\binom{0}{1} + \binom{1}{2} + \binom{3}{3} = 1
$$

This operation of finding a number for a particular subset is
commonly called "ranking". (though it is different from the
$\text{rank}(b, i)$ operation of RRR)

## Finding a block from an offset

The opposite process, "unranking" derives a block (possible
subset) from a given offset.

This process is not unlike the conversion of numbers between
different bases. Just that it involves combinatorics instead of
exponentials.

Here is an example to illustrate:

Suppose you were to find the combination of size 4 for a set of 8
elements at offset 30.

The values of $\binom{n}{4}$ for values of $n=5,6,7$ you have the
values $5,15,35$.

The largest value $15 \leq 30$ is $\binom{6}{4}$. We know that the
4th 1-bit position is $6$. ($[a,b,c,6]$)

We then do this again for $15$, the values of $\binom{n}{3}$ for
successive values of $n=4,5,6$ are $4,10,20$, so $10 \leq 15$, the
3rd 1-bit position is $5$. ($[a,b,5,6]$)

We then do this again for $5$. $n=3$, $\binom{3}{2}=3$, The 2nd
1-bit position is $3$ ($[a,3,5,6]$).

We then do this again for $2$. $n=2$, $\binom{2}{1}=4$. The 1st
1-bit position is $2$? ($[2,3,5,6]$)

Thus, the bit vector for block of size 8, with class 4, and offset
30 is:

$[0,0,1,1,0,1,1,0]$

Unranking involves an algorithm better expressed with pseudo-code:

> [!NOTE]
> If you end up with a value of $0$ but still have combination
> positions to fill. Then you would want to choose the "maximal"
> value for the combination element at index $i$. Namely, the
> maximum value for $n$ such that $\binom{n}{i+1}=0$, this would
> be $n=i$. So for the remaining combination positions, you would
> simply set the positions to be equal to $i$.


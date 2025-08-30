# Bitvector

RRR and wavelet-tree work with bits, but modern computer systems
work with bytes. Thus, a conversion layer must exist between them.

The bitvector contains 3 types of operations, $\text{get}(s, i)$,
    $\text{set}(s,i,v)$, and $\text{append}(s,v)$.

- $s$ indicates the number of bits to get/set/append.
- $i$ indicates the bit index at which to get/set the $s$ number
  of bits.
- $v$ is the unsigned integer value that contains the actual bits
  to be set/appended.

Each operand contains variants for 8, 16, 32, and 64 bit unsigned
integers respectively. This allows one to encode bitwise
information in a memory efficient manner. (instead of using at
least a byte for each bit)

The specifics on how the operands are implemented are not too
complicated, the hardest part is simply account for bits that
spillover to the next byte/unsigned integer.


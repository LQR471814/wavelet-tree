package wavelettrie

// divide bitvector into fixed-size blocks
// usually b (size) = log(n) / 2 bits
// for each block store:
//   - the # of 1's (class)
//   - the position of 1's (offset)
//       - the offset is the index of the current block in the list of possible
//       combinations of patterns with the given block size and class
//       - ex. 1 0 1 1
//       - 4 bits has C(4,3)

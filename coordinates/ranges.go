package coordinates

// A RangeID represents a sequential list of IDs. The two
// boundary numbers are inclusive. A RangeID of [3, 7] implies
// the list of IDs: {3, 4, 5, 6, 7} (in incrementing order).
type RangeID [2]int64

// A RangeIDs represents an increasing list of RangeID. A
// RangeIDs of [[2, 4], [8, 8], [13, 14]] implies the list of IDs:
// {2, 3, 4, 8, 13, 14} (in increasing order).
type RangeIDs [][2]int64

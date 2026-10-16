package jutil

// Go-only: language-adapter — Arrays.binarySearch(int[], int).

// ArraysBinarySearchInts is Arrays.binarySearch(int[] a, int key) on a
// sorted slice: the index of key, or -(insertion point) - 1 when key is
// absent, the insertion point being the index of the first element greater
// than key (len(a) if there is none), as the Javadoc specifies.
//
// The Javadoc leaves open which index is returned when key occurs several
// times. Observed behaviour (black-box): it is the one the textbook
// bisection finds, probing the middle index (low+high)/2 (rounded down) of
// the remaining range [low, high] and returning on the first equal element;
// checked against 200,000 random sorted arrays with many duplicates.
func ArraysBinarySearchInts(a []int, key int) int {
	low, high := 0, len(a)-1
	for low <= high {
		mid := int(uint(low+high) >> 1)
		switch v := a[mid]; {
		case v < key:
			low = mid + 1
		case v > key:
			high = mid - 1
		default:
			return mid
		}
	}
	return -(low + 1)
}

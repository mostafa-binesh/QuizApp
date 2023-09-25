package utils

// removeElement removes the element at the given index from the input array.
// The input array must be a slice and the index must be within the bounds of the slice.
func RemoveElement[T any](arr []T, index int) []T {
	// Copy the elements that come before the index
	result := make([]T, len(arr)-1)
	copy(result, arr[:index])

	// Copy the elements that come after the index
	copy(result[index:], arr[index+1:])

	return result
}

// removeElementByRef removes the element at the given index from the input slice.
// The index must be within the bounds of the input slice.
func RemoveElementByRef[T any](arr *[]T, index int) {
	// Move the elements after the index down by one position
	copy((*arr)[index:], (*arr)[index+1:])

	// Truncate the slice to remove the last element
	*arr = (*arr)[:len(*arr)-1]
}

// ExistsInArray checks if an element exists in an array of comparable type.
func ExistsInArray[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

// eg. gets []string and returns []*string
func ConvertSliceToPtrSlice[T any](s []T) []*T {
	ptrSlice := make([]*T, len(s))
	for i, str := range s {
		ptrSlice[i] = &str
	}
	return ptrSlice
}

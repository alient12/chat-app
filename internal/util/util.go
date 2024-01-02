package util

func InSlice(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func IsSubset(sliceA, sliceB []uint64) bool {
	for _, val := range sliceA {
		if !InSlice(sliceB, val) {
			return false
		}
	}
	return true
}

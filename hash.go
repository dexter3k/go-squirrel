package squirrel

import (
	"unsafe"
)

type Hash uint64

func HashString(source string) Hash {
	var result Hash
	count := len(source) / ((len(source) >> 5) | 1)
	for i := 0; i < count; i++ {
		result ^= (result << 5) + (result >> 2) + Hash(source[i])
	}
	return result
}

func HashInteger(source int64) Hash {
	return Hash(source)
}

func HashFloat(source float64) Hash {
	return Hash(source)
}

func HashBool(source bool) Hash {
	if source {
		return 1
	} else {
		return 0
	}
}

func HashPointer(source interface{}) Hash {
	return Hash(uintptr(unsafe.Pointer(&source))) >> 3
}

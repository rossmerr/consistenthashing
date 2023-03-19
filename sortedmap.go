package consistenthashing

import (
	"bytes"
)

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type SortedMap[K Unsigned, T any] struct {
	arr []K
	dic map[K]T
}

func NewSortedMap[K Unsigned, T any]() *SortedMap[K, T] {
	return &SortedMap[K, T]{
		arr: []K{},
		dic: map[K]T{},
	}
}

func (s *SortedMap[K, T]) Add(key K, item T) {
	if _, ok := s.dic[key]; !ok {
		s.arr = append(s.arr, key)
		s.dic[key] = item
		bubbleSortFunc(s.arr, func(a, b K) bool { return a < b })
	}
}

func (s *SortedMap[K, T]) Remove(key K) {
	if _, ok := s.dic[key]; ok {
		for i, v := range s.arr {
			if v == key {
				s.arr = append(s.arr[:i], s.arr[i+1:]...)
			}
		}
		delete(s.dic, key)
	}
}

func (s *SortedMap[K, T]) Get(key K) *T {
	if v, ok := s.dic[key]; ok {
		return &v
	}
	return nil
}

func (s *SortedMap[K, T]) Contains(key K) bool {
	_, ok := s.dic[key]
	return ok
}

func (s *SortedMap[K, T]) Keys() []K {
	return s.arr
}

func (s *SortedMap[K, T]) Empty() bool {
	return len(s.arr) == 0
}

func (s *SortedMap[K, T]) First() *T {
	if len(s.arr) > 0 {
		k := s.arr[0]
		v := s.dic[k]
		return &v
	}

	return nil
}

func (s *SortedMap[K, T]) Tail(key K) []K {
	b := intToBytes(key)
	for i, v := range s.arr {
		if bytes.HasSuffix(intToBytes(v), b) {
			return s.arr[i:]
		}
	}

	return []K{}

}

func bubbleSortFunc[T any](x []T, less func(a, b T) bool) {
	n := len(x)
	for {
		swapped := false
		for i := 1; i < n; i++ {
			if less(x[i], x[i-1]) {
				x[i-1], x[i] = x[i], x[i-1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

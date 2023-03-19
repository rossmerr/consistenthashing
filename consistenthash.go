package consistenthashing

import (
	"hash/crc32"
	"strconv"
)

type HashFunction func([]byte) uint

type ConsistentHash[T Hash] struct {
	replicate uint
	circle    *SortedMap[uint, T]
	hash      HashFunction
}

type ConsistentHashOption[T Hash] func(*ConsistentHash[T])

func WithHash[T Hash](hash HashFunction) ConsistentHashOption[T] {
	return func(s *ConsistentHash[T]) {
		s.hash = hash
	}
}

func NewConsistentHash[T Hash](replicate uint, opts ...ConsistentHashOption[T]) *ConsistentHash[T] {
	consistentHash := &ConsistentHash[T]{
		replicate: replicate,
		circle:    NewSortedMap[uint, T](),
		hash: func(b []byte) uint {
			if len(b) < 64 {
				var scratch [64]byte
				copy(scratch[:], b)
				return uint(crc32.ChecksumIEEE(scratch[:len(b)]))
			}
			return uint(crc32.ChecksumIEEE([]byte(b)))
		},
	}

	for _, opt := range opts {
		opt(consistentHash)
	}

	return consistentHash
}

func (s *ConsistentHash[T]) Add(node T) {
	for i := uint(0); i < s.replicate; i++ {
		hash := s.hash(append(intToBytes(node.Sum()), intToBytes(i)...))
		s.circle.Add(hash, node)
	}
}

func (s *ConsistentHash[T]) Remove(node T) {
	for i := uint(0); i < s.replicate; i++ {
		hash := s.hash(append(intToBytes(node.Sum()), intToBytes(i)...))
		s.circle.Remove(hash)
	}
}

func (s *ConsistentHash[T]) Get(key uint) T {
	if s.circle.Empty() {
		var noop T
		return noop
	}
	hash := s.hash(intToBytes(key))

	if !s.circle.Contains(hash) {

		tail := s.circle.Tail(hash)
		if len(tail) > 0 {
			return s.circle.Get(tail[0])
		}
		return s.circle.First()
	}
	return s.circle.Get(hash)
}

func intToBytes(num uint) []byte {
	return []byte(strconv.Itoa(int(num)))
}

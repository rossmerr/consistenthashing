package consistenthashing_test

import (
	"testing"

	"github.com/rossmerr/consistenthashing"
)

type test struct {
	sum uint
}

func (s *test) Sum() uint {
	return s.sum
}

func NewTest(s uint) *test {
	return &test{
		sum: s,
	}
}

func TestNewConsistentHash(t *testing.T) {

	tests := []struct {
		name      string
		replicate uint
		objects   []*test
		keys      []uint
		want      []uint
	}{
		{
			name:      "",
			replicate: 3,
			objects: []*test{
				NewTest(1),
				NewTest(2),
				NewTest(3),
			},
			keys: []uint{
				1, 2, 3,
			},
			want: []uint{2, 1, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := consistenthashing.NewConsistentHash[*test](tt.replicate)
			for _, v := range tt.objects {
				got.Add(v)
			}

			for i, k := range tt.keys {
				v := got.Get(k)

				if v.Sum() != tt.want[i] {
					t.Errorf("NewConsistentHash() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestNewConsistentHash_Get(t *testing.T) {

	tests := []struct {
		name      string
		replicate uint
	}{
		{
			name:      "empty circle nil response",
			replicate: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := consistenthashing.NewConsistentHash[*test](tt.replicate)
			v := got.Get(10)

			if v != nil {
				t.Errorf("NewConsistentHash() = %v, want nil", v)
			}

		})
	}
}

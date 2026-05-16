package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Parallel()

	type args struct {
		pred func(int) int
		in   []int
	}
	type expected struct {
		value []int
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "maps_values",
			args: args{
				in: []int{1, 2, 3},
				pred: func(value int) int {
					return value * 2
				},
			},
			expected: expected{
				value: []int{2, 4, 6},
			},
		},
		{
			name: "empty",
			args: args{
				in: []int{},
				pred: func(value int) int {
					return value * 2
				},
			},
			expected: expected{
				value: []int{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := Map(tt.args.in, tt.args.pred)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestSliceToMap(t *testing.T) {
	t.Parallel()

	type item struct {
		name string
		id   int
	}
	type args struct {
		pred func(item) (int, string)
		in   []item
	}
	type expected struct {
		value map[int]string
	}
	tests := []struct {
		expected expected
		name     string
		args     args
	}{
		{
			name: "maps_by_key",
			args: args{
				in: []item{
					{id: 1, name: "first"},
					{id: 2, name: "second"},
				},
				pred: func(value item) (int, string) {
					return value.id, value.name
				},
			},
			expected: expected{
				value: map[int]string{
					1: "first",
					2: "second",
				},
			},
		},
		{
			name: "overwrites_duplicate_key",
			args: args{
				in: []item{
					{id: 1, name: "first"},
					{id: 1, name: "updated"},
				},
				pred: func(value item) (int, string) {
					return value.id, value.name
				},
			},
			expected: expected{
				value: map[int]string{
					1: "updated",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := SliceToMap(tt.args.in, tt.args.pred)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	type args struct {
		pred func(int) bool
		in   []int
	}
	type expected struct {
		value []int
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "filters_values",
			args: args{
				in: []int{1, 2, 3, 4},
				pred: func(value int) bool {
					return value%2 == 0
				},
			},
			expected: expected{
				value: []int{2, 4},
			},
		},
		{
			name: "no_matches",
			args: args{
				in: []int{1, 3},
				pred: func(value int) bool {
					return value%2 == 0
				},
			},
			expected: expected{
				value: []int{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := Filter(tt.args.in, tt.args.pred)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestFind(t *testing.T) {
	t.Parallel()

	type args struct {
		pred func(int) bool
		in   []int
	}
	type expected struct {
		value int
		found bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "found",
			args: args{
				in: []int{1, 2, 3},
				pred: func(value int) bool {
					return value > 1
				},
			},
			expected: expected{
				value: 2,
				found: true,
			},
		},
		{
			name: "not_found",
			args: args{
				in: []int{1, 2, 3},
				pred: func(value int) bool {
					return value > 5
				},
			},
			expected: expected{
				value: 0,
				found: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, found := Find(tt.args.in, tt.args.pred)

			assert.Equal(t, tt.expected.value, value)
			assert.Equal(t, tt.expected.found, found)
		})
	}
}

func TestEnsureLengthNotExceeding(t *testing.T) {
	t.Parallel()

	type args struct {
		in        []int
		maxLength int
	}
	type expected struct {
		value []int
	}
	tests := []struct {
		name     string
		expected expected
		args     args
	}{
		{
			name: "shorter",
			args: args{
				in:        []int{1, 2},
				maxLength: 3,
			},
			expected: expected{
				value: []int{1, 2},
			},
		},
		{
			name: "equal",
			args: args{
				in:        []int{1, 2, 3},
				maxLength: 3,
			},
			expected: expected{
				value: []int{1, 2, 3},
			},
		},
		{
			name: "longer",
			args: args{
				in:        []int{1, 2, 3, 4},
				maxLength: 2,
			},
			expected: expected{
				value: []int{1, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := EnsureLengthNotExceeding(tt.args.in, tt.args.maxLength)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestGroupBy(t *testing.T) {
	t.Parallel()

	type item struct {
		group string
		value int
	}
	type args struct {
		keyFunc func(item) string
		in      []item
	}
	type expected struct {
		value map[string][]item
	}
	tests := []struct {
		expected expected
		name     string
		args     args
	}{
		{
			name: "groups_values",
			args: args{
				in: []item{
					{group: "odd", value: 1},
					{group: "even", value: 2},
					{group: "odd", value: 3},
				},
				keyFunc: func(value item) string {
					return value.group
				},
			},
			expected: expected{
				value: map[string][]item{
					"odd": {
						{group: "odd", value: 1},
						{group: "odd", value: 3},
					},
					"even": {
						{group: "even", value: 2},
					},
				},
			},
		},
		{
			name: "empty",
			args: args{
				in: []item{},
				keyFunc: func(value item) string {
					return value.group
				},
			},
			expected: expected{
				value: map[string][]item{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := GroupBy(tt.args.in, tt.args.keyFunc)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

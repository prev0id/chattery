package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {
	t.Parallel()

	type args struct {
		values []int
	}
	type expected struct {
		value Set[int]
	}
	tests := []struct {
		expected expected
		name     string
		args     args
	}{
		{
			name: "empty",
			args: args{
				values: nil,
			},
			expected: expected{
				value: Set[int]{},
			},
		},
		{
			name: "deduplicates_values",
			args: args{
				values: []int{1, 2, 2},
			},
			expected: expected{
				value: Set[int]{
					1: {},
					2: {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := NewSet(tt.args.values...)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestSet_Add(t *testing.T) {
	t.Parallel()

	type fields struct {
		set Set[int]
	}
	type args struct {
		value int
	}
	type expected struct {
		set Set[int]
	}
	tests := []struct {
		fields   fields
		expected expected
		name     string
		args     args
	}{
		{
			name: "adds_value",
			fields: fields{
				set: NewSet(1),
			},
			args: args{
				value: 2,
			},
			expected: expected{
				set: NewSet(1, 2),
			},
		},
		{
			name: "keeps_existing_value",
			fields: fields{
				set: NewSet(1),
			},
			args: args{
				value: 1,
			},
			expected: expected{
				set: NewSet(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.fields.set.Add(tt.args.value)

			assert.Equal(t, tt.expected.set, tt.fields.set)
		})
	}
}

func TestSet_Delete(t *testing.T) {
	t.Parallel()

	type fields struct {
		set Set[int]
	}
	type args struct {
		value int
	}
	type expected struct {
		set Set[int]
	}
	tests := []struct {
		fields   fields
		expected expected
		name     string
		args     args
	}{
		{
			name: "deletes_value",
			fields: fields{
				set: NewSet(1, 2),
			},
			args: args{
				value: 2,
			},
			expected: expected{
				set: NewSet(1),
			},
		},
		{
			name: "ignores_missing_value",
			fields: fields{
				set: NewSet(1),
			},
			args: args{
				value: 2,
			},
			expected: expected{
				set: NewSet(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.fields.set.Delete(tt.args.value)

			assert.Equal(t, tt.expected.set, tt.fields.set)
		})
	}
}

func TestSet_Contains(t *testing.T) {
	t.Parallel()

	type fields struct {
		set Set[int]
	}
	type args struct {
		value int
	}
	type expected struct {
		value bool
	}
	tests := []struct {
		fields   fields
		name     string
		args     args
		expected expected
	}{
		{
			name: "contains_value",
			fields: fields{
				set: NewSet(1, 2),
			},
			args: args{
				value: 2,
			},
			expected: expected{
				value: true,
			},
		},
		{
			name: "missing_value",
			fields: fields{
				set: NewSet(1),
			},
			args: args{
				value: 2,
			},
			expected: expected{
				value: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := tt.fields.set.Contains(tt.args.value)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestSet_Empty(t *testing.T) {
	t.Parallel()

	type fields struct {
		set Set[int]
	}
	type expected struct {
		value bool
	}
	tests := []struct {
		fields   fields
		name     string
		expected expected
	}{
		{
			name: "empty",
			fields: fields{
				set: NewSet[int](),
			},
			expected: expected{
				value: true,
			},
		},
		{
			name: "not_empty",
			fields: fields{
				set: NewSet(1),
			},
			expected: expected{
				value: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := tt.fields.set.Empty()

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

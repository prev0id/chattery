package database

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsConstraintViolation(t *testing.T) {
	t.Parallel()

	const constraintName = "users_login_key"

	wrappedErr := fmt.Errorf("wrapped: %w", &pgconn.PgError{ConstraintName: constraintName})

	type args struct {
		err            error
		constraintName string
	}
	type expected struct {
		value bool
	}
	tests := []struct {
		args     args
		name     string
		expected expected
	}{
		{
			name: "matching_constraint",
			args: args{
				err:            &pgconn.PgError{ConstraintName: constraintName},
				constraintName: constraintName,
			},
			expected: expected{
				value: true,
			},
		},
		{
			name: "wrapped_matching_constraint",
			args: args{
				err:            wrappedErr,
				constraintName: constraintName,
			},
			expected: expected{
				value: true,
			},
		},
		{
			name: "different_constraint",
			args: args{
				err:            &pgconn.PgError{ConstraintName: "other_constraint"},
				constraintName: constraintName,
			},
			expected: expected{
				value: false,
			},
		},
		{
			name: "not_pg_error",
			args: args{
				err:            assert.AnError,
				constraintName: constraintName,
			},
			expected: expected{
				value: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := IsConstraintViolation(tt.args.err, tt.args.constraintName)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	wrappedErr := fmt.Errorf("wrapped: %w", pgx.ErrNoRows)

	type args struct {
		err error
	}
	type expected struct {
		value bool
	}
	tests := []struct {
		args     args
		name     string
		expected expected
	}{
		{
			name: "nil",
			expected: expected{
				value: false,
			},
		},
		{
			name: "pgx_no_rows",
			args: args{
				err: pgx.ErrNoRows,
			},
			expected: expected{
				value: true,
			},
		},
		{
			name: "wrapped_pgx_no_rows",
			args: args{
				err: wrappedErr,
			},
			expected: expected{
				value: true,
			},
		},
		{
			name: "other_error",
			args: args{
				err: assert.AnError,
			},
			expected: expected{
				value: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := NotFound(tt.args.err)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

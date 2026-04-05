package compare

import (
	"cmp"

	"chattery/internal/domain"
)

func Servers(lhs, rhs *domain.Server) int {
	return cmp.Or(
		lhs.JoinedAt.Compare(rhs.JoinedAt),
		cmp.Compare(lhs.ID, rhs.ID),
	)
}

func Topics(lhs, rhs *domain.Topic) int {
	return cmp.Or(
		lhs.CreatedAt.Compare(rhs.CreatedAt),
		cmp.Compare(lhs.ID, rhs.ID),
	)
}

func DMs(lhs, rhs *domain.DM) int {
	return cmp.Or(
		lhs.LastActivityAt.Compare(rhs.LastActivityAt),
		cmp.Compare(lhs.ID, rhs.ID),
	)
}

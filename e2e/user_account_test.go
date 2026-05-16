//go:build e2e

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"chattery/e2e/client"
	user_api "chattery/internal/api/user"
)

func Test_GetMe(t *testing.T) {
	t.Parallel()

	var sessionA *http.Cookie
	var sessionB *http.Cookie
	var closedSession *http.Cookie
	var userA *user_api.PostCreateUserRequest
	var userB *user_api.PostCreateUserRequest

	t.Run("create_users", func(t *testing.T) {
		userA, sessionA = createUser(t, "mea")
		userB, sessionB = createUser(t, "meb")
	})

	require.NotNil(t, sessionA)
	cleanupCreatedUser(t, sessionA)
	require.NotNil(t, sessionB)
	cleanupCreatedUser(t, sessionB)

	t.Run("reject_without_session", func(t *testing.T) {
		response := client.MustGetMe(t)
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("get_first_user", func(t *testing.T) {
		requireMe(t, sessionA, userA.Username, userA.Login)
	})

	t.Run("get_second_user", func(t *testing.T) {
		requireMe(t, sessionB, userB.Username, userB.Login)
	})

	t.Run("close_seeded_user_session", func(t *testing.T) {
		closedSession = loginUser(t, "charlie@example.com", seededPassword)

		response := client.MustPostLogout(t, closedSession)
		response.RequireStatus(t, http.StatusSeeOther)
		response.RequireClearedSessionCookie(t)
		require.Equal(t, "/login", response.Header.Get("Location"))
	})

	t.Run("reject_closed_session", func(t *testing.T) {
		require.NotNil(t, closedSession)

		response := client.MustGetMe(t, closedSession)
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireClearedSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("first_session_still_works", func(t *testing.T) {
		requireMe(t, sessionA, userA.Username, userA.Login)
	})
}

func Test_UpdateUser(t *testing.T) {
	t.Parallel()

	var session *http.Cookie
	var duplicateSession *http.Cookie
	var user *user_api.PostCreateUserRequest
	var duplicate *user_api.PostCreateUserRequest
	oldLogin := ""
	updatedUsername := uniqueCreateUser(t, "updname").Username
	updatedLogin := uniqueCreateUser(t, "updlogin").Login
	updatedPassword := "NewPassword123"

	t.Run("create_users", func(t *testing.T) {
		user, session = createUser(t, "update")
		duplicate, duplicateSession = createUser(t, "upddup")
	})

	require.NotNil(t, session)
	cleanupCreatedUser(t, session)
	require.NotNil(t, duplicateSession)
	cleanupCreatedUser(t, duplicateSession)

	t.Run("reject_without_session", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			Username: updatedUsername,
		})
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("update_username_without_current_password", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			Username: updatedUsername,
		}, session)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)

		user.Username = updatedUsername
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("reject_invalid_username", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			Username: "bad username!",
		}, session)
		response.RequireStatus(t, http.StatusBadRequest)
		response.RequireErrorContains(t, "username can only contain")
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("reject_invalid_login", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			Login:           "not-email",
			CurrentPassword: validPassword,
		}, session)
		response.RequireStatus(t, http.StatusBadRequest)
		response.RequireErrorContains(t, "login must be a valid email address")
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("reject_invalid_password", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			CurrentPassword: validPassword,
			Password:        seededPassword,
		}, session)
		response.RequireStatus(t, http.StatusBadRequest)
		response.RequireErrorContains(t, "password must contain at least one uppercase letter")
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("reject_duplicate_username", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			Username: duplicate.Username,
		}, session)
		response.RequireStatus(t, http.StatusConflict)
		response.RequireErrorContains(t, "already exists")
		response.RequireErrorContains(t, duplicate.Username)
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("reject_duplicate_login", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			Login:           duplicate.Login,
			CurrentPassword: validPassword,
		}, session)
		response.RequireStatus(t, http.StatusConflict)
		response.RequireErrorContains(t, "already exists")
		response.RequireErrorContains(t, duplicate.Login)
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("reject_login_change_without_current_password", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			Login: updatedLogin,
		}, session)
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireErrorContains(t, "current password is invalid")
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("reject_password_change_with_wrong_current_password", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			CurrentPassword: "WrongPassword123",
			Password:        updatedPassword,
		}, session)
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireErrorContains(t, "current password is invalid")
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("update_login", func(t *testing.T) {
		oldLogin = user.Login

		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			Login:           updatedLogin,
			CurrentPassword: validPassword,
		}, session)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)

		user.Login = updatedLogin
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("old_login_stops_working", func(t *testing.T) {
		require.NotEmpty(t, oldLogin)

		response := client.MustPostLogin(t, &user_api.PostLoginRequest{
			Login:    oldLogin,
			Password: validPassword,
		})
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireNoSessionCookie(t)
	})

	t.Run("new_login_works", func(t *testing.T) {
		newSession := loginUser(t, user.Login, validPassword)
		cleanupSession(t, newSession)
		requireMe(t, newSession, user.Username, user.Login)
	})

	t.Run("update_password", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			CurrentPassword: validPassword,
			Password:        updatedPassword,
		}, session)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("old_password_stops_working", func(t *testing.T) {
		response := client.MustPostLogin(t, &user_api.PostLoginRequest{
			Login:    user.Login,
			Password: validPassword,
		})
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "invalid password")
	})

	t.Run("new_password_works", func(t *testing.T) {
		newSession := loginUser(t, user.Login, updatedPassword)
		cleanupSession(t, newSession)
		requireMe(t, newSession, user.Username, user.Login)
	})
}

func Test_UpdateUserAccessIsolation(t *testing.T) {
	t.Parallel()

	var owner *user_api.PostCreateUserRequest
	var ownerSession *http.Cookie
	var other *user_api.PostCreateUserRequest
	var otherSession *http.Cookie
	otherUsername := uniqueCreateUser(t, "othername").Username

	t.Run("create_users", func(t *testing.T) {
		owner, ownerSession = createUser(t, "owner")
		other, otherSession = createUser(t, "other")
	})

	require.NotNil(t, ownerSession)
	cleanupCreatedUser(t, ownerSession)
	require.NotNil(t, otherSession)
	cleanupCreatedUser(t, otherSession)

	t.Run("other_session_updates_only_other_user", func(t *testing.T) {
		response := client.MustPutUpdateUser(t, &user_api.PostUpdateUserRequest{
			Username: otherUsername,
		}, otherSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)

		other.Username = otherUsername
	})

	t.Run("owner_remains_unchanged", func(t *testing.T) {
		requireMe(t, ownerSession, owner.Username, owner.Login)
	})

	t.Run("other_is_updated", func(t *testing.T) {
		requireMe(t, otherSession, other.Username, other.Login)
	})
}

func Test_DeleteMe(t *testing.T) {
	t.Parallel()

	var owner *user_api.PostCreateUserRequest
	var ownerSession *http.Cookie
	var other *user_api.PostCreateUserRequest
	var otherSession *http.Cookie

	t.Run("create_users", func(t *testing.T) {
		owner, ownerSession = createUser(t, "delown")
		other, otherSession = createUser(t, "deloth")
	})

	require.NotNil(t, ownerSession)
	cleanupCreatedUser(t, ownerSession)
	require.NotNil(t, otherSession)
	cleanupCreatedUser(t, otherSession)

	t.Run("reject_without_session", func(t *testing.T) {
		response := client.MustDeleteMe(t)
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("other_session_deletes_only_other_user", func(t *testing.T) {
		response := client.MustDeleteMe(t, otherSession)
		response.RequireStatus(t, http.StatusOK)
		response.RequireClearedSessionCookie(t)
		require.Empty(t, response.Body)
	})

	t.Run("other_session_stops_working", func(t *testing.T) {
		response := client.MustGetMe(t, otherSession)
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireClearedSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("owner_session_still_works", func(t *testing.T) {
		requireMe(t, ownerSession, owner.Username, owner.Login)
	})

	t.Run("deleted_user_cannot_login", func(t *testing.T) {
		response := client.MustPostLogin(t, &user_api.PostLoginRequest{
			Login:    other.Login,
			Password: validPassword,
		})
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "not found")
	})
}

//go:build e2e

package e2e

import (
	"bytes"
	stdimage "image"
	"image/color"
	"image/jpeg"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"chattery/e2e/client"
	image_api "chattery/internal/api/image"
	user_api "chattery/internal/api/user"
)

func Test_ImageLifecycle(t *testing.T) {
	t.Parallel()

	var (
		user           *user_api.PostCreateUserRequest
		session        *http.Cookie
		defaultAvatar  []byte
		uploadedAvatar []byte
	)
	uploadedImage := testJPEG(t, color.RGBA{R: 230, G: 40, B: 40, A: 255})

	t.Run("create_user", func(t *testing.T) {
		user, session = createUser(t, "imglife")
		waitUserCacheSync(t)
	})

	require.NotNil(t, user)
	require.NotNil(t, session)
	cleanupCreatedUser(t, session)

	t.Run("get_default_avatar", func(t *testing.T) {
		response := client.MustGetImage(t, user.Username)
		defaultAvatar = requireJPEGResponse(t, response)
	})
	require.NotEmpty(t, defaultAvatar)

	t.Run("reject_upload_without_session", func(t *testing.T) {
		response := client.MustPostUploadImage(t, "avatar.jpeg", uploadedImage)
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("reject_delete_without_session", func(t *testing.T) {
		response := client.MustDeleteImage(t)
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("upload_image", func(t *testing.T) {
		response := client.MustPostUploadImage(t, "avatar.jpeg", uploadedImage, session)
		response.RequireStatus(t, http.StatusOK)

		var body image_api.PostUploadImageResponse
		response.RequireJSON(t, &body)
		require.Equal(t, "/v1/image/"+user.Username+".jpeg", body.Avatar)
	})
	cleanupUserImage(t, session)

	t.Run("get_uploaded_avatar", func(t *testing.T) {
		response := client.MustGetImage(t, user.Username)
		uploadedAvatar = requireJPEGResponse(t, response)
		require.NotEqual(t, defaultAvatar, uploadedAvatar)
	})
	require.NotEmpty(t, uploadedAvatar)

	t.Run("reject_invalid_image", func(t *testing.T) {
		response := client.MustPostUploadImage(t, "invalid.txt", []byte("not an image"), session)
		response.RequireStatus(t, http.StatusBadRequest)
		response.RequireErrorContains(t, "invalid image")
	})

	t.Run("delete_image", func(t *testing.T) {
		response := client.MustDeleteImage(t, session)
		response.RequireStatus(t, http.StatusOK)

		var body image_api.DeleteImageResponse
		response.RequireJSON(t, &body)
		require.Equal(t, "/v1/image/"+user.Username+".jpeg", body.Avatar)
	})

	t.Run("get_default_avatar_after_delete", func(t *testing.T) {
		response := client.MustGetImage(t, user.Username)
		resetAvatar := requireJPEGResponse(t, response)
		require.Equal(t, defaultAvatar, resetAvatar)
		require.NotEqual(t, uploadedAvatar, resetAvatar)
	})
}

func Test_ImageAccessIsolation(t *testing.T) {
	t.Parallel()

	var (
		owner        *user_api.PostCreateUserRequest
		other        *user_api.PostCreateUserRequest
		ownerSession *http.Cookie
		otherSession *http.Cookie
		ownerAvatar  []byte
		otherAvatar  []byte
	)
	ownerImage := testJPEG(t, color.RGBA{R: 30, G: 90, B: 220, A: 255})
	otherImage := testJPEG(t, color.RGBA{R: 40, G: 190, B: 80, A: 255})

	t.Run("create_users", func(t *testing.T) {
		owner, ownerSession = createUser(t, "imgown")
		other, otherSession = createUser(t, "imgoth")
		waitUserCacheSync(t)
	})

	require.NotNil(t, owner)
	require.NotNil(t, ownerSession)
	cleanupCreatedUser(t, ownerSession)
	require.NotNil(t, other)
	require.NotNil(t, otherSession)
	cleanupCreatedUser(t, otherSession)

	t.Run("upload_owner_image", func(t *testing.T) {
		response := client.MustPostUploadImage(t, "owner.jpeg", ownerImage, ownerSession)
		response.RequireStatus(t, http.StatusOK)
	})

	cleanupUserImage(t, ownerSession)

	t.Run("get_owner_image", func(t *testing.T) {
		response := client.MustGetImage(t, owner.Username)
		ownerAvatar = requireJPEGResponse(t, response)
	})
	require.NotEmpty(t, ownerAvatar)

	t.Run("upload_other_image", func(t *testing.T) {
		response := client.MustPostUploadImage(t, "other.jpeg", otherImage, otherSession)
		response.RequireStatus(t, http.StatusOK)
	})

	cleanupUserImage(t, otherSession)

	t.Run("get_other_image", func(t *testing.T) {
		response := client.MustGetImage(t, other.Username)
		otherAvatar = requireJPEGResponse(t, response)
		require.NotEqual(t, ownerAvatar, otherAvatar)
	})

	require.NotEmpty(t, otherAvatar)

	t.Run("other_upload_does_not_change_owner_image", func(t *testing.T) {
		response := client.MustGetImage(t, owner.Username)
		require.Equal(t, ownerAvatar, requireJPEGResponse(t, response))
	})

	t.Run("delete_other_image", func(t *testing.T) {
		response := client.MustDeleteImage(t, otherSession)
		response.RequireStatus(t, http.StatusOK)
	})

	t.Run("other_delete_does_not_change_owner_image", func(t *testing.T) {
		response := client.MustGetImage(t, owner.Username)
		require.Equal(t, ownerAvatar, requireJPEGResponse(t, response))
	})

	t.Run("other_image_is_reset", func(t *testing.T) {
		response := client.MustGetImage(t, other.Username)
		require.NotEqual(t, otherAvatar, requireJPEGResponse(t, response))
	})

	t.Run("reject_unknown_user_image", func(t *testing.T) {
		username := uniqueCreateUser(t, "imgnone").Username
		response := client.MustGetImage(t, username)
		response.RequireStatus(t, http.StatusNotFound)
	})
}

func cleanupUserImage(t *testing.T, session *http.Cookie) {
	t.Helper()

	t.Cleanup(func() {
		response := client.MustDeleteImage(t, session)
		if response.StatusCode == http.StatusUnauthorized {
			return
		}

		response.RequireStatus(t, http.StatusOK)
	})
}

func testJPEG(t testing.TB, fill color.RGBA) []byte {
	t.Helper()

	img := stdimage.NewRGBA(stdimage.Rect(0, 0, 8, 8))
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.SetRGBA(x, y, fill)
		}
	}

	var body bytes.Buffer
	require.NoError(t, jpeg.Encode(&body, img, &jpeg.Options{Quality: 90}))

	return body.Bytes()
}

func requireJPEGResponse(t testing.TB, response *client.Response) []byte {
	t.Helper()

	response.RequireStatus(t, http.StatusOK)
	require.Equal(t, "no-store", response.Header.Get("Cache-Control"))
	require.Equal(t, "image/jpeg", response.Header.Get("Content-Type"))
	require.NotEmpty(t, response.Body)

	_, err := jpeg.Decode(bytes.NewReader(response.Body))
	require.NoError(t, err)

	return response.Body
}

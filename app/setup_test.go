package app_test

import (
	"fmt"
	"instagram-lite/app"
	"instagram-lite/entity"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	t.Run("should make new account when NewAccount is called", func(t *testing.T) {
		acc := app.NewAccount(&entity.User{Name: "aditbuddy"})

		assert.NotEmpty(t, acc)
	})

	t.Run("should return error when new account has the same username with an existing account", func(t *testing.T) {
		ar := app.NewAccRegistry()
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		_, _ = ar.Record(acc1)
		_, err := ar.Record(acc2)

		assert.ErrorIs(t, err, app.ErrUserExist)
	})

	t.Run("should return photo status when GetPhoto is called", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		result := acc1.GetPhoto()

		assert.NotEmpty(t, result)
	})

	t.Run("should show account in follow list when following", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test"})
		expected1 := acc2
		expected2 := acc1

		result1, result2, err := acc1.Follow(acc2)

		assert.Equal(t, expected1, result1[0])
		assert.Equal(t, expected2, result2[0])
		assert.Nil(t, err)
	})

	t.Run("should return error when user follow themselves", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		_, _, err := acc1.Follow(acc2)

		assert.ErrorIs(t, err, app.ErrSameAccount)
	})

	t.Run("should return error when user follow another user that has been followed", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test"})

		_, _, _ = acc1.Follow(acc2)
		_, _, err := acc1.Follow(acc2)

		assert.ErrorIs(t, err, app.ErrAlreadyFollowed)
	})

	t.Run("should upload photo when user post a photo", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		expected1 := app.NewActivity(acc1, "upload", acc1)
		expected2 := true

		result1, err := acc1.Post()
		result2 := acc1.HasUploadPhoto()

		assert.Equal(t, expected1, result1[0])
		assert.Equal(t, expected2, result2)
		assert.Nil(t, err)
	})

	t.Run("should return error when user post a photo twice", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		_, _ = acc1.Post()
		_, err := acc1.Post()

		assert.ErrorIs(t, err, app.ErrUploadTwice)
	})

	t.Run("should get notified when following account post a photo", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test1"})
		acc3 := app.NewAccount(&entity.User{Name: "test2"})
		expected1 := app.NewActivity(acc1, "upload", acc1)
		expected2 := app.NewActivity(acc1, "upload", acc2)
		expected3 := app.NewActivity(acc1, "upload", acc3)

		_, _, _ = acc2.Follow(acc1)
		_, _, _ = acc3.Follow(acc1)
		result1, _ := acc1.Post()
		result2 := acc2.GetActivity()
		result3 := acc3.GetActivity()

		assert.Equal(t, expected1, result1[0])
		assert.Equal(t, expected2, result2[0])
		assert.Equal(t, expected3, result3[0])
	})

	t.Run("should like a photo when user like another user photo", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test"})
		expected := app.NewActivity(acc1, "like", acc2)

		_, _ = acc2.Post()
		_, _, _ = acc1.Follow(acc2)
		result1, result2, err := acc1.Like(acc2)

		assert.Equal(t, expected, result1[0])
		assert.Equal(t, expected, result2[1])
		assert.Nil(t, err)
	})

	t.Run("should return error when user like same account with no photo uploaded", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		_, _, err := acc1.Like(acc1)

		assert.ErrorIs(t, err, app.ErrNoPhoto)
	})

	t.Run("should return error when user like same account twice", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		_, _ = acc1.Post()
		_, _, _ = acc1.Like(acc1)
		_, _, err := acc1.Like(acc1)

		assert.ErrorIs(t, err, app.ErrLikedTwice)
	})

	t.Run("should return error when user like another account wihtout following", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test"})

		_, _, err := acc1.Like(acc2)

		assert.Equal(t, fmt.Errorf("unable to like %s's photo", acc2.GetUsername()), err)
	})

	t.Run("should return error when user like another account with no photo uploaded", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test"})

		_, _, _ = acc1.Follow(acc2)
		_, _, err := acc1.Like(acc2)

		assert.Equal(t, fmt.Errorf("%s doesn't have a photo", acc2.GetUsername()), err)
	})

	t.Run("should return error when user like another account twice", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test"})

		_, _, _ = acc1.Follow(acc2)
		_, _ = acc2.Post()
		_, _, _ = acc1.Like(acc2)
		_, _, err := acc1.Like(acc2)

		assert.ErrorIs(t, err, app.ErrLikedTwice)
	})

	t.Run("should get notified when following account like a photo", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test1"})
		acc3 := app.NewAccount(&entity.User{Name: "test2"})
		expected1 := app.NewActivity(acc1, "upload", acc1)
		expected2 := true
		expected3 := app.NewActivity(acc3, "like", acc1)

		_, _, _ = acc2.Follow(acc3)
		_, _, _ = acc3.Follow(acc1)
		result1, _ := acc1.Post()
		result3, _, _ := acc3.Like(acc1)
		act2 := acc2.GetActivity()
		result2 := act2[0].IsSameActivity(app.NewActivity(acc3, "like", acc1))

		assert.Equal(t, expected1, result1[0])
		assert.Equal(t, expected2, result2)
		assert.Equal(t, expected3, result3[1])
	})

	t.Run("should filter like photo notification when two account follow each other", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test1"})
		acc3 := app.NewAccount(&entity.User{Name: "test2"})
		expected1 := app.NewActivity(acc2, "like", acc1)
		expected2 := app.NewActivity(acc2, "like", acc1)
		expected3 := app.NewActivity(acc2, "like", acc3)

		_, _, _ = acc1.Follow(acc2)
		_, _, _ = acc2.Follow(acc1)
		_, _ = acc1.Post()
		_, _ = acc3.Post()
		_, _, _ = acc2.Follow(acc3)
		result2, result1, _ := acc2.Like(acc1)
		result3, _, _ := acc2.Like(acc3)

		assert.Equal(t, expected1, result1[1])
		assert.Equal(t, expected2, result2[1])
		assert.Equal(t, expected3, result3[2])
	})

	t.Run("should filter like photo notification when two account follow each other", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test1"})
		acc3 := app.NewAccount(&entity.User{Name: "test2"})
		expected := true

		_, _, _ = acc2.Follow(acc1)
		_, _, _ = acc3.Follow(acc1)
		_, _, _ = acc1.Follow(acc2)
		_, _ = acc1.Post()
		_, _ = acc2.Post()
		_, _, _ = acc2.Like(acc1)
		_, _, _ = acc3.Like(acc1)
		_, _, _ = acc1.Like(acc2)
		result := acc1.HasMorePhotoLike(acc2)

		assert.Equal(t, expected, result)
	})
}

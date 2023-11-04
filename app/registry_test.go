package app_test

import (
	"instagram-lite/app"
	"instagram-lite/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegitstry(t *testing.T) {
	t.Run("should make account registry when NewAccRegistry funtion is called", func(t *testing.T) {
		r := app.NewAccRegistry()

		assert.NotEmpty(t, r)
	})

	t.Run("should assign new account to account registry", func(t *testing.T) {
		r := app.NewAccRegistry()
		acc := app.NewAccount(&entity.User{Name: "aditbuddy"})
		expected := acc

		result, err := r.Record(acc)

		assert.Equal(t, expected, result[0])
		assert.Nil(t, err)
	})

	t.Run("should return user exist if username already used", func(t *testing.T) {
		r := app.NewAccRegistry()
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		_, _ = r.Record(acc1)
		_, err := r.Record(acc2)

		assert.ErrorIs(t, err, app.ErrUserExist)
	})

	t.Run("should return false if username unknown", func(t *testing.T) {
		r := app.NewAccRegistry()
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test"})
		expected := false

		_, _ = r.Record(acc1)
		_, result := r.IsAccountExist(acc2)

		assert.Equal(t, expected, result)
	})

	t.Run("should return sorted account by photo like", func(t *testing.T) {
		r := app.NewAccRegistry()
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})
		acc2 := app.NewAccount(&entity.User{Name: "test"})
		expected := acc1

		_, _ = r.Record(acc1)
		_, _ = r.Record(acc2)
		_, _, _ = acc1.Follow(acc2)
		_, _, _ = acc2.Follow(acc1)
		_, _ = acc1.Post()
		_, _ = acc2.Post()
		_, _, _ = acc2.Like(acc1)
		result := r.GetLeaderboard()

		assert.Equal(t, expected, result[0])
	})
}

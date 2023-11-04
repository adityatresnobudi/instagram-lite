package app_test

import (
	"instagram-lite/app"
	"instagram-lite/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivity(t *testing.T) {
	t.Run("should make new activity when NewActivity funtion is called", func(t *testing.T) {
		acc := app.NewAccount(&entity.User{Name: "aditbuddy"})
		act := app.NewActivity(acc, "upload", acc)

		assert.NotEmpty(t, act)
	})

	t.Run("should return activity accDo status when GetAccDo is called", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		_, _ = acc1.Post()
		act := acc1.GetActivity()
		result := act[0].GetAccDo()

		assert.NotEmpty(t, result)
	})

	t.Run("should return activity action status when GetAction is called", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		_, _ = acc1.Post()
		act := acc1.GetActivity()
		result := act[0].GetAction()

		assert.NotEmpty(t, result)
	})

	t.Run("should return activity accTo status when GetAccTo is called", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "aditbuddy"})

		_, _ = acc1.Post()
		act := acc1.GetActivity()
		result := act[0].GetAccTo()

		assert.NotEmpty(t, result)
	})
}

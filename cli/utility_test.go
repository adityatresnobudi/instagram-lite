package cli_test

import (
	"fmt"
	"instagram-lite/app"
	"instagram-lite/cli"
	"instagram-lite/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtility(t *testing.T) {
	t.Run("should return error when input is empty", func(t *testing.T) {
		relation := ""

		_, _, _, err := cli.HandleSetup(relation)

		assert.ErrorIs(t, err, cli.ErrInvalidInput)
	})

	t.Run("should return error when input doesn't contain the word 'follows'", func(t *testing.T) {
		relation := "adit talk budi"

		_, _, _, err := cli.HandleSetup(relation)

		assert.ErrorIs(t, err, cli.ErrInvalidKeyword)
	})

	t.Run("should return error when input is less or more than 3 words", func(t *testing.T) {
		relation := "adit budi"

		_, _, _, err := cli.HandleSetup(relation)

		assert.ErrorIs(t, err, cli.ErrInvalidInput)
	})

	t.Run("should return error when input is less or more than 3 words", func(t *testing.T) {
		relation := "adit budi follow adit"

		_, _, _, err := cli.HandleSetup(relation)

		assert.ErrorIs(t, err, cli.ErrInvalidInput)
	})

	t.Run("should return error when the second word is not 'follows'", func(t *testing.T) {
		relation := "adit follow budi"

		_, _, _, err := cli.HandleSetup(relation)

		assert.ErrorIs(t, err, cli.ErrInvalidKeyword)
	})

	t.Run("should properly follow and record when HandleSetup is called", func(t *testing.T) {
		relation1 := "Alice follows Bob"
		acc1 := app.NewAccount(&entity.User{Name: "Alice"})
		acc2 := app.NewAccount(&entity.User{Name: "Bob"})
		expected := true

		res, _, _, _ := cli.HandleSetup(relation1)
		result1 := res[0].HasFollow(res[1])
		result2 := res[0].IsSameAccount(acc1)
		result3 := res[1].IsSameAccount(acc2)

		assert.Equal(t, expected, result1)
		assert.Equal(t, expected, result2)
		assert.Equal(t, expected, result3)
	})

	t.Run("should Record account that is not recorded yet when HandleSetup is called", func(t *testing.T) {
		relation2 := "Alice follows Bill"
		acc := app.NewAccount(&entity.User{Name: "Bill"})
		expected := true

		res, _, _, _ := cli.HandleSetup(relation2)
		result := res[2].IsSameAccount(acc)

		assert.Equal(t, expected, result)
	})

	t.Run("should return error when action is empty", func(t *testing.T) {
		action := ""

		_, _, err := cli.HandleAction(action)

		assert.ErrorIs(t, err, cli.ErrInvalidInput)
	})

	t.Run("should return error when input has keyword likes and uploaded", func(t *testing.T) {
		action := "adit likes uploaded photo"

		_, _, err := cli.HandleAction(action)

		assert.ErrorIs(t, err, cli.ErrInvalidKeyword)
	})

	t.Run("should return error when like input has a typo", func(t *testing.T) {
		action := "adit like budi photo"

		_, _, err := cli.HandleAction(action)

		assert.ErrorIs(t, err, cli.ErrInvalidKeyword)
	})

	t.Run("should return error when like input has a typo", func(t *testing.T) {
		action := "adit likes budi photos"

		_, _, err := cli.HandleAction(action)

		assert.ErrorIs(t, err, cli.ErrInvalidKeyword)
	})

	t.Run("should return error when input has a nonexistent user", func(t *testing.T) {
		acc := app.NewAccount(&entity.User{Name: "adit"})
		action := "adit likes budi photo"

		_, _, err := cli.HandleAction(action)

		assert.Equal(t, fmt.Errorf("unknown user %s", acc.GetUsername()), err)
	})

	t.Run("should return error when upload input has a typo", func(t *testing.T) {
		action := "adit upload photo"

		_, _, err := cli.HandleAction(action)

		assert.ErrorIs(t, err, cli.ErrInvalidKeyword)
	})

	t.Run("should return error when upload input has a typo", func(t *testing.T) {
		action := "adit uploaded photos"

		_, _, err := cli.HandleAction(action)

		assert.ErrorIs(t, err, cli.ErrInvalidKeyword)
	})

	t.Run("should return error when upload input has a nonexistent user", func(t *testing.T) {
		acc := app.NewAccount(&entity.User{Name: "adit"})
		action := "adit uploaded photo"

		_, _, err := cli.HandleAction(action)

		assert.Equal(t, fmt.Errorf("unknown user %s", acc.GetUsername()), err)
	})

	t.Run("should upload when HandleAction handling action upload", func(t *testing.T) {
		relation3 := "John follows Bob"
		relation4 := "Bob follows Alice"
		relation5 := "Bob follows Bill"
		relation6 := "John follows Alice"
		action1 := "Alice uploaded photo"
		acc := app.NewAccount(&entity.User{Name: "Alice"})
		expected := true

		_, _, _, _ = cli.HandleSetup(relation3)
		_, _, _, _ = cli.HandleSetup(relation4)
		_, _, _, _ = cli.HandleSetup(relation5)
		_, _, _, _ = cli.HandleSetup(relation6)
		_, res, _ := cli.HandleAction(action1)
		result := res[0].IsSameActivity(app.NewActivity(acc, "upload", acc))

		assert.Equal(t, expected, result)
	})

	t.Run("should like when HandleAction handling action like", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "Bob"})
		acc2 := app.NewAccount(&entity.User{Name: "Alice"})
		action2 := "Bob likes Alice photo"
		expected := true

		res1, res2, _ := cli.HandleAction(action2)
		result1 := res1[1].IsSameActivity(app.NewActivity(acc1, "like", acc2))
		result2 := res2[1].IsSameActivity(app.NewActivity(acc1, "like", acc2))

		assert.Equal(t, expected, result1)
		assert.Equal(t, expected, result2)
	})

	t.Run("should return error when display is empty", func(t *testing.T) {
		display := ""

		_, err := cli.HandleDisplay(display)

		assert.ErrorIs(t, err, cli.ErrInvalidInput)
	})

	t.Run("should return error when display input is a nonexistent user", func(t *testing.T) {
		acc := app.NewAccount(&entity.User{Name: "adit"})
		display := "adit"

		_, err := cli.HandleDisplay(display)

		assert.Equal(t, fmt.Errorf("unknown user %s", acc.GetUsername()), err)
	})

	t.Run("should display activity when HandleDisplay is called", func(t *testing.T) {
		action3 := "Bill uploaded photo"
		action4 := "Bob likes Bill photo"
		action5 := "Bill likes Bill photo"
		action6 := "Alice likes Bill photo"
		acc1 := app.NewAccount(&entity.User{Name: "Bob"})
		expected := "\nBob activities:\n" +
			"Alice uploaded photo\n" +
			"You liked Alice's photo\n" +
			"Bill uploaded photo\n" +
			"You liked Bill's photo\n" +
			"Bill liked Bill's photo\n" +
			"Alice liked Bill's photo\n"

		_, _, _ = cli.HandleAction(action3)
		_, _, _ = cli.HandleAction(action4)
		_, _, _ = cli.HandleAction(action5)
		_, _, _ = cli.HandleAction(action6)
		result, _ := cli.HandleDisplay(acc1.GetUsername())

		assert.Equal(t, expected, result)
	})

	t.Run("should display activity when HandleDisplay is called", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "Alice"})
		expected := "\nAlice activities:\n" +
			"You uploaded photo\n" +
			"Bob liked your photo\n" +
			"Bill uploaded photo\n" +
			"Bob liked Bill's photo\n" +
			"Bill liked Bill's photo\n" +
			"You liked Bill's photo\n"

		result, _ := cli.HandleDisplay(acc1.GetUsername())

		assert.Equal(t, expected, result)
	})

	t.Run("should display activity when HandleDisplay is called", func(t *testing.T) {
		acc1 := app.NewAccount(&entity.User{Name: "Bill"})
		expected := "\nBill activities:\n" +
			"You uploaded photo\n" +
			"Bob liked your photo\n" +
			"You liked your photo\n" +
			"Alice liked your photo\n"

		result, _ := cli.HandleDisplay(acc1.GetUsername())

		assert.Equal(t, expected, result)
	})

	t.Run("should return trending photo leaderboard when HandleTrending is called", func(t *testing.T) {
		expected := "Trending photos:\n" +
			"1. Bill photo got 3 likes\n" +
			"2. Alice photo got 1 likes\n"

		result := cli.HandleTrending()

		assert.Equal(t, expected, result)
	})
}

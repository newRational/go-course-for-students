package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser_EmptyNickname(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("", "jenny@gmail.com")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUser_TooLongNickname(t *testing.T) {
	client := getTestClient()

	nick := strings.Repeat("a", 51)

	_, err := client.createUser(nick, "jenny@gmail.com")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUser_EmptyEmail(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("jenny", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUser_TooLongEmail(t *testing.T) {
	client := getTestClient()

	username := strings.Repeat("a", 51)

	_, err := client.createUser("jenny", username+"@gmail.com")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("jenny", "invalid_email")
	fmt.Println()
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser_EmptyNickname(t *testing.T) {
	client := getTestClient()

	resp, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	_, err = client.updateUser(resp.Data.ID, "", "jenny@gmail.com")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser_TooLongNickname(t *testing.T) {
	client := getTestClient()

	resp, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	nick := strings.Repeat("a", 51)

	_, err = client.updateUser(resp.Data.ID, nick, "jenny@gmail.com")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser_EmptyEmail(t *testing.T) {
	client := getTestClient()

	resp, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	_, err = client.updateUser(resp.Data.ID, "jenny", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser_TooLongEmail(t *testing.T) {
	client := getTestClient()

	resp, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	username := strings.Repeat("a", 51)

	_, err = client.updateUser(resp.Data.ID, "jenny", username+"@gmail.com")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser_InvalidEmail(t *testing.T) {
	client := getTestClient()

	resp, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	_, err = client.updateUser(resp.Data.ID, "jenny", "invalid_email")
	assert.ErrorIs(t, err, ErrBadRequest)
}

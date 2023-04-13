package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	// create user
	_, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	// create user
	_, err = client.createUser("polly", "polly@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	resp, err = client.changeAdStatus(1, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	// create user
	_, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	// create user
	_, err = client.createUser("polly", "polly@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(1, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	// create user
	_, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	// added
	_, err = client.createAd(123, "hello", "world")
	assert.ErrorIs(t, err, ErrBadRequest)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}

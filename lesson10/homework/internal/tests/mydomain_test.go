package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAdByNonExistentUser(t *testing.T) {
	client := getTestHTTPClient()

	// create user
	_, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	// added
	_, err = client.createAd(123, "hello", "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestChangeStatusAdOfNonExistentUser(t *testing.T) {
	client := getTestHTTPClient()

	// create user
	_, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	resp, err = client.changeAdStatus(1, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAdOfNonExistentUser(t *testing.T) {
	client := getTestHTTPClient()

	// create user
	_, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	// create user
	_, err = client.createUser("polly", "polly@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(123, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrBadRequest)
}

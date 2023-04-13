package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	resp, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)
	assert.Zero(t, resp.Data.ID)
	assert.Equal(t, resp.Data.Nickname, "jenny")
	assert.Equal(t, resp.Data.Email, "jenny@gmail.com")
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	resp, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	resp, err = client.updateUser(resp.Data.ID, "polly", "polly@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.Nickname, "polly")
	assert.Equal(t, resp.Data.Email, "polly@gmail.com")
}

func TestShowAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	ad0, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	ad1, err := client.createAd(0, "hello", "friend")
	assert.NoError(t, err)

	assert.Zero(t, ad0.Data.ID)
	assert.Equal(t, ad0.Data.Title, "hello")
	assert.Equal(t, ad0.Data.Text, "world")
	assert.Equal(t, ad0.Data.AuthorID, int64(0))
	assert.False(t, ad0.Data.Published)

	assert.Equal(t, ad1.Data.ID, int64(1))
	assert.Equal(t, ad1.Data.Title, "hello")
	assert.Equal(t, ad1.Data.Text, "friend")
	assert.Equal(t, ad1.Data.AuthorID, int64(0))
	assert.False(t, ad1.Data.Published)
}

func TestListAdsWithParams(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("jenny", "jenny@gmail.com")
	assert.NoError(t, err)

	_, err = client.createUser("polly", "polly@gmail.com")
	assert.NoError(t, err)

	ad0, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(1, "hello", "friend")
	assert.NoError(t, err)

	ad1, err := client.createAd(0, "hello", "привет")
	assert.NoError(t, err)

	resp, err := client.createAd(1, "goodbye", "friend")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(resp.Data.AuthorID, resp.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.listAds(map[string]string{"published": "false", "title": "hello", "user_id": "0", "created": "2023-04-14"})
	assert.NoError(t, err)

	fmt.Fprintln(os.Stderr, "ads:", ads.Data)
	assert.Len(t, ads.Data, 1)

	assert.Contains(t, ads.Data, ad0.Data)
	assert.Contains(t, ads.Data, ad1.Data)
}

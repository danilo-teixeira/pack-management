package pack_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pack-management/internal/domain/pack"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePack(t *testing.T) {
	t.Run("Shoud create a pack successfully", func(t *testing.T) {
		defer gock.Off()

		gock.New(dogApiURL).
			Get("/facts").
			MatchParam("limit", "1").
			Reply(http.StatusOK).
			JSON(`{
				"data": [
					{
						"id": "cb382e94-d7e2-415b-b943-085960f3819a",
						"type": "fact",
						"attributes": {
							"body": "Toto in The Wizard of Oz was played by a female Cairn Terrier named Terry."
						}
					}
				]
			}`)

		gock.New(negerDateAPIURL).
			Get("/PublicHolidays/2025/BR").
			Reply(http.StatusOK).
			JSON(`[
				{
					"date": "2025-01-01",
					"localName": "Confraternização Universal",
					"name": "New Year's Day",
					"countryCode": "BR",
					"fixed": false,
					"global": true,
					"counties": null,
					"launchYear": null,
					"types": [
						"Public"
					]
				}
			]`)

		resp, err := clientApp(httptest.NewRequest(
			http.MethodPost,
			"/packs",
			bytes.NewBuffer([]byte(`{
				"description": "Livros para entrega",
				"sender": "Loja ABC",
				"recipient": "João Silva",
				"estimated_delivery_date": "2025-04-02"
			}`)),
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		packJSON := pack.CreatePackResponse{}
		err = json.NewDecoder(resp.Body).Decode(&packJSON)
		assert.Nil(t, err)

		assert.NotEmpty(t, packJSON.ID)
		assert.Equal(t, "Livros para entrega", packJSON.Description)
		assert.Equal(t, "Loja ABC", packJSON.SenderName)
		assert.Equal(t, "João Silva", packJSON.ReceiverName)
		assert.Equal(t, "CREATED", packJSON.Status.String())
		assert.NotEmpty(t, packJSON.CreatedAt)
		assert.NotEmpty(t, packJSON.UpdateAt)

		time.Sleep(1 * time.Second) // wait for the gock to finish
		assert.True(t, gock.IsDone())
	})

	t.Run("Shoud return error when missing required fields", func(t *testing.T) {
		body := []byte(`{}`)

		resp, err := clientApp(httptest.NewRequest(http.MethodPost, "/packs", bytes.NewBuffer(body)))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestUpdatePackStatus(t *testing.T) {
	t.Run("Shoud update a pack status from CREATED to IN_TRANSIT successfully", func(t *testing.T) {
		createdPack := createPack(t)

		resp, err := clientApp(httptest.NewRequest(
			http.MethodPatch,
			"/packs/"+createdPack.ID,
			bytes.NewBuffer([]byte(`{
				"status": "IN_TRANSIT"
			}`)),
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		packJSON := pack.UpdatePackResponse{}
		err = json.NewDecoder(resp.Body).Decode(&packJSON)
		assert.Nil(t, err)

		assert.NotEmpty(t, packJSON.ID)
		assert.Equal(t, createdPack.Description, packJSON.Description)
		assert.Equal(t, createdPack.SenderName, packJSON.SenderName)
		assert.Equal(t, createdPack.ReceiverName, packJSON.ReceiverName)
		assert.Equal(t, "IN_TRANSIT", packJSON.Status.String())
		assert.NotEmpty(t, packJSON.CreatedAt)
		assert.NotEmpty(t, packJSON.UpdateAt)
		assert.Empty(t, packJSON.DeliveredAt)
	})

	t.Run("Shoud update a pack status from IN_TRANSIT to DELIVERED successfully", func(t *testing.T) {
		createdPack := createPack(t)

		resp, err := clientApp(httptest.NewRequest(
			http.MethodPatch,
			"/packs/"+createdPack.ID,
			bytes.NewBuffer([]byte(`{
				"status": "IN_TRANSIT"
			}`)),
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = clientApp(httptest.NewRequest(
			http.MethodPatch,
			"/packs/"+createdPack.ID,
			bytes.NewBuffer([]byte(`{
				"status": "DELIVERED"
			}`)),
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		packJSON := pack.UpdatePackResponse{}
		err = json.NewDecoder(resp.Body).Decode(&packJSON)
		assert.Nil(t, err)

		assert.NotEmpty(t, packJSON.ID)
		assert.Equal(t, createdPack.Description, packJSON.Description)
		assert.Equal(t, createdPack.SenderName, packJSON.SenderName)
		assert.Equal(t, createdPack.ReceiverName, packJSON.ReceiverName)
		assert.Equal(t, "DELIVERED", packJSON.Status.String())
		assert.NotEmpty(t, packJSON.CreatedAt)
		assert.NotEmpty(t, packJSON.UpdateAt)
		assert.NotEmpty(t, packJSON.DeliveredAt)
	})

	t.Run("Shoud return error when skip from CREATED to DELIVERED", func(t *testing.T) {
		createdPack := createPack(t)

		resp, err := clientApp(httptest.NewRequest(
			http.MethodPatch,
			"/packs/"+createdPack.ID,
			bytes.NewBuffer([]byte(`{
				"status": "DELIVERED"
			}`)),
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Shoud return error when status is invalid", func(t *testing.T) {
		resp, err := clientApp(httptest.NewRequest(
			http.MethodPatch,
			"/packs/pack_1",
			bytes.NewBuffer([]byte(`{
				"status": "INVALID_STATUS"
			}`)),
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Shoud return error when pack not found", func(t *testing.T) {
		resp, err := clientApp(httptest.NewRequest(
			http.MethodPatch,
			"/packs/pack_not_found_1",
			bytes.NewBuffer([]byte(`{
				"status": "IN_TRANSIT"
			}`)),
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func createPack(t *testing.T) pack.CreatePackResponse {
	defer gock.Off()

	gock.New(dogApiURL).
		Get("/facts").
		MatchParam("limit", "1").
		Reply(http.StatusOK).
		JSON(`{
			"data": [
				{
					"id": "cb382e94-d7e2-415b-b943-085960f3819a",
					"type": "fact",
					"attributes": {
						"body": "Toto in The Wizard of Oz was played by a female Cairn Terrier named Terry."
					}
				}
			]
		}`)

	gock.New(negerDateAPIURL).
		Get("/PublicHolidays/2025/BR").
		Reply(http.StatusOK).
		JSON(`[
			{
				"date": "2025-01-01",
				"localName": "Confraternização Universal",
				"name": "New Year's Day",
				"countryCode": "BR",
				"fixed": false,
				"global": true,
				"counties": null,
				"launchYear": null,
				"types": [
					"Public"
				]
			}
		]`)

	resp, err := clientApp(httptest.NewRequest(
		http.MethodPost,
		"/packs",
		bytes.NewBuffer([]byte(`{
			"description": "Livros para entrega",
			"sender": "Loja ABC",
			"recipient": "João Silva",
			"estimated_delivery_date": "2025-04-02"
		}`)),
	))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	packJSON := pack.CreatePackResponse{}
	err = json.NewDecoder(resp.Body).Decode(&packJSON)
	assert.Nil(t, err)

	return packJSON
}

package pack_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pack-management/internal/domain/pack"
	"testing"

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

		assert.True(t, gock.IsDone())
	})

	t.Run("Shoud return error when missing required fields", func(t *testing.T) {
		body := []byte(`{}`)

		resp, err := clientApp(httptest.NewRequest(http.MethodPost, "/packs", bytes.NewBuffer(body)))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

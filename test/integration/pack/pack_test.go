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

		packJSON := pack.PackJSON{}
		err = json.NewDecoder(resp.Body).Decode(&packJSON)
		assert.Nil(t, err)

		assert.NotEmpty(t, packJSON.ID)
		assert.Equal(t, "Livros para entrega", packJSON.Description)
		assert.Equal(t, "Loja ABC", packJSON.SenderName)
		assert.Equal(t, "João Silva", packJSON.ReceiverName)
		assert.Equal(t, "CREATED", packJSON.Status.String())
		assert.NotEmpty(t, packJSON.CreatedAt)
		assert.NotEmpty(t, packJSON.UpdateAt)
		assert.Empty(t, packJSON.DeliveredAt)
		assert.Empty(t, packJSON.CanceledAt)
		assert.Empty(t, packJSON.Events)

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

func TestGetPackByID(t *testing.T) {
	t.Run("Shoud get a pack successfully", func(t *testing.T) {
		createdPack := createPack(t, nil)

		resp, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs/"+createdPack.ID,
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		packJSON := pack.PackJSON{}
		err = json.NewDecoder(resp.Body).Decode(&packJSON)
		assert.Nil(t, err)

		assert.NotEmpty(t, packJSON.ID)
		assert.Equal(t, createdPack.Description, packJSON.Description)
		assert.Equal(t, createdPack.SenderName, packJSON.SenderName)
		assert.Equal(t, createdPack.ReceiverName, packJSON.ReceiverName)
		assert.Equal(t, "CREATED", packJSON.Status.String())
		assert.NotEmpty(t, packJSON.CreatedAt)
		assert.NotEmpty(t, packJSON.UpdateAt)
		assert.Empty(t, packJSON.DeliveredAt)
		assert.Empty(t, packJSON.CanceledAt)
		assert.Empty(t, packJSON.Events)
	})

	t.Run("Shoud get a pack successfully with events", func(t *testing.T) {
		createdPack := createPack(t, nil)
		createEvent(t, createdPack.ID)

		resp, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs/"+createdPack.ID+"?with_events=true",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		packJSON := pack.PackJSON{}
		err = json.NewDecoder(resp.Body).Decode(&packJSON)
		assert.Nil(t, err)

		assert.NotEmpty(t, packJSON.ID)
		assert.Equal(t, createdPack.Description, packJSON.Description)
		assert.Equal(t, createdPack.SenderName, packJSON.SenderName)
		assert.Equal(t, createdPack.ReceiverName, packJSON.ReceiverName)
		assert.Equal(t, "CREATED", packJSON.Status.String())
		assert.NotEmpty(t, packJSON.CreatedAt)
		assert.NotEmpty(t, packJSON.UpdateAt)
		assert.Empty(t, packJSON.DeliveredAt)
		assert.Empty(t, packJSON.CanceledAt)

		eventDate, _ := time.Parse(time.RFC3339, "2025-01-20T15:13:59Z")
		assert.Len(t, packJSON.Events, 1)
		assert.NotEmpty(t, packJSON.Events[0].ID)
		assert.Equal(t, "Pacote chegou ao centro de distribuição", packJSON.Events[0].Description)
		assert.Equal(t, "Centro de Distribuição São Paulo", packJSON.Events[0].Location)
		assert.Equal(t, eventDate, packJSON.Events[0].Date)
	})

	t.Run("Shoud return error when pack not found", func(t *testing.T) {
		resp, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs/pack_not_found_1",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestListPacks(t *testing.T) {

	t.Run("Shoud list packs successfully", func(t *testing.T) {
		resp, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respJSON := pack.ListPackJSON{}
		err = json.NewDecoder(resp.Body).Decode(&respJSON)
		assert.Nil(t, err)

		assert.Len(t, respJSON.Items, 3)
		assert.Equal(t, 100, respJSON.Metadata.PageSize)
		assert.Empty(t, respJSON.Metadata.NextCursor)
		assert.Empty(t, respJSON.Metadata.PrevCursor)
	})

	t.Run("Shoud list packs successfully with page_size filter", func(t *testing.T) {
		resp, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs?page_size=1",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respJSON := pack.ListPackJSON{}
		err = json.NewDecoder(resp.Body).Decode(&respJSON)
		assert.Nil(t, err)

		assert.Len(t, respJSON.Items, 1)
		assert.Equal(t, 1, respJSON.Metadata.PageSize)
		assert.NotEmpty(t, respJSON.Metadata.NextCursor)
		assert.Empty(t, respJSON.Metadata.PrevCursor)
	})

	t.Run("Shoud list packs successfully with page_cursor filter", func(t *testing.T) {
		respPage1, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs?page_size=1",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, respPage1.StatusCode)

		respPage1JSON := pack.ListPackJSON{}
		err = json.NewDecoder(respPage1.Body).Decode(&respPage1JSON)
		assert.Nil(t, err)

		assert.Len(t, respPage1JSON.Items, 1)
		assert.Equal(t, 1, respPage1JSON.Metadata.PageSize)
		assert.NotEmpty(t, respPage1JSON.Metadata.NextCursor)
		assert.Empty(t, respPage1JSON.Metadata.PrevCursor)

		respPage2, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs?page_size=1&page_cursor="+(respPage1JSON.Metadata.NextCursor),
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, respPage2.StatusCode)

		respPage2JSON := pack.ListPackJSON{}
		err = json.NewDecoder(respPage2.Body).Decode(&respPage2JSON)
		assert.Nil(t, err)

		assert.Len(t, respPage2JSON.Items, 1)
		assert.Equal(t, 1, respPage2JSON.Metadata.PageSize)
		assert.NotEmpty(t, respPage2JSON.Metadata.NextCursor)
		assert.NotEmpty(t, respPage2JSON.Metadata.PrevCursor)

		respPage3, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs?page_size=1&page_cursor="+(respPage2JSON.Metadata.PrevCursor),
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, respPage3.StatusCode)

		respPage3JSON := pack.ListPackJSON{}
		err = json.NewDecoder(respPage3.Body).Decode(&respPage3JSON)
		assert.Nil(t, err)

		assert.Len(t, respPage3JSON.Items, 1)
		assert.Equal(t, 1, respPage3JSON.Metadata.PageSize)
		assert.Equal(t, respPage1JSON.Metadata.NextCursor, respPage3JSON.Metadata.NextCursor)
		assert.Empty(t, respPage3JSON.Metadata.PrevCursor)
		assert.Equal(t, respPage1JSON.Items[0].ID, respPage3JSON.Items[0].ID)
	})

	t.Run("Shoud list packs successfully with sender_name filter", func(t *testing.T) {
		createPack(t, &createPackParams{SenderName: "test_sender"})

		resp, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs?sender_name=test_sender",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respJSON := pack.ListPackJSON{}
		err = json.NewDecoder(resp.Body).Decode(&respJSON)
		assert.Nil(t, err)

		assert.Len(t, respJSON.Items, 1)
		assert.Equal(t, 100, respJSON.Metadata.PageSize)
		assert.Empty(t, respJSON.Metadata.NextCursor)
		assert.Empty(t, respJSON.Metadata.PrevCursor)
		assert.Equal(t, "test_sender", respJSON.Items[0].SenderName)
	})

	t.Run("Shoud list packs successfully with sender_name filter", func(t *testing.T) {
		createPack(t, &createPackParams{RecipientName: "Recipient_sender"})

		resp, err := clientApp(httptest.NewRequest(
			http.MethodGet,
			"/packs?recipient_name=Recipient_sender",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respJSON := pack.ListPackJSON{}
		err = json.NewDecoder(resp.Body).Decode(&respJSON)
		assert.Nil(t, err)

		assert.Len(t, respJSON.Items, 1)
		assert.Equal(t, 100, respJSON.Metadata.PageSize)
		assert.Empty(t, respJSON.Metadata.NextCursor)
		assert.Empty(t, respJSON.Metadata.PrevCursor)
		assert.Equal(t, "Recipient_sender", respJSON.Items[0].ReceiverName)
	})
}

func TestUpdatePackStatus(t *testing.T) {
	t.Run("Shoud update a pack status from CREATED to IN_TRANSIT successfully", func(t *testing.T) {
		createdPack := createPack(t, nil)

		resp, err := clientApp(httptest.NewRequest(
			http.MethodPatch,
			"/packs/"+createdPack.ID,
			bytes.NewBuffer([]byte(`{
				"status": "IN_TRANSIT"
			}`)),
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		packJSON := pack.PackJSON{}
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
		assert.Empty(t, packJSON.CanceledAt)
		assert.Empty(t, packJSON.Events)
	})

	t.Run("Shoud update a pack status from IN_TRANSIT to DELIVERED successfully", func(t *testing.T) {
		createdPack := createPack(t, nil)

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

		packJSON := pack.PackJSON{}
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
		assert.Empty(t, packJSON.CanceledAt)
		assert.Empty(t, packJSON.Events)
	})

	t.Run("Shoud return error when skip from CREATED to DELIVERED", func(t *testing.T) {
		createdPack := createPack(t, nil)

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

func TestCancelPack(t *testing.T) {
	t.Run("Shoud cancel a pack successfully", func(t *testing.T) {
		createdPack := createPack(t, nil)

		resp, err := clientApp(httptest.NewRequest(
			http.MethodPost,
			"/packs/"+createdPack.ID+"/cancel",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		packJSON := pack.PackJSON{}
		err = json.NewDecoder(resp.Body).Decode(&packJSON)
		assert.Nil(t, err)

		assert.NotEmpty(t, packJSON.ID)
		assert.Equal(t, createdPack.Description, packJSON.Description)
		assert.Equal(t, createdPack.SenderName, packJSON.SenderName)
		assert.Equal(t, createdPack.ReceiverName, packJSON.ReceiverName)
		assert.Equal(t, "CANCELED", packJSON.Status.String())
		assert.NotEmpty(t, packJSON.CreatedAt)
		assert.NotEmpty(t, packJSON.UpdateAt)
		assert.Empty(t, packJSON.DeliveredAt)
		assert.NotEmpty(t, packJSON.CanceledAt)
		assert.Empty(t, packJSON.Events)
	})

	t.Run("Shoud return error when try to cancel a in transit pack", func(t *testing.T) {
		createdPack := createPack(t, nil)

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
			http.MethodPost,
			"/packs/"+createdPack.ID+"/cancel",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Shoud return error when try to cancel a delivered pack", func(t *testing.T) {
		createdPack := createPack(t, nil)

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

		resp, err = clientApp(httptest.NewRequest(
			http.MethodPost,
			"/packs/"+createdPack.ID+"/cancel",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Shoud return error when pack not found", func(t *testing.T) {
		resp, err := clientApp(httptest.NewRequest(
			http.MethodPost,
			"/packs/pack_not_found_1/cancle",
			nil,
		))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

type createPackParams struct {
	SenderName    string
	RecipientName string
}

func createPack(t *testing.T, params *createPackParams) pack.PackJSON {
	defer gock.Off()

	if params == nil {
		params = &createPackParams{}
	}

	if params.SenderName == "" {
		params.SenderName = "Loja ABC"
	}

	if params.RecipientName == "" {
		params.RecipientName = "João Silva"
	}

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

	resp, err := clientApp(httptest.NewRequest(
		http.MethodPost,
		"/packs",
		bytes.NewBuffer([]byte(`{
			"description": "Livros para entrega",
			"sender": "`+params.SenderName+`",
			"recipient": "`+params.RecipientName+`",
			"estimated_delivery_date": "2025-04-02"
		}`)),
	))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	packJSON := pack.PackJSON{}
	err = json.NewDecoder(resp.Body).Decode(&packJSON)
	assert.Nil(t, err)

	time.Sleep(1 * time.Millisecond) // wait for the gock to finish
	assert.True(t, gock.IsDone())

	return packJSON
}

func createEvent(t *testing.T, packID string) {
	resp, err := clientApp(httptest.NewRequest(
		http.MethodPost,
		"/pack_events",
		bytes.NewBuffer([]byte(`{
			"pack_id": "`+packID+`",
			"description": "Pacote chegou ao centro de distribuição",
			"location": "Centro de Distribuição São Paulo",
			"date": "2025-01-20T15:13:59Z"
		}`)),
	))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

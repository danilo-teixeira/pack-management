package pack_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pack-management/internal/services/pack"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePack(t *testing.T) {
	t.Run("Shoud create a pack successfully", func(t *testing.T) {
		body := []byte(`{
			"description": "Livros para entrega",
			"sender": "Loja ABC",
			"recipient": "João Silva",
			"estimated_delivery_date": "2025-04-02"
		}`)

		resp, err := clientApp(httptest.NewRequest(http.MethodPost, "/packs", bytes.NewBuffer(body)))
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
	})
	
	t.Run("Shoud return error when missing required fields", func(t *testing.T) {
		body := []byte(`{}`)

		resp, err := clientApp(httptest.NewRequest(http.MethodPost, "/packs", bytes.NewBuffer(body)))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

package render

import (
	"net/http"
	"testing"

	"github.com/janomonje/bed-n-breakfast/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var templateData models.TemplateData
	request, err := getSession()
	if err != nil {
		t.Fatal(err)
	}

	session.Put(request.Context(), "flash", "123")

	result := AddDefaultData(&templateData, request)

	if result.Flash != "123" {
		t.Error("flsh value of 123 not found in session")
	}
}

func getSession() (*http.Request, error) {
	request, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	ctx := request.Context()
	ctx, _ = session.Load(ctx, request.Header.Get("X-Session"))
	request = request.WithContext(ctx)
	return request, nil
}

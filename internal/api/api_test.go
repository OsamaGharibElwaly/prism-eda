package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"prismeda/internal/eda"
	"prismeda/internal/service"

	"github.com/gin-gonic/gin"
)

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(service.New(1<<20, "does-not-exist.csv"), "does-not-exist")
}

func upload(t *testing.T, router *gin.Engine, name, body string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", name)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte(body))
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/datasets", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func TestUploadProfilePreviewAndScatter(t *testing.T) {
	router := testRouter()
	res := upload(t, router, "sales.csv", "price,area\n10,100\n20,200\n30,300\n")
	if res.Code != http.StatusCreated {
		t.Fatalf("upload status = %d: %s", res.Code, res.Body.String())
	}
	var profile eda.Profile
	if err := json.Unmarshal(res.Body.Bytes(), &profile); err != nil {
		t.Fatal(err)
	}
	if profile.ID == "" || profile.Rows != 3 || profile.Columns[0].Name != "price" {
		t.Fatalf("bad profile: %+v", profile)
	}

	for _, path := range []string{
		"/api/datasets/" + profile.ID,
		"/api/datasets/" + profile.ID + "/preview?limit=2",
		"/api/datasets/" + profile.ID + "/scatter?x=price&y=area",
		"/api/health",
	} {
		get := httptest.NewRecorder()
		router.ServeHTTP(get, httptest.NewRequest(http.MethodGet, path, nil))
		if get.Code != http.StatusOK {
			t.Fatalf("GET %s = %d: %s", path, get.Code, get.Body.String())
		}
	}

	missing := httptest.NewRecorder()
	router.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/api/datasets/unknown", nil))
	if missing.Code != http.StatusNotFound {
		t.Fatalf("unknown dataset status = %d", missing.Code)
	}
}

func TestUploadValidation(t *testing.T) {
	router := testRouter()
	if res := upload(t, router, "notes.pdf", "a,b\n1,2\n"); res.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("bad extension status = %d", res.Code)
	}
	if res := upload(t, router, "empty.csv", ""); res.Code != http.StatusBadRequest {
		t.Fatalf("empty file status = %d", res.Code)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/datasets", strings.NewReader("not multipart"))
	req.Header.Set("Content-Type", "text/plain")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("non-multipart status = %d", res.Code)
	}
}

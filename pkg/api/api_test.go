package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	token = "123123123"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("API", func() {
	api := New(token, "")

	Describe("RequireToken middleware", func() {
		It("should return 401 Unauthorized when token is not provided", func() {
			handlerFunc := api.RequireToken(testHandler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/hello", nil)

			handlerFunc(rec, req)

			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 401 Unauthorized when token is invalid", func() {
			handlerFunc := api.RequireToken(testHandler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/hello", nil)
			req.Header.Set("Authorization", "Bearer 123")

			handlerFunc(rec, req)

			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should return 200 OK when token is valid", func() {
			handlerFunc := api.RequireToken(testHandler)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/hello", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			handlerFunc(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("Listen address handling", func() {
		It("should default an empty address", func() {
			Expect(NormalizeListenAddress("")).To(Equal(":8080"))
		})

		It("should prefix a bare port with a colon", func() {
			Expect(NormalizeListenAddress("9789")).To(Equal(":9789"))
		})

		It("should keep an address with a port untouched", func() {
			Expect(NormalizeListenAddress(":9789")).To(Equal(":9789"))
		})

		It("should keep a host:port pair untouched", func() {
			Expect(NormalizeListenAddress("127.0.0.1:9789")).To(Equal("127.0.0.1:9789"))
		})

		It("should normalize the address passed to New", func() {
			Expect(New(token, "9789").ListenAddress).To(Equal(":9789"))
		})
	})
})

func testHandler(w http.ResponseWriter, req *http.Request) {
	_, _ = io.WriteString(w, "Hello!")
}

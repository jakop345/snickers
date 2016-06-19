package snickers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/flavioribeiro/snickers/db/memory"
	"github.com/flavioribeiro/snickers/rest"
	"github.com/flavioribeiro/snickers/types"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rest API", func() {
	Context("/presets location", func() {
		var (
			response   *httptest.ResponseRecorder
			server     *mux.Router
			dbInstance *memory.Database
		)

		BeforeEach(func() {
			response = httptest.NewRecorder()
			server = rest.NewRouter()
			dbInstance, _ = memory.GetDatabase()
			dbInstance.ClearDatabase()
		})

		It("GET should return application/json on its content type", func() {
			request, _ := http.NewRequest("GET", "/presets", nil)
			server.ServeHTTP(response, request)
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		It("GET should return stored presets", func() {
			examplePreset := types.Preset{
				Name: "examplePreset",
			}
			dbInstance.StorePreset(examplePreset)
			expected := `[{"name":"examplePreset","video":{},"audio":{}}]`

			request, _ := http.NewRequest("GET", "/presets", nil)
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(string(response.Body.String())).To(Equal(expected))
		})

		It("POST should save a new preset", func() {
			preset := []byte(`{"name": "storedPreset", "video": {},"audio": {}}`)
			request, _ := http.NewRequest("POST", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			Expect(len(dbInstance.GetPresets())).To(Equal(1))
		})

		It("POST with malformed preset should return bad request", func() {
			preset := []byte(`{"neime: "badPreset}}`)
			request, _ := http.NewRequest("POST", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})
	})
})

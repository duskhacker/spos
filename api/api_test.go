package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/pborman/uuid"

	"github.com/duskhacker/cqrsnu/internal/github.com/gin-gonic/gin"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/ginkgo"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/gomega"
)

var pf = fmt.Printf

func performRequest(r http.Handler, method, path string, postData string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(postData))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

var _ = Describe("API", func() {
	var (
		router *gin.Engine
		tabID  string
	)

	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode)
		router = GinEngine()

		w := performRequest(router, "POST", "/opentab", `{"waitstaff":"Kinessa"}`)
		Expect(w.Code).To(Equal(http.StatusAccepted))

		response := struct {
			TabID string `json:"tabID" binding:"required"`
		}{}
		json.Unmarshal(w.Body.Bytes(), &response)
		tabID = response.TabID

		json := fmt.Sprintf(`{"tabID":"%s","items":[0,1]}`, tabID)
		w = performRequest(router, "POST", "/placeorder", json)
		Expect(w.Code).To(Equal(http.StatusAccepted))

	})

	It("returns the chef todo list ", func() {
		time.Sleep(time.Millisecond * 2)

		w := performRequest(router, "GET", "/cheftodolist", "")

		res := ChefTodoListResponse{}
		json.Unmarshal(w.Body.Bytes(), &res)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(uuid.Parse(res.Tabs[0].TabID)).ToNot(BeNil())
		Expect(res.Tabs[0].Items[0].MenuNumber).To(Equal(0))
		Expect(res.Tabs[0].Items[0].Description).To(Equal("Steak"))
		Expect(res.Tabs[0].Items[1].MenuNumber).To(Equal(1))
		Expect(res.Tabs[0].Items[1].Description).To(Equal("Burger"))

	})
})

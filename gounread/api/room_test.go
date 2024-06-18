package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"go.uber.org/mock/gomock"
	"gounread/api"
	"gounread/service/svcmock"
	"io"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Room", func() {
	var (
		svr      *api.Server
		m        *svcmock.MockServicer
		req      *http.Request
		recorder *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		var t gomock.TestReporter
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recorder = httptest.NewRecorder()
		m = svcmock.NewMockServicer(ctrl)
		svr = api.NewServer(m)
		svr.SetupRouter()
	})

	When("create room", func() {
		Context("", func() {
			type body struct {
				UserIDs []string `json:"user_ids"`
			}
			var userIDFixed = func() string {
				return "7e16486d-9f03-46d1-a429-531710f3d311"
			}
			var userID string
			BeforeEach(func() {
				userID = uuid.NewString()
			})

			DescribeTable("", func(auth func(*http.Request), bd *body, mock func([]string), respCheck func(*body)) {
				b, err := json.Marshal(bd)
				Expect(err).ShouldNot(HaveOccurred())
				req, err = http.NewRequest(http.MethodPost, "/rooms", bytes.NewReader(b))
				auth(req)
				mock(bd.UserIDs)
				svr.Router.ServeHTTP(recorder, req)
				respCheck(bd)
			},
				Entry("ok",
					func(r *http.Request) {
						r.Header.Set("user", userID)
					},
					&body{UserIDs: []string{uuid.NewString(), uuid.NewString()}},
					func(users []string) {
						m.EXPECT().CreateRoom(gomock.Any()).Times(1).Return(nil)
					},
					func(bd *body) {
						Expect(recorder.Code).To(Equal(http.StatusCreated))
					},
				),
				Entry("not enough user ids",
					func(r *http.Request) {
						r.Header.Set("user", userID)
					},
					&body{UserIDs: []string{}},
					func(users []string) {
						m.EXPECT().CreateRoom(gomock.Any()).Times(1).Return(fmt.Errorf("minimal user count 2. but: %d", len(users)+1))
					},
					func(bd *body) {
						Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
						var m map[string]string
						b, err := io.ReadAll(recorder.Body)
						Expect(err).ShouldNot(HaveOccurred())
						err = json.Unmarshal(b, &m)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(m).To(MatchAllKeys(Keys{
							"error": Equal(fmt.Sprintf("minimal user count 2. but: %d", len(bd.UserIDs)+1)),
						}))
					},
				),
				Entry("user ids nil",
					func(r *http.Request) {
						r.Header.Set("user", userID)
					},
					&body{UserIDs: nil},
					func(users []string) {
						m.EXPECT().CreateRoom(gomock.Any()).Times(0)
					},
					func(bd *body) {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
						var m map[string]string
						b, err := io.ReadAll(recorder.Body)
						Expect(err).ShouldNot(HaveOccurred())
						err = json.Unmarshal(b, &m)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(m).To(MatchAllKeys(Keys{
							"error": Equal("Key: 'UserIDs' Error:Field validation for 'UserIDs' failed on the 'required' tag"),
						}))
					},
				),
				Entry("including myself in invite users",
					func(r *http.Request) {
						r.Header.Set("user", userIDFixed())
					},
					&body{UserIDs: []string{userIDFixed()}}, // access field not working(empty string)
					func(users []string) {
						m.EXPECT().CreateRoom(gomock.Any()).Times(0)
					},
					func(bd *body) {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
						var m map[string]string
						b, err := io.ReadAll(recorder.Body)
						Expect(err).ShouldNot(HaveOccurred())
						err = json.Unmarshal(b, &m)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(m).To(MatchAllKeys(Keys{
							"error": Equal("must be exclude yourself in invite user list"),
						}))
					},
				),
			)
		})
	})
})

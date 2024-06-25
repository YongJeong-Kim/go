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
	"gounread/service"
	"gounread/service/svcmock"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
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
		ctrl := gomock.NewController(GinkgoT())
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

			DescribeTable("", func(auth func(*http.Request), bd *body, mock func([]string), checkResp func(*body)) {
				b, err := json.Marshal(bd)
				Expect(err).ShouldNot(HaveOccurred())
				req, err = http.NewRequest(http.MethodPost, "/rooms", bytes.NewReader(b))
				auth(req)
				mock(bd.UserIDs)
				svr.Router.ServeHTTP(recorder, req)
				checkResp(bd)
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

	When("list rooms by user id", func() {
		Context("", func() {
			url := func(userID string) string {
				return fmt.Sprintf("/users/%s/rooms", userID)
			}

			var (
				method string
			)

			BeforeEach(func() {
				method = http.MethodGet
			})

			It("ok", func() {
				userID := uuid.NewString()
				result := &service.ListRoomsByUserIDResult{
					RoomID:        uuid.NewString(),
					Time:          time.Now().UTC(),
					RecentMessage: "asdf",
					UnreadCount:   "2",
				}
				m.EXPECT().ListRoomsByUserID(userID).Times(1).Return([]*service.ListRoomsByUserIDResult{result}, nil)
				req, err := http.NewRequest(method, url(userID), nil)
				Expect(err).ShouldNot(HaveOccurred())
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusOK))

				b, err := io.ReadAll(recorder.Body)
				var resp []*service.ListRoomsByUserIDResult
				err = json.Unmarshal(b, &resp)
				Expect(err).ShouldNot(HaveOccurred())
				for _, r := range resp {
					Expect(r).To(PointTo(MatchAllFields(Fields{
						"RoomID":        Equal(result.RoomID),
						"Time":          Equal(result.Time),
						"RecentMessage": Equal(result.RecentMessage),
						"UnreadCount":   Equal(result.UnreadCount),
					})))
				}
			})

			It("bind invalid user id uri failed", func() {
				userID := "invalid user id"

				m.EXPECT().ListRoomsByUserID(userID).Times(0)
				req, err := http.NewRequest(method, url(userID), nil)
				Expect(err).ShouldNot(HaveOccurred())
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				b, err := io.ReadAll(recorder.Body)
				var resp map[string]string
				err = json.Unmarshal(b, &resp)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).To(MatchAllKeys(Keys{
					"error": Equal(fmt.Sprintf("invalid UUID length: %d", len(userID))),
				}))
			})

			It("user not found", func() {
				userID := uuid.NewString()

				m.EXPECT().ListRoomsByUserID(userID).Times(1).Return(nil, fmt.Errorf(fmt.Sprintf("user not found. %s", userID)))
				req, err := http.NewRequest(method, url(userID), nil)
				Expect(err).ShouldNot(HaveOccurred())
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				b, err := io.ReadAll(recorder.Body)
				var resp map[string]string
				err = json.Unmarshal(b, &resp)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).To(MatchAllKeys(Keys{
					"error": Equal(fmt.Sprintf("user not found. %s", userID)),
				}))
			})
		})
	})
})

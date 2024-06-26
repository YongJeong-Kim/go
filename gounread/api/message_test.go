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
	"gounread/embedded/notifymock"
	"gounread/repository"
	"io"

	//. "github.com/onsi/gomega/gstruct"
	"go.uber.org/mock/gomock"
	"gounread/api"
	"gounread/service"
	"gounread/service/svcmock"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Message", func() {
	var (
		m   *svcmock.MockServicer
		nm  *notifymock.MockNotifier
		svr *api.Server
	)
	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		ctrl := gomock.NewController(GinkgoT())
		defer ctrl.Finish()

		nm = notifymock.NewMockNotifier(ctrl)
		m = svcmock.NewMockServicer(ctrl)
		svr = api.NewServer(m, nm)
		svr.SetupRouter()
	})

	When("send message", func() {
		sendURI := func(roomID string) string {
			return fmt.Sprintf("/rooms/%s/send", roomID)
		}
		method := http.MethodPost
		var roomID, userID string

		BeforeEach(func() {
			roomID = uuid.NewString()
			userID = uuid.NewString()
		})

		It("ok", func() {
			msg := "asdf"
			b, err := json.Marshal(map[string]string{
				"message": msg,
			})
			Expect(err).ShouldNot(HaveOccurred())
			req, err := http.NewRequest(method, sendURI(roomID), bytes.NewReader(b))
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("user", userID)

			sent := time.Now().UTC()
			/*				param := &service.SendMessageParam{
							RoomID:  roomID,
							Sender:  userID,
							Message: msg,
							Sent:    sent,
						}*/
			result := &service.Payload{
				RoomID:  roomID,
				Message: msg,
				Sender:  userID,
				Sent:    sent,
				Unread:  []string{uuid.NewString()},
			}
			// exactly pass eq param only
			// if at least one field value is different, it will not passed.
			// m.EXPECT().SendMessage(param).Times(1).Return(result, nil)
			m.EXPECT().SendMessage(gomock.Any()).Times(1).Return(result, nil)
			nm.EXPECT().Publish("room."+roomID, gomock.Any()).Times(1).Return(nil)
			recorder := httptest.NewRecorder()
			svr.Router.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(http.StatusCreated))
		})

		Context("fail", func() {
			It("missing message", func() {
				req, err := http.NewRequest(method, sendURI(roomID), nil)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("user", userID)

				m.EXPECT().SendMessage(gomock.Any()).Times(0)
				nm.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(0)
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				b, err := io.ReadAll(recorder.Body)
				Expect(err).ShouldNot(HaveOccurred())
				var resp map[string]string
				err = json.Unmarshal(b, &resp)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).To(MatchAllKeys(Keys{
					"error": Equal("invalid request"),
				}))
			})

			It("send message", func() {
				msg := "asdf"
				b, err := json.Marshal(map[string]string{
					"message": msg,
				})
				Expect(err).ShouldNot(HaveOccurred())
				req, err := http.NewRequest(method, sendURI(roomID), bytes.NewReader(b))
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("user", userID)

				m.EXPECT().SendMessage(gomock.Any()).Times(1).Return(nil, fmt.Errorf("get users by room id error. invalid UUID \"invalid room id\""))
				nm.EXPECT().Publish(gomock.Any(), gomock.Any()).Times(0)
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				b, err = io.ReadAll(recorder.Body)
				var m map[string]string
				err = json.Unmarshal(b, &m)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(m).To(MatchAllKeys(Keys{
					"error": Equal("get users by room id error. invalid UUID \"invalid room id\""),
				}))
			})

			It("publish room", func() {
				msg := "asdf"
				b, err := json.Marshal(map[string]string{
					"message": msg,
				})
				Expect(err).ShouldNot(HaveOccurred())
				req, err := http.NewRequest(method, sendURI(roomID), bytes.NewReader(b))
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("user", userID)

				sent := time.Now().UTC()
				/*				param := &service.SendMessageParam{
								RoomID:  roomID,
								Sender:  userID,
								Message: msg,
								Sent:    sent,
							}*/
				result := &service.Payload{
					RoomID:  roomID,
					Message: msg,
					Sender:  userID,
					Sent:    sent,
					Unread:  []string{uuid.NewString()},
				}
				// exactly pass eq param only
				// if at least one field value is different, it will not passed.
				// m.EXPECT().SendMessage(param).Times(1).Return(result, nil)
				m.EXPECT().SendMessage(gomock.Any()).Times(1).Return(result, nil)
				nm.EXPECT().Publish("room."+roomID, gomock.Any()).Times(1).Return(fmt.Errorf("connection error"))
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				b, err = io.ReadAll(recorder.Body)
				Expect(err).ShouldNot(HaveOccurred())
				var m map[string]string
				err = json.Unmarshal(b, &m)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(m).To(MatchAllKeys(Keys{
					"error": Equal("room publish error. connection error"),
				}))
			})
		})
	})

	When("read message", func() {
		readURI := func(roomID string) string {
			return fmt.Sprintf("/rooms/%s/read", roomID)
		}
		method := http.MethodPut
		var roomID, userID string

		BeforeEach(func() {
			roomID = uuid.NewString()
			userID = uuid.NewString()
		})

		Context("ok", func() {
			It("new messages", func() {
				req, err := http.NewRequest(method, readURI(roomID), nil)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("user", userID)

				messages := []*repository.GetMessagesByRoomIDAndTimeResult{
					{
						RoomID: roomID,
						Sender: uuid.NewString(),
						Sent:   time.Now().UTC(),
						Msg:    "asd",
						Unread: []string{uuid.NewString()},
					},
				}
				m.EXPECT().ReadMessage(roomID, userID).Times(1).Return(messages, nil)
				m.EXPECT().GetRecentMessages(roomID, 10).Times(0)
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusOK))
				b, err := io.ReadAll(recorder.Body)
				var result []*repository.GetMessagesByRoomIDAndTimeResult
				err = json.Unmarshal(b, &result)
				Expect(err).ShouldNot(HaveOccurred())
				for i, r := range result {
					Expect(r).To(PointTo(MatchAllFields(Fields{
						"RoomID": Equal(messages[i].RoomID),
						"Sender": Equal(messages[i].Sender),
						"Sent":   Equal(messages[i].Sent),
						"Msg":    Equal(messages[i].Msg),
						"Unread": ContainElements(messages[i].Unread),
					})))
				}
			})

			It("no message", func() {
				req, err := http.NewRequest(method, readURI(roomID), nil)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("user", userID)

				messages := []*repository.GetRecentMessagesResult{
					{
						RoomID: roomID,
						Sender: uuid.NewString(),
						Sent:   time.Now().UTC(),
						Msg:    "asd",
					},
				}
				m.EXPECT().ReadMessage(roomID, userID).Times(1).Return(nil, nil)
				m.EXPECT().GetRecentMessages(roomID, 10).Times(1).Return(messages, nil)
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusOK))
				b, err := io.ReadAll(recorder.Body)
				var result []*repository.GetRecentMessagesResult
				err = json.Unmarshal(b, &result)
				Expect(err).ShouldNot(HaveOccurred())
				for i, r := range result {
					Expect(r).To(PointTo(MatchAllFields(Fields{
						"RoomID": Equal(messages[i].RoomID),
						"Sender": Equal(messages[i].Sender),
						"Sent":   Equal(messages[i].Sent),
						"Msg":    Equal(messages[i].Msg),
					})))
				}
			})
		})

		Context("fail", func() {
			It("get recent messages", func() {
				req, err := http.NewRequest(method, readURI(roomID), nil)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("user", userID)

				m.EXPECT().ReadMessage(roomID, userID).Times(1).Return(nil, nil)
				m.EXPECT().GetRecentMessages(roomID, 10).Times(1).Return(nil, fmt.Errorf("get recent messages next error. invalid UUID \"invalid room id\""))
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				b, err := io.ReadAll(recorder.Body)
				var m map[string]string
				err = json.Unmarshal(b, &m)
				Expect(m).To(MatchAllKeys(Keys{
					"error": Equal("get recent messages next error. invalid UUID \"invalid room id\""),
				}))
			})

			It("read message", func() {
				req, err := http.NewRequest(method, readURI(roomID), nil)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("user", userID)

				m.EXPECT().ReadMessage(roomID, userID).Times(1).Return(nil, fmt.Errorf("read message get message read time error. get message read time failed. invalid UUID \"invalid room id\""))
				m.EXPECT().GetRecentMessages(gomock.Any(), gomock.Any()).Times(0)
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				b, err := io.ReadAll(recorder.Body)
				var m map[string]string
				err = json.Unmarshal(b, &m)
				Expect(m).To(MatchAllKeys(Keys{
					"error": Equal("read message get message read time error. get message read time failed. invalid UUID \"invalid room id\""),
				}))
			})
		})
	})
})

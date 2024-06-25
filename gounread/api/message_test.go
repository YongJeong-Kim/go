package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		svr *api.Server
	)
	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		defer ctrl.Finish()

		m = svcmock.NewMockServicer(ctrl)
		svr = api.NewServer(m)
		svr.SetupRouter()
	})

	When("send message", func() {
		url := func(roomID string) string {
			return fmt.Sprintf("/rooms/%s/send", roomID)
		}
		method := http.MethodPost
		var roomID, userID string

		BeforeEach(func() {
			roomID = uuid.NewString()
			userID = uuid.NewString()
		})

		Context("", func() {
			/*	It("missing message", func() {
				req, err := http.NewRequest(method, url(roomID), nil)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("user", userID)

				m.EXPECT().SendMessage(gomock.Any()).Times(0)
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
			})*/

			It("ok", func() {
				msg := "asdf"
				b, err := json.Marshal(map[string]string{
					"message": msg,
				})
				Expect(err).ShouldNot(HaveOccurred())
				req, err := http.NewRequest(method, url(roomID), bytes.NewReader(b))
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
				m.EXPECT().SendMessage(&service.SendMessageParam{
					RoomID:  roomID,
					Sender:  userID,
					Message: msg,
					//Sent:    time.Now().UTC(),
				}).Times(1).Return(result, nil)
				recorder := httptest.NewRecorder()
				svr.Router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})
		})
	})
})

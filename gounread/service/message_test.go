package service_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gounread/api"
	"gounread/service"
	"log"
)

var _ = Describe("Message Service", func() {
	var svc service.Servicer
	BeforeEach(func(ctx SpecContext) {
		log.Println("asd")
		svc = service.NewService(api.NewSession())
	})

	When("send message", func() {
		Context("", func() {
			It("invalid uuid", func(ctx SpecContext) {
				err := svc.SendMessage(&service.SendMessageParam{
					RoomID:  "",
					Sender:  "",
					Message: "",
				})
				Expect(err).Should(MatchError("send message failed. invalid UUID \"\""))
			})

			It("invalid sender", func(ctx SpecContext) {
				err := svc.SendMessage(&service.SendMessageParam{
					RoomID:  uuid.NewString(),
					Sender:  "",
					Message: "",
				})
				Expect(err).Should(MatchError("send message failed. invalid UUID \"\""))
			})

			It("empty message", func(ctx SpecContext) {
				err := svc.SendMessage(&service.SendMessageParam{
					RoomID:  uuid.NewString(),
					Sender:  uuid.NewString(),
					Message: "",
				})
				Expect(err).To(Succeed())
			})

			It("ok", func(ctx SpecContext) {
				err := svc.SendMessage(&service.SendMessageParam{
					RoomID:  uuid.NewString(),
					Sender:  uuid.NewString(),
					Message: "asdiw",
				})
				Expect(err).To(Succeed())
			})
		})
	})

	When("read message", func() {
		Context("", func() {
			It("", func() {
				err := svc.ReadMessage(uuid.NewString(), uuid.NewString())
				Expect(err).To(Succeed())
			})
		})
	})

	When("get all rooms by user", func() {
		Context("", func() {
			It("", func() {
				//rooms := svc.GetRoomsByUserID("febba554-152e-496a-add5-31d0672fdc2a")
				times := svc.GetAllRoomsReadMessageTime("febba554-152e-496a-add5-31d0672fdc2a")
				counts := svc.GetRoomsUnreadMessageCount(times)
				log.Println(counts)
			})
		})
	})
})

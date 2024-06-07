package service_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gounread/api"
	"gounread/service"
	"log"
	"time"
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
				Contain
			})
		})
		Context("dd", func() {
			It("dd", func(ctx SpecContext) {
				Expect(true, true)
			}, SpecTimeout(time.Second*3))
		})
	})
})

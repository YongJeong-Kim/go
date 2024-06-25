package service_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"gounread/api"
	"gounread/repository"
	"gounread/service"
	"time"
)

var _ = Describe("Message Service", func() {
	var (
		svc         service.Servicer
		sender      string
		inviteUsers []string
		allUsers    []string
	)

	BeforeEach(func(ctx SpecContext) {
		var repo repository.Repositorier = repository.NewRepository(api.NewSession())
		svc = service.NewService(repo)
		sender = uuid.NewString()
		inviteUsers = []string{uuid.NewString(), uuid.NewString()}
		allUsers = append(allUsers, sender)
		allUsers = append(allUsers, inviteUsers...)
	})
	AfterEach(func() {
		allUsers = nil
	})

	When("send message", func() {
		Context("ok", func() {
			It("ok", func(ctx SpecContext) {
				createRoomForTest(svc, allUsers)
				r1, err := svc.GetRoomsByUserID(sender)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r1).To(HaveLen(1))

				param := &service.SendMessageParam{
					RoomID:  r1[0].RoomID,
					Sender:  sender,
					Message: "asdiw",
					Sent:    time.Now().UTC(),
				}
				payload, err := svc.SendMessage(param)
				Expect(payload).To(PointTo(MatchAllFields(Fields{
					"RoomID":  Equal(param.RoomID),
					"Sender":  Equal(param.Sender),
					"Message": Equal(param.Message),
					"Sent":    Equal(param.Sent),
					"Unread":  ContainElements(inviteUsers),
				})))
			})

			It("empty message(ok)", func(ctx SpecContext) {
				createRoomForTest(svc, allUsers)
				u1 := allUsers[0]
				r1, err := svc.GetRoomsByUserID(u1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r1).To(HaveLen(1))

				param := &service.SendMessageParam{
					RoomID:  r1[0].RoomID,
					Sender:  sender,
					Message: "",
					Sent:    time.Now().UTC(),
				}
				payload, err := svc.SendMessage(param)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(payload).To(PointTo(MatchAllFields(Fields{
					"RoomID":  Equal(param.RoomID),
					"Sender":  Equal(param.Sender),
					"Message": Equal(param.Message),
					"Sent":    Equal(param.Sent),
					"Unread":  ContainElements(inviteUsers),
				})))
			})
		})

		Context("fail", func() {
			It("invalid room id", func(ctx SpecContext) {
				createRoomForTest(svc, allUsers)
				r1, err := svc.GetRoomsByUserID(sender)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r1).To(HaveLen(1))

				param := &service.SendMessageParam{
					RoomID:  "invalid room id",
					Sender:  sender,
					Message: "asdiw",
					Sent:    time.Now().UTC(),
				}
				payload, err := svc.SendMessage(param)
				Expect(err).To(MatchError("get users by room id error. invalid UUID \"invalid room id\""))
				Expect(payload).To(BeNil())
			})

			It("invalid sender", func(ctx SpecContext) {
				createRoomForTest(svc, allUsers)
				r1, err := svc.GetRoomsByUserID(sender)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r1).To(HaveLen(1))

				param := &service.SendMessageParam{
					RoomID:  r1[0].RoomID,
					Sender:  "",
					Message: "asdiw",
					Sent:    time.Now().UTC(),
				}
				payload, err := svc.SendMessage(param)
				Expect(err).Should(MatchError("send message failed. invalid UUID \"\""))
				Expect(payload).To(BeNil())
			})
		})
	})

	When("read message", func() {
		It("ok", func() {
			createRoomForTest(svc, allUsers)
			_, u2, u3 := allUsers[0], allUsers[1], allUsers[2]
			r2, err := svc.GetRoomsByUserID(u2)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r2).To(HaveLen(1))

			msg := "asdadwd"
			sent := time.Now().UTC()
			sendMessageForTest(svc, &repository.CreateMessageParam{
				RoomID:      r2[0].RoomID,
				Sender:      sender,
				Message:     msg,
				Sent:        sent,
				UnreadUsers: inviteUsers,
			})

			messages, err := svc.ReadMessage(r2[0].RoomID, u2)
			Expect(err).ShouldNot(HaveOccurred())
			for _, m := range messages {
				Expect(m).To(PointTo(MatchAllFields(Fields{
					"RoomID": Equal(r2[0].RoomID),
					"Sender": Equal(sender),
					"Msg":    Equal(msg),
					"Sent":   BeTemporally("~", sent, time.Millisecond),
					"Unread": ContainElements([]string{u3}),
				})))
			}

			messages, err = svc.ReadMessage(r2[0].RoomID, u3)
			Expect(err).ShouldNot(HaveOccurred())
			for _, m := range messages {
				Expect(m).To(PointTo(MatchAllFields(Fields{
					"RoomID": Equal(r2[0].RoomID),
					"Sender": Equal(sender),
					"Msg":    Equal(msg),
					"Sent":   BeTemporally("~", sent, time.Millisecond),
					"Unread": BeEmpty(),
				})))
			}
		})

		Context("fail", func() {
			It("invalid room id", func() {
				createRoomForTest(svc, allUsers)
				_, u2, _ := allUsers[0], allUsers[1], allUsers[2]
				r1, err := svc.GetRoomsByUserID(u2)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r1).To(HaveLen(1))

				sendMessageForTest(svc, &repository.CreateMessageParam{
					RoomID:      r1[0].RoomID,
					Sender:      sender,
					Message:     "egfhow",
					Sent:        time.Now().UTC(),
					UnreadUsers: inviteUsers,
				})

				messages, err := svc.ReadMessage("invalid room id", u2)
				Expect(err).To(MatchError("read message get message read time error. get message read time failed. invalid UUID \"invalid room id\""))
				Expect(messages).To(BeNil())
				Expect(messages).To(BeEmpty())
			})

			It("invalid user id", func() {
				createRoomForTest(svc, allUsers)
				_, u2, _ := allUsers[0], allUsers[1], allUsers[2]
				r1, err := svc.GetRoomsByUserID(u2)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r1).To(HaveLen(1))

				sendMessageForTest(svc, &repository.CreateMessageParam{
					RoomID:      r1[0].RoomID,
					Sender:      sender,
					Message:     "egfhow",
					Sent:        time.Now().UTC(),
					UnreadUsers: inviteUsers,
				})

				messages, err := svc.ReadMessage(r1[0].RoomID, "invalid user id")
				Expect(err).To(MatchError("read message get message read time error. get message read time failed. invalid UUID \"invalid user id\""))
				Expect(messages).To(BeNil())
				Expect(messages).To(BeEmpty())
			})
		})
	})

	When("get recent messages", func() {
		It("ok", func() {
			createRoomForTest(svc, allUsers)
			r1, err := svc.GetRoomsByUserID(sender)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r1).To(HaveLen(1))
			sent := time.Now().UTC()
			msg := "sdofihowe"
			sendMessageForTest(svc, &repository.CreateMessageParam{
				RoomID:      r1[0].RoomID,
				Sender:      sender,
				Message:     msg,
				Sent:        sent,
				UnreadUsers: inviteUsers,
			})

			messages, err := svc.GetRecentMessages(r1[0].RoomID, 10)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(messages).To(HaveLen(1))
			for _, m := range messages {
				Expect(m).To(PointTo(MatchAllFields(Fields{
					"RoomID": Equal(r1[0].RoomID),
					"Sender": Equal(sender),
					"Msg":    Equal(msg),
					"Sent":   BeTemporally("~", sent, time.Millisecond),
				})))
			}
		})

		Context("fail", func() {
			It("invalid room id", func() {
				createRoomForTest(svc, allUsers)
				r1, err := svc.GetRoomsByUserID(sender)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r1).To(HaveLen(1))

				sendMessageForTest(svc, &repository.CreateMessageParam{
					RoomID:      r1[0].RoomID,
					Sender:      sender,
					Message:     "sdlfihios",
					Sent:        time.Now().UTC(),
					UnreadUsers: inviteUsers,
				})

				messages, err := svc.GetRecentMessages("invalid room id", 10)
				Expect(err).To(MatchError("get recent messages next error. invalid UUID \"invalid room id\""))
				Expect(messages).Should(BeNil())
				Expect(messages).Should(BeEmpty())
			})

			It("room not found", func() {
				createRoomForTest(svc, allUsers)
				r1, err := svc.GetRoomsByUserID(sender)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r1).To(HaveLen(1))

				sendMessageForTest(svc, &repository.CreateMessageParam{
					RoomID:      r1[0].RoomID,
					Sender:      sender,
					Message:     "sdlfihios",
					Sent:        time.Now().UTC(),
					UnreadUsers: inviteUsers,
				})

				roomNotFound := uuid.NewString()
				messages, err := svc.GetRecentMessages(roomNotFound, 10)
				Expect(err).To(MatchError(fmt.Sprintf("room not found. %s", roomNotFound)))
				Expect(messages).Should(BeNil())
				Expect(messages).Should(BeEmpty())
			})
		})
	})
})

func sendMessageForTest(svc service.Servicer, param *repository.CreateMessageParam) {
	payload, err := svc.SendMessage(&service.SendMessageParam{
		RoomID:  param.RoomID,
		Sender:  param.Sender,
		Message: param.Message,
		Sent:    param.Sent,
	})
	Expect(err).ShouldNot(HaveOccurred())
	Expect(payload).To(PointTo(MatchAllFields(Fields{
		"RoomID":  Equal(param.RoomID),
		"Sender":  Equal(param.Sender),
		"Message": Equal(param.Message),
		"Sent":    Equal(param.Sent),
		"Unread":  ContainElements(param.UnreadUsers),
	})))
}

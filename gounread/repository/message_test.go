package repository_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"gounread/api"
	"gounread/repository"
	"sort"
	"time"
)

var _ = Describe("Message", func() {
	var (
		repo        repository.Repositorier
		roomID      string
		sender      string
		inviteUsers []string
	)
	BeforeEach(func() {
		session := api.NewSession()
		repo = repository.NewRepository(session)
		//defer session.Close()
		roomID = uuid.NewString()
		sender = uuid.NewString()
		inviteUsers = []string{uuid.NewString(), uuid.NewString()}
	})
	AfterEach(func() {
		inviteUsers = nil
	})

	When("create message", func() {
		Context("", func() {
			It("ok", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)

				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "dd",
					UnreadUsers: inviteUsers,
				}, 1)

				count, err := repo.GetMessageCountByRoomIDAndSent(roomID, sent)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("invalid room id", func() {
				err := repo.CreateMessage(&repository.CreateMessageParam{
					RoomID:      "invalid room id",
					Sender:      sender,
					Sent:        time.Now().UTC(),
					Message:     "dd",
					UnreadUsers: []string{},
				})
				Expect(err).To(MatchError(fmt.Errorf("send message failed. invalid UUID \"invalid room id\"")))
			})
		})
	})

	When("update recent message", func() {
		Context("", func() {
			It("ok", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				newMessage := "new message"
				err := repo.UpdateRecentMessage(roomID, newMessage)
				Expect(err).ShouldNot(HaveOccurred())

				recent, err := repo.GetRecentMessageByRoomID(roomID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(recent).To(PointTo(MatchAllFields(Fields{
					"RoomID":        Equal(roomID),
					"RecentMessage": Equal(newMessage),
				})))
			})

			It("room not found", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				err := repo.UpdateRecentMessage("room not found", "new message")
				Expect(err).To(MatchError(fmt.Errorf("update recent message failed. invalid UUID \"room not found\"")))
			})
		})
	})

	When("get unread message count", func() {
		Context("", func() {
			It("ok", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        time.Now().UTC(),
					Message:     "asdfasdf",
					UnreadUsers: inviteUsers,
				}, 1)

				t, err := repo.GetMessageReadTime(roomID, inviteUsers[0])
				Expect(err).ShouldNot(HaveOccurred())

				count, err := repo.GetUnreadMessageCount(roomID, t)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("unread message not found", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        time.Now().UTC(),
					Message:     "asdfasdf",
					UnreadUsers: inviteUsers,
				}, 1)

				err := repo.UpdateMessageReadTime(roomID, inviteUsers[0], time.Now().UTC())
				Expect(err).ShouldNot(HaveOccurred())
				t, err := repo.GetMessageReadTime(roomID, inviteUsers[0])
				Expect(err).ShouldNot(HaveOccurred())

				count, err := repo.GetUnreadMessageCount(roomID, t)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(count).To(BeZero())
			})

			It("invalid room id", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        time.Now().UTC(),
					Message:     "asdfasdf",
					UnreadUsers: inviteUsers,
				}, 1)

				t, err := repo.GetMessageReadTime(roomID, inviteUsers[0])
				Expect(err).ShouldNot(HaveOccurred())

				count, err := repo.GetUnreadMessageCount("invalid room id", t)
				Expect(err).To(MatchError(fmt.Errorf("get message status unread count failed. invalid UUID \"invalid room id\"")))
				Expect(count).To(BeZero())
			})
		})
	})

	When("get messages by room id and time", func() {
		Context("", func() {
			It("ok", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sendCount := 5
				start := time.Now().UTC()
				var times []time.Time
				for range sendCount {
					n := time.Now().UTC()
					createMessageForTest(repo, &repository.CreateMessageParam{
						RoomID:      roomID,
						Sender:      sender,
						Sent:        n,
						Message:     "asdf",
						UnreadUsers: inviteUsers,
					}, 1)
					times = append(times, n)
				}
				sort.Slice(times, func(i, j int) bool {
					return times[i].After(times[j])
				})
				end := time.Now().UTC()
				messages := repo.GetMessagesByRoomIDAndTime(roomID, start, end)
				Expect(messages).To(HaveLen(sendCount))

				for i, m := range messages {
					Expect(m).To(PointTo(MatchAllFields(Fields{
						"RoomID": Equal(roomID),
						"Sent":   BeTemporally("~", times[i], 500*time.Millisecond),
						"Sender": Equal(sender),
						"Msg":    Equal("asdf"),
						"Unread": ContainElements(inviteUsers),
					})))
				}
			})

			It("room not found", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				start := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      uuid.NewString(),
					Sender:      sender,
					Sent:        time.Now().UTC(),
					Message:     "asdf",
					UnreadUsers: inviteUsers,
				}, 1)
				end := time.Now().UTC()
				messages := repo.GetMessagesByRoomIDAndTime(roomID, start, end)
				Expect(messages).To(BeZero())
			})

			It("invalid room id", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				start := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        time.Now().UTC(),
					Message:     "asdf",
					UnreadUsers: inviteUsers,
				}, 1)
				end := time.Now().UTC()
				messages := repo.GetMessagesByRoomIDAndTime("invalid room id", start, end)
				Expect(messages).To(BeNil())
				Expect(messages).To(BeZero())
			})
		})
	})

	When("get message read time", func() {
		Context("", func() {
			It("yourself", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aaa",
					UnreadUsers: inviteUsers,
				}, 1)
				err := repo.UpdateMessageReadTime(roomID, sender, sent)
				Expect(err).ShouldNot(HaveOccurred())

				t, err := repo.GetMessageReadTime(roomID, sender)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(sent).To(BeTemporally("~", t, time.Millisecond))
			})

			It("read user", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aaa",
					UnreadUsers: inviteUsers,
				}, 1)
				u1 := inviteUsers[0]
				err := repo.UpdateMessageReadTime(roomID, u1, sent)
				Expect(err).ShouldNot(HaveOccurred())

				t, err := repo.GetMessageReadTime(roomID, u1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(sent).To(BeTemporally("~", t, time.Millisecond))
			})

			It("invalid room id", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aaa",
					UnreadUsers: inviteUsers,
				}, 1)
				err := repo.UpdateMessageReadTime(roomID, sender, sent)
				Expect(err).ShouldNot(HaveOccurred())

				t, err := repo.GetMessageReadTime("invalid room id", sender)
				Expect(err).To(MatchError(fmt.Errorf("get message read time failed. invalid UUID \"invalid room id\"")))
				Expect(t).To(Equal(time.Time{}))
			})

			It("room not found", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aaa",
					UnreadUsers: inviteUsers,
				}, 1)
				err := repo.UpdateMessageReadTime(roomID, sender, sent)
				Expect(err).ShouldNot(HaveOccurred())

				t, err := repo.GetMessageReadTime(uuid.NewString(), sender)
				Expect(err).To(MatchError(fmt.Errorf("get message read time failed. not found")))
				Expect(t).To(Equal(time.Time{}))
			})

			It("invalid user id", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aaa",
					UnreadUsers: inviteUsers,
				}, 1)
				err := repo.UpdateMessageReadTime(roomID, sender, sent)
				Expect(err).ShouldNot(HaveOccurred())

				t, err := repo.GetMessageReadTime(roomID, "invalid user id")
				Expect(err).To(MatchError(fmt.Errorf("get message read time failed. invalid UUID \"invalid user id\"")))
				Expect(t).To(Equal(time.Time{}))
			})

			It("user not found", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aaa",
					UnreadUsers: inviteUsers,
				}, 1)
				err := repo.UpdateMessageReadTime(roomID, sender, sent)
				Expect(err).ShouldNot(HaveOccurred())

				t, err := repo.GetMessageReadTime(roomID, uuid.NewString())
				Expect(err).To(MatchError(fmt.Errorf("get message read time failed. not found")))
				Expect(t).To(Equal(time.Time{}))
			})
		})
	})

	When("get recent message by room id", func() {
		Context("", func() {
			It("ok", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				msg := "asdfasdf"
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     msg,
					UnreadUsers: inviteUsers,
				}, 1)

				recent, err := repo.GetRecentMessageByRoomID(roomID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(recent).To(PointTo(MatchAllFields(Fields{
					"RoomID":        Equal(roomID),
					"RecentMessage": Equal(msg),
				})))
			})

			It("invalid room id", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				msg := "asdfasdf"
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     msg,
					UnreadUsers: inviteUsers,
				}, 1)

				recent, err := repo.GetRecentMessageByRoomID("invalid room id")
				Expect(err).To(MatchError(fmt.Errorf("get room recent message failed. invalid UUID \"invalid room id\"")))
				Expect(recent).To(BeNil())
			})

			It("room not found", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				msg := "asdfasdf"
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     msg,
					UnreadUsers: inviteUsers,
				}, 1)

				recent, err := repo.GetRecentMessageByRoomID(uuid.NewString())
				Expect(err).To(MatchError(fmt.Errorf("get room recent message failed. not found")))
				Expect(recent).To(BeNil())
			})
		})
	})

	When("get all rooms read message time", func() {
		Context("", func() {
			It("ok", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aqaq",
					UnreadUsers: inviteUsers,
				}, 1)
				updateReadTime := time.Now().UTC()
				err := repo.UpdateMessageReadTime(roomID, sender, updateReadTime)
				Expect(err).ShouldNot(HaveOccurred())

				readTimes := repo.GetAllRoomsReadMessageTime(sender)
				for _, r := range readTimes {
					Expect(r).To(PointTo(MatchAllFields(Fields{
						"RoomID":   Equal(roomID),
						"ReadTime": BeTemporally("~", r.ReadTime, 500*time.Millisecond),
					})))
				}
			})

			It("user not found", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aqaq",
					UnreadUsers: inviteUsers,
				}, 1)
				updateReadTime := time.Now().UTC()
				err := repo.UpdateMessageReadTime(roomID, sender, updateReadTime)
				Expect(err).ShouldNot(HaveOccurred())

				readTimes := repo.GetAllRoomsReadMessageTime(uuid.NewString())
				Expect(readTimes).To(BeNil())
				Expect(readTimes).To(HaveLen(0))
			})

			It("invalid user id", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aqaq",
					UnreadUsers: inviteUsers,
				}, 1)
				updateReadTime := time.Now().UTC()
				err := repo.UpdateMessageReadTime(roomID, sender, updateReadTime)
				Expect(err).ShouldNot(HaveOccurred())

				readTimes := repo.GetAllRoomsReadMessageTime("invalid user id")
				Expect(readTimes).To(BeNil())
				Expect(readTimes).To(HaveLen(0))
			})
		})
	})

	/*When("update unread message batch", func() {
		Context("ok", func() {
			It("ok", func() {
				createRoomForTest(repo, roomID, sender, inviteUsers)
				sent := time.Now().UTC()
				createMessageForTest(repo, &repository.CreateMessageParam{
					RoomID:      roomID,
					Sender:      sender,
					Sent:        sent,
					Message:     "aqaq",
					UnreadUsers: inviteUsers,
				}, 5)
				updateReadTime := time.Now().UTC()
				u1 := inviteUsers[0]
				err := repo.UpdateMessageReadTime(roomID, u1, updateReadTime)
				Expect(err).ShouldNot(HaveOccurred())

				err = repo.UpdateUnreadMessageBatch(&repository.UpdateUnreadMessageBatchParam{
					UserID:   u1,
					Messages: repo.GetMessagesByRoomIDAndTime(roomID, updateReadTime, time.Now().UTC()),
				})
				Expect(err).ShouldNot(HaveOccurred())

				err := repo.UpdateMessageReadTime(roomID, u1, updateReadTime)
				Expect(err).ShouldNot(HaveOccurred())

				repo.GetMessagesByRoomIDAndTime(roomID, start, end)
			})
		})

		Context("", func() {

		})
	})*/
})

func createMessageForTest(repo repository.Repositorier, param *repository.CreateMessageParam, repeat int) {
	if repeat < 1 {
		panic("create message for test minimum repeat is 1")
	}

	for range repeat {
		err := repo.CreateMessage(param)
		Expect(err).ShouldNot(HaveOccurred())
		err = repo.UpdateRecentMessage(param.RoomID, param.Message)
		Expect(err).ShouldNot(HaveOccurred())
	}
}

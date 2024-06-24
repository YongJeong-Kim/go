package repository_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gounread/api"
	"gounread/repository"
	"time"
)

var _ = Describe("Room", func() {
	var repo repository.Repositorier
	var users []string
	BeforeEach(func() {
		repo = repository.NewRepository(api.NewSession())
		for range 5 {
			users = append(users, uuid.NewString())
		}
		Expect(users).ShouldNot(BeNil())
	})
	AfterEach(func() {
		users = nil
	})

	When("create room", func() {
		Context("create room with 5 users", func() {
			It("ok", func() {
				createRoomForTest(repo, uuid.NewString(), uuid.NewString(), users)
			})

			It("invalid user id", func() {
				users[0] = "invalid user id"
				err := repo.CreateRoom(uuid.NewString(), users)
				Expect(err).To(BeEquivalentTo(fmt.Errorf("create room error. invalid UUID \"invalid user id\"")))
			})
		})
	})

	When("get rooms by user id", func() {
		Context("get rooms by 5 user ids", func() {
			BeforeEach(func() {
				err := repo.CreateRoom(uuid.NewString(), users)
				Expect(err).Should(BeNil())
			})

			It("ok", func() {
				for _, u := range users {
					rooms, err := repo.GetRoomsByUserID(u)
					Expect(err).Should(BeNil())
					for _, r := range rooms {
						Expect(r.RoomID).ShouldNot(BeNil())
						Expect(r.Time).ShouldNot(BeNil())
						Expect(r.RecentMessage).Should(BeEmpty())
					}
				}
			})

			It("invalid user id", func() {
				for _, u := range users {
					u = "invalid user id"
					rooms, err := repo.GetRoomsByUserID(u)
					Expect(rooms).To(BeNil())
					Expect(err).To(MatchError("GetRoomsByUserID next failed. invalid UUID \"invalid user id\""))
				}
			})
		})
	})

	When("get users by room id", func() {
		BeforeEach(func() {
			err := repo.CreateRoom(uuid.NewString(), users)
			Expect(err).To(Succeed())
		})

		Context("get user in room <-> get room by user", func() {
			It("ok", func() {
				for _, u := range users {
					rooms, err := repo.GetRoomsByUserID(u)
					Expect(err).Should(BeNil())
					for _, r := range rooms {
						Expect(r.RoomID).ShouldNot(BeNil())
						Expect(r.Time).ShouldNot(BeNil())
						Expect(r.RecentMessage).Should(BeEmpty())

						usersInRoom, err := repo.GetUsersByRoomID(r.RoomID)
						Expect(err).Should(BeNil())
						Expect(usersInRoom).To(HaveLen(5))
						Expect(usersInRoom).Should(ContainElement(u))
					}
				}
			})

			It("invalid room id", func() {
				for _, u := range users {
					rooms, err := repo.GetRoomsByUserID(u)
					Expect(err).Should(BeNil())
					for _, r := range rooms {
						Expect(r.RoomID).ShouldNot(BeNil())
						Expect(r.Time).ShouldNot(BeNil())
						Expect(r.RecentMessage).Should(BeEmpty())

						usersInRoom, err := repo.GetUsersByRoomID("invalid room id")
						Expect(usersInRoom).Should(BeNil())
						Expect(err).To(Equal(fmt.Errorf("get users by room id error. invalid UUID \"invalid room id\"")))
					}
				}
			})
		})
	})
})

func createRoomForTest(repo repository.Repositorier, roomID, owner string, inviteUsers []string) {
	allUsers := append(inviteUsers, owner)
	err := repo.CreateRoom(roomID, allUsers)
	Expect(err).ShouldNot(HaveOccurred())

	usersInRoom, err := repo.GetUsersByRoomID(roomID)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(usersInRoom).To(HaveLen(len(allUsers)))

	now := time.Now().UTC()
	for _, u := range allUsers {
		err = repo.UpdateMessageReadTime(roomID, u, now)
		Expect(err).ShouldNot(HaveOccurred())
	}
}

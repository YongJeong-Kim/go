package service_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	"gounread/api"
	"gounread/repository"
	"gounread/service"
	"strconv"
)

var _ = Describe("Room", func() {
	var svc service.Servicer
	var users []string

	BeforeEach(func() {
		var repo repository.Repositorier = repository.NewRepository(api.NewSession())
		svc = service.NewService(repo)

		for range 3 {
			users = append(users, uuid.NewString())
		}
	})
	AfterEach(func() {
		users = nil
	})

	Context("room count", func() {
		It("user1, user2, user3 room count", func() {
			createRoomForTest(svc, users)
			createRoomForTest(svc, users)

			u1, u2, u3 := users[0], users[1], users[2]
			rooms1, err := svc.GetRoomsByUserID(u1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rooms1).To(HaveLen(2))

			rooms2, err := svc.GetRoomsByUserID(u2)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rooms2).To(HaveLen(2))

			rooms3, err := svc.GetRoomsByUserID(u3)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rooms3).To(HaveLen(2))
		})

		It("user not found. empty room", func() {
			uID := "sdfseoh"
			rooms, err := svc.GetRoomsByUserID(uID)
			Expect(err).To(MatchError(fmt.Sprintf("GetRoomsByUserID next failed. invalid UUID \"%s\"", uID)))
			Expect(rooms).To(BeNil())
			Expect(rooms).To(BeEmpty())
		})
	})

	DescribeTable(
		"create room",
		func(users []string, matcher types.GomegaMatcher) {
			err := svc.CreateRoom(users)
			Expect(err).Should(matcher)
		},
		Entry("ok",
			[]string{uuid.NewString(), uuid.NewString()},
			BeNil(),
		),
		Entry("invalid user id",
			[]string{"no uuid", ""},
			BeEquivalentTo(fmt.Errorf("invalid user id: invalid UUID length: %d. %s", len("no uuid"), "no uuid")),
		),
		Entry("empty user id",
			[]string{"", "no uuid"},
			BeEquivalentTo(fmt.Errorf("invalid user id: invalid UUID length: %d. %s", len(""), "")),
		),
		Entry("equal greater than 2 user count",
			[]string{uuid.NewString()},
			BeEquivalentTo(fmt.Errorf("minimal user count 2. but: %d", 1)),
		),
	)

	When("list rooms by user id", func() {
		It("ok", func() {
			createRoomForTest(svc, users)
			u1, u2, u3 := users[0], users[1], users[2]
			r1, err := svc.GetRoomsByUserID(u1)
			Expect(err).ShouldNot(HaveOccurred())
			r2, err := svc.GetRoomsByUserID(u2)
			Expect(err).ShouldNot(HaveOccurred())
			r3, err := svc.GetRoomsByUserID(u3)
			Expect(err).ShouldNot(HaveOccurred())

			rooms1, err := svc.ListRoomsByUserID(u1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rooms1).To(HaveLen(1))
			for _, r := range rooms1 {
				Expect(r).To(PointTo(MatchAllFields(Fields{
					"RoomID":        Equal(r1[0].RoomID),
					"Time":          Equal(r1[0].Time),
					"RecentMessage": Equal(r1[0].RecentMessage),
					"UnreadCount":   Equal(strconv.Itoa(0)),
				})))
			}

			rooms2, err := svc.ListRoomsByUserID(u2)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rooms2).To(HaveLen(1))
			for _, r := range rooms2 {
				Expect(r).To(PointTo(MatchAllFields(Fields{
					"RoomID":        Equal(r2[0].RoomID),
					"Time":          Equal(r2[0].Time),
					"RecentMessage": Equal(r2[0].RecentMessage),
					"UnreadCount":   Equal(strconv.Itoa(0)),
				})))
			}

			rooms3, err := svc.ListRoomsByUserID(u3)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rooms3).To(HaveLen(1))
			for _, r := range rooms3 {
				Expect(r).To(PointTo(MatchAllFields(Fields{
					"RoomID":        Equal(r3[0].RoomID),
					"Time":          Equal(r3[0].Time),
					"RecentMessage": Equal(r3[0].RecentMessage),
					"UnreadCount":   Equal(strconv.Itoa(0)),
				})))
			}
		})

		Context("fail", func() {
			It("user not found", func() {
				createRoomForTest(svc, users)

				userNotFound := uuid.NewString()
				rooms, err := svc.ListRoomsByUserID(userNotFound)
				Expect(err).To(MatchError(fmt.Sprintf("user not found. %s", userNotFound)))
				Expect(rooms).To(BeNil())
				Expect(rooms).To(BeEmpty())
			})

			It("invalid user id", func() {
				createRoomForTest(svc, users)

				rooms, err := svc.ListRoomsByUserID("invalid user id")
				Expect(err).To(MatchError("GetRoomsByUserID next failed. invalid UUID \"invalid user id\""))
				Expect(rooms).To(BeNil())
				Expect(rooms).To(BeEmpty())
			})
		})
	})

	/*Context("concurrency set data type", func() {
		It("100", func() {
			count := 3000
			var users []string
			for _ = range count {
				users = append(users, uuid.NewString())
			}
			err := svc.CreateRoom(users)
			Expect(err).To(Succeed())

			room, err := svc.GetRoomsByUserID(users[0])

			sess := api.NewSession()
			q := `UPDATE room SET users = users - ? WHERE room_id = ?`
			var wg sync.WaitGroup
			for i := range count {
				wg.Add(1)
				go func(ii int) {
					defer wg.Done()
					err := sess.Query(q, nil).Bind([]string{users[ii]}, room[0].RoomID).Exec()
					if err != nil {
						log.Println(err)
					}
					Expect(err).To(Succeed())
				}(i)
			}
			wg.Wait()
		})
	})*/
})

func createRoomForTest(svc service.Servicer, users []string) {
	err := svc.CreateRoom(users)
	Expect(err).ShouldNot(HaveOccurred())
}

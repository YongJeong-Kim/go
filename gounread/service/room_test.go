package service_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"gounread/api"
	"gounread/service"
)

var _ = Describe("Room", func() {
	var svc service.Servicer
	var user1, user2 string

	BeforeEach(func() {
		svc = service.NewService(api.NewSession())
		user1 = "05f84f46-d4ad-42db-af4f-da63cffcb721"
		user2 = "febba554-152e-496a-add5-31d0672fdc2a"
	})

	Context("room count", func() {
		It("user1, user2 room count", func() {
			rooms1 := svc.GetRoomsByUserID(user1)
			Expect(len(rooms1)).Should(BeEquivalentTo(2))
			rooms2 := svc.GetRoomsByUserID(user2)
			Expect(len(rooms2)).Should(BeEquivalentTo(2))
		})
		It("user not found. empty room", func() {
			rooms := svc.GetRoomsByUserID("sdfseoh")
			Expect(len(rooms)).Should(BeEquivalentTo(0))
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
})

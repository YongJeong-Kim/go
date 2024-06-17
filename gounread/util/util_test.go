package util_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"gounread/util"
)

var _ = Describe("Util", func() {
	type Str struct {
		Name string
		Age  int
	}

	DescribeTable("filter", func(samples []any, input any, matcher types.GomegaMatcher) {
		r := util.Filter(samples, func(s any) bool {

			switch s.(type) {
			case *Str:
				return input == s.(*Str).Name
			default:
				return s == input
			}
		})
		Expect(r).To(matcher)
	},
		Entry(
			"string filter",
			[]any{"a", "b"},
			"a",
			BeEquivalentTo([]any{"a"}),
		),
		Entry(
			"int filter",
			[]any{1, 2, 3, 4, 5},
			2,
			BeEquivalentTo([]any{2}),
		),
		Entry(
			"struct filter",
			[]any{
				&Str{
					Name: "aaa",
					Age:  11,
				},
				&Str{
					Name: "bbb",
					Age:  22,
				},
			},
			"bbb",
			BeEquivalentTo([]any{&Str{
				Name: "bbb",
				Age:  22,
			}}),
		),
	)
})

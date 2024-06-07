package service_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGounread(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gounread Suite")
}

package messages_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMessages(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Messages Suite")
}

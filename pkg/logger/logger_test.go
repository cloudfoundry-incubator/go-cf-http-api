package logger_test

import (
	"bytes"

	"../logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {
	It("prints logs", func() {
		var b bytes.Buffer
		logger.Init(nil, &b)

		logger.Out.Println("some-log")

		Expect(string(b.Bytes())).To(ContainSubstring("some-log"))
	})

	It("prints error logs", func() {
		var b bytes.Buffer
		logger.Init(&b, nil)

		logger.Err.Println("some-err")

		Expect(string(b.Bytes())).To(ContainSubstring("logger_test.go"))
		Expect(string(b.Bytes())).To(ContainSubstring("some-err"))
	})
})

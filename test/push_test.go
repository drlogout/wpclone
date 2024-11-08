package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run wpclone push", func() {
	Context("without flags", func() {
		BeforeEach(func() {
			_, err := setupLocalFiles("push")
			Expect(err).To(BeNil())

			err = startRemoteDocker()
			Expect(err).To(BeNil())
		})

		It("should push local files and db", func() {
			_, err := wpcloneCLI("pull")
			Expect(err).To(BeNil())

			err = modifyLocalWP()
			Expect(err).To(BeNil())

			_, err = wpcloneCLI("push", "--force")
			Expect(err).To(BeNil())

			ok, err := remoteHasContent("Modified locally")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
})

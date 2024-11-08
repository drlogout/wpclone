package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run wpclone pull", func() {
	Context("local without flags", func() {
		BeforeEach(func() {
			_, err := setupLocalFiles("pull")
			Expect(err).To(BeNil())

			err = startRemoteDocker()
			Expect(err).To(BeNil())
		})

		It("should pull remote files and db", func() {
			_, err := wpcloneCLI("pull", "--force")
			Expect(err).To(BeNil())

			ok, err := localHasContent("## REMOTE ##")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})

	// Context("docker without flags", func() {
	// 	BeforeEach(func() {
	// 		_, err := setupLocalFiles("pull-docker")
	// 		Expect(err).To(BeNil())

	// 		err = startRemoteDocker()
	// 		Expect(err).To(BeNil())
	// 	})

	// 	It("should pull remote files and db", func() {
	// 		_, err := wpcloneCLI("pull", "--force")
	// 		Expect(err).To(BeNil())

	// 		ok, err := localHasContent("## REMOTE ##")
	// 		Expect(err).To(BeNil())
	// 		Expect(ok).To(BeTrue())
	// 	})
	// })
})

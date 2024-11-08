package main_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run wpclone init", func() {
	Context("without flags", func() {
		var paths Paths

		BeforeEach(func() {
			var err error

			paths, err = setupLocalFiles()
			Expect(err).To(BeNil())
		})

		It("should create wpclone.yml", func() {
			_, err := wpcloneCLI("init")
			Expect(err).To(BeNil())

			_, err = os.Stat(paths.LocalConfigFile)
			os.IsNotExist(err)
			Expect(err).To(BeNil())

			cfg, err := parseConfig(paths.LocalConfigFile)
			Expect(err).To(BeNil())

			Expect(cfg.Local.Path).To(Equal("/home/user/local_wp"))
			Expect(cfg.Remote.Path).To(Equal("/var/www/html"))
		})
	})
})

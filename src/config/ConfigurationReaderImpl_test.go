package config

import (
	"path"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"chatServer/testhelpers"
)

func TestConfigurationReaderImpl(t *testing.T) {

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "ConfigurationReaderImpl unit test suite")
}

func createReader() ConfigurationReader {
	return NewReaderImpl()
}

var _ = ginkgo.Describe("ConfigurationReaderImpl", func() {

	ginkgo.Context("Read", func() {

		ginkgo.It("should read the json config file and return the configuration", func() {

			reader := createReader()
			cfg := reader.Read(path.Join(testhelpers.GetServerRootDir(), "/resources/config/config.json"))
			gomega.Expect(cfg.Host).To(gomega.Equal("localhost"))
			gomega.Expect(cfg.Port).To(gomega.Equal("9080"))
			gomega.Expect(cfg.ConnectionType).To(gomega.Equal("tcp"))
		})
	})
})
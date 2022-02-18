package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"scooter_micro/config"
)

var _ = Describe("Config", func() {

	Describe(".GetStringParameter() with cleared environment parameters", func() {

		BeforeEach(func() {
			os.Clearenv()
		})

		DescribeTable("environment parameters are clear",
			func(envParam, def string) {
				Expect(os.Getenv(envParam)).To(Equal(""))
				Expect(config.GetStringParameter(os.Getenv(envParam), def)).To(Equal(def))
			},
			EntryDescription("Calling GetStringParameter with ENV:%v Default:%v"),
			Entry(nil, "HTTP_PORT", "defaultParameter"),
			Entry(nil, "PG_HOST", "defaultParameter"),
			Entry(nil, "PG_PORT", "defaultParameter"),
		)
	})

	When(".GetStringParameter() with environment parameters", func() {
		BeforeEach(func() {
			err := os.Setenv("HTTP_PORT", "HTTP_PORT")
			Expect(err).NotTo(HaveOccurred())
			err = os.Setenv("PG_HOST", "PG_HOST")
			Expect(err).NotTo(HaveOccurred())
			err = os.Setenv("PG_PORT", "PG_PORT")
			Expect(err).NotTo(HaveOccurred())
		})

		DescribeTable("parameters have set",
			func(envParam, def string) {
				Expect(os.Getenv(envParam)).NotTo(BeNil())
				Expect(config.GetStringParameter(envParam, def)).To(Equal(envParam))
			},
			EntryDescription("Calling GetStringParameter with ENV:%v Default:%v"),
			Entry(nil, "HTTP_PORT", "defaultParameter"),
			Entry(nil, "PG_HOST", "defaultParameter"),
			Entry(nil, "PG_PORT", "defaultParameter"),
		)
	})
})

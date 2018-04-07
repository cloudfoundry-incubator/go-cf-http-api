package viewer_test

import (
	"../viewer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Viewer", func() {
	It("parses the template", func() {
		template := `<div>{{test}}</div>`
		expected := `<div>value</div>`

		parsedTemplate := viewer.Parse(template, map[string]string{"test": "value"})
		Expect(parsedTemplate).To(Equal(expected))
	})

	It("does not error if the metadata is nonexistent", func() {

		template := `<div>{{test}}</div>`
		expected := `<div>{{test}}</div>`

		parsedTemplate := viewer.Parse(template, map[string]string{})
		Expect(parsedTemplate).To(Equal(expected))
	})
})

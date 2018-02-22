package ocl_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/opsidian/ocl/ocl"
	parsley_parser "github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

func evaluate(filename string) (interface{}, *parsley_parser.History, reader.Error) {
	path := filepath.Join("fixtures", filename)
	reader, err := text.NewFileReader(path, true)
	Expect(err).ToNot(HaveOccurred())
	s := parsley.NewSentence(ocl.NewParser())
	return s.Evaluate(reader, nil)
}

func evaluateText(tmpl string) (interface{}, *parsley_parser.History, reader.Error) {
	reader := text.NewReader([]byte(tmpl), "", true)
	s := parsley.NewSentence(ocl.NewParser())
	return s.Evaluate(reader, nil)
}

var _ = Describe("OCL Parser", func() {
	It("should parse the full syntax from file", func() {
		res, _, err := evaluate("complete.ocl")
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"var_int":         1,
			"var_float":       1.2,
			"var_string":      "string",
			"var_bool_false":  false,
			"var_bool_true":   true,
			"var_array":       []interface{}{1, "string2"},
			"var_empty_array": []interface{}{},
			"var_empty_obj":   nil,
			"var_flint_value": 3,
			"var_obj": map[string]interface{}{
				"var_int":        2,
				"var_float":      2.3,
				"var_string":     "string3",
				"var_bool_false": false,
				"var_bool_true":  true,
				"var_array":      []interface{}{3, "string4"},
				"var_empty_obj":  nil,
			},
			"var_obj2": map[string]interface{}{
				"a": 1,
			},
		}))
	})

	It("should parse an empty file", func() {
		res, _, err := evaluateText("")
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{}))
	})

	It("should parse an array", func() {
		res, _, err := evaluateText(`
			a = []
			b = [1]
			c = [1, 2]
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": []interface{}{},
			"b": []interface{}{1},
			"c": []interface{}{1, 2},
		}))
	})

	It("should parse an object", func() {
		res, _, err := evaluateText(`
			a = {}
			b = {
			}
			c = {
				d = 1
			}
			e = {
				f = 2
				g = 3
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": nil,
			"b": nil,
			"c": map[string]interface{}{"d": 1},
			"e": map[string]interface{}{"f": 2, "g": 3},
		}))
	})

	It("should parse an object with multiple keys", func() {
		res, _, err := evaluateText(`
			a b {
				c = 1
				d = 2
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": map[string]interface{}{"b": map[string]interface{}{"c": 1, "d": 2}},
		}))
	})

	It("should throw an error if a key is empty", func() {
		_, _, err := evaluateText(`
			a "" {}
		`)
		Expect(err).To(MatchError("key can not be empty at 2:6"))
	})

	It("should merge objects", func() {
		res, _, err := evaluateText(`
			a a1 {
				b = 1
			}
			a a2 {
				c = 2
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": map[string]interface{}{
				"a1": map[string]interface{}{"b": 1},
				"a2": map[string]interface{}{"c": 2},
			},
		}))
	})

	It("should throw an error if multikeys are defined twice", func() {
		_, _, err := evaluateText(`
			a a1 {}
			a a1 {}
		`)
		Expect(err).To(MatchError("key was already defined at 3:6"))
	})

	It("should parse an array append", func() {
		res, _, err := evaluateText(`
			a[] = 1
			a[] = 2
			b[] {
				c = 1
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": []interface{}{1, 2},
			"b": []interface{}{map[string]interface{}{"c": 1}},
		}))
	})

	It("should throw an error if you try to append to a non-array type", func() {
		_, _, err := evaluateText(`
			a = 1
			a[] = 2
		`)
		Expect(err).To(MatchError("can not append to not array type at 3:4"))
	})

	It("should throw an error if a non-last key has []", func() {
		_, _, err := evaluateText(`
			a[] b {}
		`)
		Expect(err).To(MatchError("only the last key can end with [] at 2:4"))
	})

	It("should parse an integer", func() {
		res, _, err := evaluateText(`
			a = 1
			b = [2]
			c {
				d = 3
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": 1,
			"b": []interface{}{2},
			"c": map[string]interface{}{"d": 3},
		}))
	})

	It("should parse a float", func() {
		res, _, err := evaluateText(`
			a = 1.1
			b = [2.1]
			c {
				d = 3.1
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": 1.1,
			"b": []interface{}{2.1},
			"c": map[string]interface{}{"d": 3.1},
		}))
	})

	It("should parse a string", func() {
		res, _, err := evaluateText(`
			a = "a1"
			b = ["b1"]
			c {
				d = "d1"
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": "a1",
			"b": []interface{}{"b1"},
			"c": map[string]interface{}{"d": "d1"},
		}))
	})

	It("should parse a boolean", func() {
		res, _, err := evaluateText(`
			a = false
			b = [true]
			c {
				d = false
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": false,
			"b": []interface{}{true},
			"c": map[string]interface{}{"d": false},
		}))
	})

	It("should parse a flint expression", func() {
		res, _, err := evaluateText(`
			a = {{ 1 + 2 }}
			b = [{{ 2 + 3 }}]
			c {
				d = {{ 3 + 4 }}
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": 3,
			"b": []interface{}{5},
			"c": map[string]interface{}{"d": 7},
		}))
	})

	It("should parse a quoted flint expression", func() {
		res, _, err := evaluateText(`
			a = "{{ 1 + 2 }}"
			b = ["{{ 2 + 3 }}"]
			c {
				d = "{{ 3 + 4 }}"
			}
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": "3",
			"b": []interface{}{"5"},
			"c": map[string]interface{}{"d": "7"},
		}))
	})

	It("should convert flint types to string", func() {
		res, _, err := evaluateText(`
			a = "{{ 1 }}"
			b = "{{ 1.2 }}"
			c = "{{ true }}"
			d = "{{ "str" }}"
			e = "{{ nil }}"
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": "1",
			"b": "1.2",
			"c": "true",
			"d": "str",
			"e": "nil",
		}))
	})

	It("should ignore comments", func() {
		res, _, err := evaluateText(`
// this is a comment
# this is also a comment
/*
and this is too
*/
a = 1
		`)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(map[string]interface{}{
			"a": 1,
		}))
	})

	Context("a variable is defined multiple times", func() {
		It("should throw an error", func() {
			_, _, err := evaluateText(`
				a = 1
				a = 2
			`)
			Expect(err).To(MatchError("key was already defined at 3:5"))
		})
	})

	Context("an object key is defined multiple times", func() {
		It("should throw an error", func() {
			_, _, err := evaluateText(`
				a {
					b = 1
					b = 2
				}
			`)
			Expect(err).To(MatchError("key was already defined at 4:6"))
		})
	})

	Context("evaluation error", func() {
		It("should throw an error when a Flint expression fails to evaluate", func() {
			_, _, err := evaluateText(`
				a = {{ 1 / 0 }}
			`)
			Expect(err).To(MatchError("divison by zero at 2:16"))
		})

		It("should throw an error when a quoted Flint expression fails to evaluate", func() {
			_, _, err := evaluateText(`
				a = "{{ 1 / 0 }}"
			`)
			Expect(err).To(MatchError("divison by zero at 2:17"))
		})

		It("should throw an error when a quoted Flint expression a non-printable type", func() {
			_, _, err := evaluateText(`
				a = "{{ [1, 2] }}"
			`)
			Expect(err).To(MatchError("an array can not be converted to string at 2:14"))
		})

		It("should throw an error when a Flint expression fails to evaluate in an array", func() {
			_, _, err := evaluateText(`
				a = [{{ 1 / 0 }}]
			`)
			Expect(err).To(MatchError("divison by zero at 2:17"))
		})

		It("should throw an error when a Flint expression fails to evaluate in an object", func() {
			_, _, err := evaluateText(`
				a = {
					b = {{ 1 / 0 }}
				}
			`)
			Expect(err).To(MatchError("divison by zero at 3:17"))
		})
	})

	Context("invalid syntax", func() {
		It("should throw an error when a variable doesn't have a value", func() {
			_, _, err := evaluateText(`
				a =
				b = 2
			`)
			Expect(err).To(MatchError("Failed to parse the input: was expecting value at 2:8"))
		})

		It("should throw an error when a bracket is not closed", func() {
			_, _, err := evaluateText(`
				a = {
				b = 2
			`)
			Expect(err).To(MatchError("Failed to parse the input: was expecting '}' at 4:1"))
		})

		It("should throw an error for an unexpected value", func() {
			_, _, err := evaluateText(`
				a = b
			`)
			Expect(err).To(MatchError("Failed to parse the input: was expecting value at 2:8"))
		})
	})

})

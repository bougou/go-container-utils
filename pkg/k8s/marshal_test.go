package k8s

import (
	"testing"
)

func TestYAMLMarshal(t *testing.T) {
	tests := []struct {
		expected string
		object   interface{}
	}{
		{
			expected: `age: 20
name: tom
nested:
  age: 20
  name: tom
  nestedL:
  - age: 20
    name: tom
  - age: 20
    name: tom
    nestedL:
    - age: 20
      name: tom
    - age: 20
      name: tom
`,

			object: map[string]interface {
			}{
				"status":            "",
				"name":              "tom",
				"age":               20,
				"creationTimestamp": nil,
				"nested": map[string]interface{}{
					"name":              "tom",
					"age":               20,
					"creationTimestamp": nil,
					"nestedL": []map[string]interface{}{
						{
							"name":              "tom",
							"age":               20,
							"creationTimestamp": nil,
						},
						{
							"name":              "tom",
							"age":               20,
							"creationTimestamp": nil,
							"nestedL": []map[string]interface{}{
								{
									"name":              "tom",
									"age":               20,
									"creationTimestamp": nil,
								},
								{
									"name":              "tom",
									"age":               20,
									"creationTimestamp": nil,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		b, err := YAMLMarshal(tt.object,
			WithTransform(TransformRemoveCreationTimestampRecursive),
			WithTransform(TransformRemoveStatus),
		)
		if err != nil {
			t.Error(err)
			return
		}

		if tt.expected != string(b) {
			t.Error("not matched")
		}
	}
}

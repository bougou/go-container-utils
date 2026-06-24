package k8s

// Ref: https://github.com/kubernetes-sigs/controller-tools/blob/release-0.15/pkg/genall/genall.go#L129

import (
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"
)

// Usage:
//
//	b, err := YAMLMarshal(vm,
//		WithTransform(TransformRemoveCreationTimestampRecursive),
//		WithTransform(TransformRemoveStatus))
//	if err != nil {
//	}
func YAMLMarshal(o interface{}, options ...*TransformYAMLOption) ([]byte, error) {
	j, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("error marshaling into JSON: %v", err)
	}

	return yamlJSONToYAMLWithFilter(j, options...)
}

// yamlJSONToYAMLWithFilter is based on sigs.k8s.io/yaml.JSONToYAML, but allows for transforming the final data before writing.
func yamlJSONToYAMLWithFilter(j []byte, options ...*TransformYAMLOption) ([]byte, error) {
	// Convert the JSON to an object.
	var jsonObj map[string]interface{}
	// We are using yaml.Unmarshal here (instead of json.Unmarshal) because the
	// Go JSON library doesn't try to pick the right number type (int, float,
	// etc.) when unmarshalling to interface{}, it just picks float64
	// universally. go-yaml does go through the effort of picking the right
	// number type, so we can preserve number type throughout this process.
	if err := yaml.Unmarshal(j, &jsonObj); err != nil {
		return nil, err
	}

	for _, option := range options {
		if option.transform != nil {
			if err := option.transform(jsonObj); err != nil {
				return nil, err
			}
		}
	}

	// Marshal this object into YAML.
	return yaml.Marshal(jsonObj)
}

type TransformYAMLOption struct {
	transform func(obj map[string]interface{}) error
}

// WithTransform applies a transformation to objects just before writing them.
func WithTransform(transform func(obj map[string]interface{}) error) *TransformYAMLOption {
	return &TransformYAMLOption{
		transform: transform,
	}
}

// TransformRemoveStatus remove status field.
func TransformRemoveStatus(obj map[string]interface{}) error {
	delete(obj, "status")
	return nil
}

// TransformRemoveCreationTimestamp ensures we do not write the metadata.creationTimestamp field.
func TransformRemoveCreationTimestamp(obj map[string]interface{}) error {
	metadata := obj["metadata"].(map[string]interface{})
	delete(metadata, "creationTimestamp")
	return nil
}

// TransformRemoveCreationTimestampRecursive removes all fields "creationTimestamp" recusively.
func TransformRemoveCreationTimestampRecursive(obj map[string]interface{}) error {
	return RemoveFieldRecusive(obj, "creationTimestamp")
}

func TransformRemoveCreationTimestampRecursiveL(objs []interface{}) error {
	return RemoveFieldRecusiveL(objs, "creationTimestamp")
}

func RemoveFieldRecusive(obj map[string]interface{}, field string) error {
	for k, v := range obj {
		if k == field {
			delete(obj, k)
			continue
		}
		if vAsM, ok := v.(map[string]interface{}); ok {
			if err := RemoveFieldRecusive(vAsM, field); err != nil {
				return err
			}
			obj[k] = vAsM
			continue
		}
		if vAsL, ok := v.([]interface{}); ok {
			if err := RemoveFieldRecusiveL(vAsL, field); err != nil {
				return err
			}
			obj[k] = vAsL
			continue
		}
	}
	return nil
}

func RemoveFieldRecusiveL(objs []interface{}, field string) error {
	for i, e := range objs {
		if eAsM, ok := e.(map[string]interface{}); ok {
			if err := RemoveFieldRecusive(eAsM, field); err != nil {
				return err
			}
			objs[i] = eAsM
			continue
		}

		if eAsL, ok := e.([]interface{}); ok {
			if err := RemoveFieldRecusiveL(eAsL, field); err != nil {
				return err
			}
			objs[i] = eAsL
			continue
		}
	}
	return nil
}

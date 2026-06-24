package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	k8syaml "sigs.k8s.io/yaml"
)

// Ref:
// - https://github.com/kubernetes/client-go/issues/193

func UnmarshalToUnstructured(objYaml string) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := k8syaml.Unmarshal([]byte(objYaml), &u); err != nil {
		return nil, err
	}

	return u, nil
}

// UnmarshalToRuntimeObject unmarshal yaml to runtime.Object.
//
// after you got the returned runObj runtime.Object, you can apply type switch to runObj to get concrete type, eg:
//
//	 switch obj := runObj.(type) {
//		case *corev1.Service:
//			clusterIP := obj.Spec.ClusterIP
//			port := obj.Spec.Ports[0].Port
//		default:
//	}
//
// or apply type assertion, eg:
//
//	if obj, ok := runObj.(*corev1.Service); ok {
//	}
func UnmarshalToRuntimeObject(objYaml string) (runtime.Object, *schema.GroupVersionKind, error) {
	decode := scheme.Codecs.UniversalDeserializer()
	return decode.Decode([]byte(objYaml), nil, nil)
}

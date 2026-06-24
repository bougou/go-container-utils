package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var getPort = func(u *unstructured.Unstructured) interface{} {
	return u.Object["spec"].(map[string]interface{})["ports"].([]interface{})[0].(map[string]interface{})["port"]
}

var jsonStr = `{
	"apiVersion": "v1",
	"kind": "Service",
	"metadata": {
		"name": "test-svc"
	},
	"spec": {
		"ports": [
			{
				"port": 8080,
				"protocol": "UDP"
			}
		]
	}
}`

func Test_UnmarshalToRuntimeObject(t *testing.T) {
	runObj, gvk, err := UnmarshalToRuntimeObject(jsonStr)
	if err != nil {
		t.Error(err)
	}

	switch runObj.(type) {
	case *corev1.Service:
	default:
		t.Error("not corev1.Service type")
	}

	assert.Equal(t, "Service", gvk.Kind, "gvk.Kind is not 'Service'")

	if obj, ok := runObj.(*corev1.Service); ok {
		assert.Equal(t, "Service", obj.Kind, "not equal")
	}
}

func Test_UnmarshalToUnstructured(t *testing.T) {
	u, err := UnmarshalToUnstructured(jsonStr)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, int64(8080), getPort(u).(int64), "not equal")
}

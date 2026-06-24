package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8syaml "sigs.k8s.io/yaml"
)

// apimachinery/pkg defines
//
//	type NamespacedName struct {
//		Namespace string
//		Name      string
//	}
//
// controller-runtime defines
//
//	type ObjectKey = types.NamespacedName
//
// They are the same thing.
func WithNamespacedName(ns string, name string) types.NamespacedName {
	return types.NamespacedName{
		Namespace: ns,
		Name:      name,
	}
}

func GetNamespacedName(k8sObjYaml string) (types.NamespacedName, error) {
	o := types.NamespacedName{}

	type metaObj struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`
	}

	someObj := &metaObj{}
	if err := k8syaml.Unmarshal([]byte(k8sObjYaml), someObj); err != nil {
		return o, err
	}

	// apiVersion := someObj.APIVersion
	// kind := someObj.Kind

	name := someObj.ObjectMeta.Name
	namespace := someObj.ObjectMeta.Namespace

	return types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, nil
}

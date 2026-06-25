package resourcequantity

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
)

func Test_Quantity_String(t *testing.T) {
	s := "50000Mi"
	quantity, err := resource.ParseQuantity(s)
	if err != nil {
		t.Error(err)
	}

	got := showQuantity(quantity)
	want := `
String: 50000Mi
Value: 52428800000
Scale(Kilo): 52428800
Scale(Mega): 52429
Scale(Giga): 53
Scale(Tera): 1
Scale(Peta): 1
Scale(Exa): 1

IEC(Ki): 51200000
IEC(Mi): 50000
IEC(Gi): 48
IEC(Ti): 0
IEC(Pi): 0
IEC(Ei): 0
`
	if got != want {
		t.Errorf("showQuantity() mismatch\ngot:\n%q\nwant:\n%q", got, want)
	}

}

func showQuantity(quantity resource.Quantity) string {
	return fmt.Sprintf(`
String: %s
Value: %d
Scale(Kilo): %d
Scale(Mega): %d
Scale(Giga): %d
Scale(Tera): %d
Scale(Peta): %d
Scale(Exa): %d

IEC(Ki): %d
IEC(Mi): %d
IEC(Gi): %d
IEC(Ti): %d
IEC(Pi): %d
IEC(Ei): %d
`,

		quantity.String(),
		quantity.Value(),
		quantity.ScaledValue(resource.Kilo),
		quantity.ScaledValue(resource.Mega),
		quantity.ScaledValue(resource.Giga),
		quantity.ScaledValue(resource.Tera),
		quantity.ScaledValue(resource.Peta),
		quantity.ScaledValue(resource.Exa),
		quantity.Value()/1024,
		quantity.Value()/1024/1024,
		quantity.Value()/1024/1024/1024,
		quantity.Value()/1024/1024/1024/1024,
		quantity.Value()/1024/1024/1024/1024/1024,
		quantity.Value()/1024/1024/1024/1024/1024/1024,
	)
}

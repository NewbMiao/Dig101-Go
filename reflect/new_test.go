package reflectUtils

import (
	"testing"

	"github.com/NewbMiao/Dig101-Go/reflect/testData"
)

func TestSetStructPtrUnExportedStrField(t *testing.T) {
	var eg testData.Example

	err := SetStructPtrUnExportedStrField(&eg, "a", "test")
	if err != nil {
		t.Fatal(err)
	}
	if GetStructPtrUnExportedField(&eg, "a").String() != "test" {
		t.Errorf("SetStructPtrUnExportedStrField failed:\t GetStructPtrUnExportedField got %+v", eg)
	}
	if fieldVal, _ := GetStructUnExportedField(eg, "a"); fieldVal.String() != "test" {
		t.Errorf("SetStructPtrUnExportedStrField failed:\t GetStructUnExportedField got %+v", eg)
	}

}

func TestSetStructUnExportedStrField(t *testing.T) {
	var eg testData.Example
	// did not change eg, just create newV to modify
	newV, err := SetStructUnExportedStrField(eg, "a", "test")
	if err != nil {
		t.Fatal(err)
	}
	if newV.FieldByName("a").String() != "test" {
		t.Errorf("SetStructPtrUnExportedStrField failed:\t got %+v", eg)
	}
}

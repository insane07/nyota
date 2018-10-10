package handler

import (
	"nyota/backend/model"
	"testing"
)

func TestPrepNewUnclassifiedData(t *testing.T) {
	var newUnclassifiedDataArray model.APIDataArray
	result := prepNewUnclassifiedData(newUnclassifiedDataArray)
	if len(result) != len(newUnclassifiedDataArray) {
		t.Errorf("Expecting same count")
	}
	if len(result) != 0 {
		t.Errorf("Expecting count == 0")
	}

	apidata := model.KeySortData(11, 111)
	newUnclassifiedDataArray = append(newUnclassifiedDataArray, apidata)
	result = prepNewUnclassifiedData(newUnclassifiedDataArray)
	if len(result) != 1 {
		t.Errorf("Expecting count == 1")
	}
	if result[0].Value != 111 {
		t.Errorf("Expecting Value == 111, returned %d", result[0].Value)
	}

	apidata = model.KeySortData(12, 114)
	newUnclassifiedDataArray = append(newUnclassifiedDataArray, apidata)
	result = prepNewUnclassifiedData(newUnclassifiedDataArray)
	if len(result) != 2 {
		t.Errorf("Expecting count == 2")
	}
	if result[1].Value != 3 {
		t.Errorf("Expecting Value == 3, returned %d", result[1].Value)
	}

	apidata = model.KeySortData(13, 118)
	newUnclassifiedDataArray = append(newUnclassifiedDataArray, apidata)
	result = prepNewUnclassifiedData(newUnclassifiedDataArray)
	if len(result) != 2 {
		t.Errorf("Expecting count == 2")
	}
	if result[1].Value != 4 {
		t.Errorf("Expecting Value == 4, returned %d", result[1].Value)
	}

	apidata = model.KeySortData(14, 114)
	newUnclassifiedDataArray = append(newUnclassifiedDataArray, apidata)
	result = prepNewUnclassifiedData(newUnclassifiedDataArray)
	if len(result) != 3 {
		t.Errorf("Expecting count == 3")
	}
	if result[1].Value != 4 {
		t.Errorf("Expecting Value == 4, returned %d", result[1].Value)
	}
	if result[2].Value != 0 {
		t.Errorf("Expecting Value == 0, returned %d", result[2].Value)
	}

	apidata = model.KeySortData(15, 146)
	newUnclassifiedDataArray = append(newUnclassifiedDataArray, apidata)

	apidata = model.KeySortData(16, 122)
	newUnclassifiedDataArray = append(newUnclassifiedDataArray, apidata)

	apidata = model.KeySortData(17, 101)
	newUnclassifiedDataArray = append(newUnclassifiedDataArray, apidata)

	result = prepNewUnclassifiedData(newUnclassifiedDataArray)
	if len(result) != 6 {
		t.Errorf("Expecting count == 6")
	}
	if result[5].Value != 0 {
		t.Errorf("Expecting Value == 0, returned %d", result[5].Value)
	}
}

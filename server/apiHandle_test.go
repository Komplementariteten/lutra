package server

import (
	"encoding/base64"
	"testing"
)

func Test_Base64(t *testing.T) {

	str1 := "/9j/4AAQSkZJRgABAQAASABIAAD/4QDKRXhpZgAATU0AKgAAAAgABwESAAMAAAABAAEAAAEaAAUAAAABAAAAYgEbAAUAAAABAAAAagEoAAMAAAABAAIAAAExAAIAAAARAAAAcgEyAAIAAAAUAAAAhIdpAAQAAAABAAAAmAAAAAAAAABIAAAAAQAAAEgAAAABUGl4ZWxtYXRvciAzLjkuMgAAMjAyMTowNDoxNCAyMjowNDozMwAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAA9aADAAQAAAABAAAAWAAAAAD/4QlvaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wLwA8P3hwYWNrZXQgYmVnaW49Iu+7vyIgaWQ9Ilc1TTBNcENlaGlIenJlU3pOVGN6a2M5ZCI/PiA8eDp4bXBtZXRhIHhtbG5zOng9ImFkb2JlOm5zOm1ldGEvIiB4OnhtcHRrPSJYTVAgQ29yZSA2LjAuMCI+IDxyZGY6UkRGIHhtbG5zOnJkZj0iaHR0cDovL3d3dy53My5vcmcvMTk5OS8wMi8yMi1yZGYtc3ludGF4LW5zIyI+IDxyZGY6RGVzY3JpcHRpb24gcmRmOmFib3V0PSIiIHhtbG5zOnhtcD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wLyIgeG1wOk1vZGlmeURhdGU9IjIwMjEtMDQtMTRUMjI6MDQ6MzMiLz4gPC9yZGY6UkRGPiA8L3g6eG1wbWV0YT4gICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8P3hwYWNrZXQgZW5kPSJ3Ij8+AP/AABEIAFgA9QMBIgACEQEDEQH/xAAfAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgv/xAC1EAAC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"
	str2 := "/9j/4AAQSkZJRgABAQAASABIAAD/4QDKRXhpZgAATU0AKgAAAAgABwESAAMAAAABAAEAAAEaAAUAAAABAAAAYgEbAAUAAAABAAAAagEoAAMAAAABAAIAAAExAAIAAAARAAAAcgEyAAIAAAAUAAAAhIdpAAQAAAABAAAAmAAAAAAAAABIAAAAAQAAAEgAAAABUGl4ZWxtYXRvciAzLjkuMgAAMjAyMTowNDoxNCAyMjowNDozMwAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAA9aADAAQAAAABAAAAWAAAAAD/4QlvaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wLwA8P3hwYWNrZXQgYmVnaW49Iu+7vyIgaWQ9Ilc1TTBNcENlaGlIenJlU3pOVGN6a2M5ZCI/PiA8eDp4bXBtZXRhIHhtbG5zOng9ImFkb2JlOm5zOm1ldGEvIiB4OnhtcHRrPSJYTVAgQ29yZSA2LjAuMCI+IDxyZGY6UkRGIHhtbG5zOnJkZj0iaHR0cDovL3d3dy53My5vcmcvMTk5OS8wMi8yMi1yZGYtc3ludGF4LW5zIyI+IDxyZGY6RGVzY3JpcHRpb24gcmRmOmFib3V0PSIiIHhtbG5zOnhtcD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wLyIgeG1wOk1vZGlmeURhdGU9IjIwMjEtMDQtMTRUMjI6MDQ6MzMiLz4gPC9yZGY6UkRGPiA8L3g6eG1wbWV0YT4gICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8P3hwYWNrZXQgZW5kPSJ3Ij8+AP/AABEIAFgA9QMBIgACEQEDEQH/xAAfAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgv/xAC1EAAC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"

	_, err := base64.StdEncoding.DecodeString(str1)
	if err != nil {
		t.Log(err.Error())
	}

	_, err = base64.StdEncoding.DecodeString(str2)
	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestApiHandle_CheckJson_VaildRequest(t *testing.T) {
	valid_json := `{"Type":"testdata", "Content":{"a":12,"b":13}}`
	var response ApiRequestResponse
	content := checkJson([]byte(valid_json), &response)
	if content == nil {
		t.Fatal("Content could not be read")
	}
	if response.Error {
		t.Fatal("Valid json returns error")
	}
}

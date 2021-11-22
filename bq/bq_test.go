package bq

import (
	"errors"
	"testing"
)

func TestNotExistErrorChecker(t *testing.T) {
	not404 := errors.New("Token is expired")
	b := isNotExist(not404)
	if b {
		t.Fatalf("%v is Not Found Error", not404)
	}
	notFound := errors.New("This is the error of non-existing dataset xyz, googleapi: Error 404: Not found: Dataset acuit-data-lake-dev:xyz, notFound")
	b = isNotExist(notFound)
	if !b {
		t.Fatalf("%v is Found Error", notFound)
	}
}

package yesql_test

import (
	"testing"

	"github.com/izolate/yesql"
)

func TestFoo(t *testing.T) {
	db, err := yesql.Open("postgres", "host=localhost")
	if err != nil {
		t.Fatalf(err.Error())
	}
	db.DB.Begin()
}

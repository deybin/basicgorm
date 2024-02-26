package test

import (
	"fmt"
	"testing"

	"github.com/deybin/basicgorm/test/table"
)

func TestModels(t *testing.T) {

	sucursal := new(table.Sucursal).New()

	fmt.Println(string(sucursal.GetTableName()))

	r := ""

	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

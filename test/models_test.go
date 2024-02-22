package test

import (
	"fmt"
	"testing"

	"github.com/deybin/basicgorm/test/table"
)

type Cliente struct {
	Name      string `json:"name" basicgorm:"lowerCase;max:20"`
	LastName  string `json:"last_name" basicgorm:"lowerCase;max:80"`
	BirthDate string `json:"birth_date" basicgorm:"date"`
	Year      uint64 `json:"year"`
}

func TestModels(t *testing.T) {

	sucursal, _ := table.GetSucursal()

	fmt.Println(string(sucursal[0].Type))

	r := ""

	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

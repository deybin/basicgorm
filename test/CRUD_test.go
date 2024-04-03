package test

import (
	"testing"

	"github.com/deybin/basicgorm"
	"github.com/deybin/basicgorm/test/table"
)

func TestCRUD_Single(t *testing.T) {
	dataInsert := map[string]interface{}{
		"c_sucu": "003",
		"l_sucu": "sucursal de prueba",
		"l_dire": "sin información",
	}

	crud := basicgorm.SqlExecSingle{}

	err := crud.New(new(table.Sucursal).New(), dataInsert).Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Exec("new_capital")
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

func TestCRUD_Single_Update(t *testing.T) {
	dataInsert := map[string]interface{}{
		"l_alma": "principal",
		"where":  map[string]interface{}{"c_sucu": "001", "c_alma": "002"},
	}

	crud := basicgorm.SqlExecSingle{}
	err := crud.New(new(table.Store).New(), dataInsert).Update()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Exec("new_capital")
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

func TestCRUD_Multiple(t *testing.T) {
	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "004",
		"l_sucu": "sucursal de prueba",
		"l_dire": "sin información",
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "005",
		"l_sucu": "sucursal de prueba",
		"l_dire": "sin información",
	})
	// dataInsert = append(dataInsert, map[string]interface{}{
	// 	"c_sucu": "003",
	// 	"l_sucu": "sucursal de prueba",
	// 	"l_dire": "sin información",
	// })

	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "004",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})

	crud := basicgorm.SqlExecMultiple{}
	crud.New("new_capital")

	trSucursal := crud.SetInfo(new(table.Sucursal).New(), dataInsert...)
	trAlmacen := crud.SetInfo(new(table.Store).New(), dataInsertAlma...)
	err := trAlmacen.Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	err = trSucursal.Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Exec()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

func TestCRUD_Multiple_transaction(t *testing.T) {
	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "004",
		"l_sucu": "sucursal de prueba",
		"l_dire": "indefinido",
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "005",
		"l_sucu": "sucursal de prueba",
		"l_dire": "indefinido",
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "003",
		"l_sucu": "sucursal de prueba",
		"l_dire": "indefinido",
	})

	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "004",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})

	// schema, tableName := table.GetSucursal()
	// schemaAlma, tableNameAlma := table.GetAlmacen()
	crud := new(basicgorm.SqlExecMultiple).New("new_capital")
	TransactionAlmacen := crud.New("new_capital").SetInfo(new(table.Store).New(), dataInsertAlma...)
	TransactionSucursal := crud.SetInfo(new(table.Sucursal).New(), dataInsert...)

	err := TransactionAlmacen.Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	err = TransactionSucursal.Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.ExecTransaction(TransactionAlmacen)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.ExecTransaction(TransactionSucursal)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Commit()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
}

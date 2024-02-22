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
	}
	schema, tableName := table.GetSucursal()
	crud := basicgorm.SqlExecSingle{}

	err := crud.New(tableName, dataInsert).Insert(schema)
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
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "005",
		"l_sucu": "sucursal de prueba",
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "003",
		"l_sucu": "sucursal de prueba",
	})

	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "004",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})

	schema, tableName := table.GetSucursal()
	schemaAlma, tableNameAlma := table.GetAlmacen()
	crud := basicgorm.SqlExecMultiple{}
	crud.New("new_capital")

	trAlmacen := crud.SetInfo(tableNameAlma, dataInsertAlma...)
	trSucursal := crud.SetInfo(tableName, dataInsert...)
	err := trAlmacen.Insert(schemaAlma)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	err = trSucursal.Insert(schema)
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
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "005",
		"l_sucu": "sucursal de prueba",
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "003",
		"l_sucu": "sucursal de prueba",
	})

	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "004",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})

	schema, tableName := table.GetSucursal()
	schemaAlma, tableNameAlma := table.GetAlmacen()
	crud := new(basicgorm.SqlExecMultiple).New("new_capital")
	TransactionAlmacen := crud.New("new_capital").SetInfo(tableNameAlma, dataInsertAlma...)
	TransactionSucursal := crud.SetInfo(tableName, dataInsert...)

	err := TransactionAlmacen.Insert(schemaAlma)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	err = TransactionSucursal.Insert(schema)
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

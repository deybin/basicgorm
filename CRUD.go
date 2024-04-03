package basicgorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type SqlExecSingle struct {
	ob     []map[string]interface{} //datos para observación
	data   []map[string]interface{} //datos para insertar o actualizar o eliminar
	query  []map[string]interface{}
	schema Schema
	action string
}

type SqlExecMultiple struct {
	tx          *sql.Tx
	database    string
	transaction []*Transaction
}

type Transaction struct {
	ob     []map[string]interface{} //datos para observación
	data   []map[string]interface{} //datos para insertar o actualizar o eliminar
	query  []map[string]interface{}
	schema Schema
	action string
}

/*
New crea una nueva instancia de SqlExecSingle con el esquema y los datos proporcionados.

	Parámetros
		* s {Schema}: esquema de la tabla
		* datos {[]map[string]interface{}}: datos a insertar, actualizar o eliminar

	Return
		- (*SqlExecSingle) retorna  puntero *SqlExecSingle struct
*/
func (sq *SqlExecSingle) New(s Schema, datos ...map[string]interface{}) *SqlExecSingle {
	sq.ob = datos
	sq.schema = s
	return sq
}

/*
Valida los datos para insertar y crea el query para insertar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *SqlExecSingle) Insert() error {
	sqlExec, data_insert, err := _insert(sq.schema.GetTableName(), sq.ob, sq.schema.GetSchemaInsert())
	if err != nil {
		return err
	}
	sq.query = sqlExec
	sq.data = data_insert
	sq.action = "INSERT"
	return nil
}

/*
Valida los datos para actualizar y crea el query para actualizar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *SqlExecSingle) Update() error {
	sqlExec, data_update, err := _update(sq.schema.GetTableName(), sq.ob, sq.schema.GetSchemaUpdate())
	if err != nil {
		return err
	}
	sq.query = sqlExec
	sq.data = data_update
	sq.action = "UPDATE"
	return nil
}

/*
Valida los datos para Eliminar y crea el query para Eliminar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *SqlExecSingle) Delete() error {
	sqlExec, data_delete, err := _delete(sq.schema.GetTableName(), sq.ob, sq.schema.GetSchemaDelete())
	if err != nil {
		return err
	}
	sq.query = sqlExec
	sq.data = data_delete
	sq.action = "DELETE"
	return nil
}

/*
Retorna los datos que se enviaron o enviaran para ser insertados, modificados o eliminados

	Return
		- []map[string]interface{}
*/
func (sq *SqlExecSingle) GetData() []map[string]interface{} {
	return sq.data
}

/*
Ejecuta el query

	Return
		- returns {error}: retorna errores ocurridos durante la ejecución
*/
func (sq *SqlExecSingle) Exec(database string, params ...bool) error {
	cnn, err := Connection(database)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err_cnn := cnn.PingContext(ctx)
	if err_cnn != nil {
		return errors.New(fmt.Sprint("Error Sql PING: ", err_cnn))
	}
	cross := false
	if len(params) == 1 {
		cross = params[0]
	}
	dataExec := sq.query
	defer cnn.Close()
	for _, item := range dataExec {
		sqlPre := item["sqlPreparate"].(string)
		if cross {
			if sq.action == "UPDATE" {
				sqlPre = Query_Cross_Update(sqlPre)
			}
		}
		// fmt.Println("PREPARED: ", sqlPre)
		stmt, err_prepare := cnn.Prepare(sqlPre)
		if err_prepare != nil {
			return fmt.Errorf("error sql prepare: %s ", err_prepare.Error())
		}
		valuesExec := item["valuesExec"].([]interface{})
		_, err_exec := stmt.Exec(valuesExec...)
		if err_exec != nil {
			return fmt.Errorf("error sql %s: %s", sq.action, err_exec.Error())
		}
	}
	return nil
}

/*
Crea una nueva instancia de SqlExecMultiple con el nombre de la base de datos proporcionado.

	Parámetros
	  * name {string}: database
	Return
	  - (*SqlExecMultiple) retorna  puntero *SqlExecMultiple struct
*/
func (sq *SqlExecMultiple) New(database string) *SqlExecMultiple {
	sq.database = database
	// sq.transaction = make(map[string]*Transaction)
	return sq
}

/*
SetInfo establece la información para una nueva transacción en SqlExecMultiple.

	Recibe un esquema (s Schema) y datos (datos ...map[string]interface{}) para la transacción.
	Retorna un puntero a la transacción creada.
	Parámetros
		* s {Schema}: esquema de la tabla
		* datos {[]map[string]interface{}}: datos a insertar, actualizar o eliminar
	Return
		- (*Transaction) retorna puntero *Transaction
*/
func (sq *SqlExecMultiple) SetInfo(s Schema, datos ...map[string]interface{}) *Transaction {
	key := len(sq.transaction)
	sq.transaction = append(sq.transaction, &Transaction{
		ob:     datos,
		schema: s,
	})

	return sq.transaction[key]
}

/*
*
Ejecuta el query

	Return
		- (error): retorna errores ocurridos durante la ejecución
*/
func (sq *SqlExecMultiple) Exec(params ...bool) error {
	cnn, err := Connection(sq.database)
	if err != nil {
		return err
	}

	ctx := context.Background()
	tx, err := cnn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error sql tx: %s ", err.Error())
	}

	err = cnn.PingContext(ctx)
	if err != nil {
		return errors.New(fmt.Sprint("Error Sql PING: ", err))
	}
	cross := false
	if len(params) == 1 {
		cross = params[0]
	}
	defer cnn.Close()

	for _, t := range sq.transaction {
		for _, item := range t.query {
			sqlPre := item["sqlPreparate"].(string)
			if cross {
				if t.action == "UPDATE" {
					sqlPre = Query_Cross_Update(sqlPre)
				}
			}
			valuesExec := item["valuesExec"].([]interface{})
			_, err := tx.Exec(sqlPre, valuesExec...)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error sql %s: %s", t.action, err.Error())
			}
		}
	}

	//Commit para confirmar la transacción
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error sql ping: %s", err.Error())
	}

	return nil
}

func (t *Transaction) Insert() error {
	sqlExec, data_insert, err := _insert(t.schema.GetTableName(), t.ob, t.schema.GetSchemaInsert())
	if err != nil {
		return err
	}
	t.query = sqlExec
	t.data = data_insert
	t.action = "INSERT"
	return nil
}

func (t *Transaction) Update() error {
	sqlExec, data_update, err := _update(t.schema.GetTableName(), t.ob, t.schema.GetSchemaUpdate())
	if err != nil {
		return err
	}
	t.query = sqlExec
	t.data = data_update
	t.action = "UPDATE"
	return nil
}

func (t *Transaction) Delete() error {
	sqlExec, data_delete, err := _delete(t.schema.GetTableName(), t.ob, t.schema.GetSchemaDelete())
	if err != nil {
		return err
	}
	t.query = sqlExec
	t.data = data_delete
	t.action = "DELETE"
	return nil
}

func (t *Transaction) GetData() []map[string]interface{} {
	return t.data
}

func (sq *SqlExecMultiple) ExecTransaction(t *Transaction) error {
	if sq.tx == nil {
		cnn, err := Connection(sq.database)
		if err != nil {
			return err
		}

		ctx := context.Background()
		sq.tx, err = cnn.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("error sql tx: %s ", err.Error())
		}

		err = cnn.PingContext(ctx)
		if err != nil {
			return errors.New(fmt.Sprint("Error Sql PING: ", err))
		}

		defer cnn.Close()
	}

	for _, item := range t.query {
		sqlPre := item["sqlPreparate"].(string)
		valuesExec := item["valuesExec"].([]interface{})
		_, err := sq.tx.Exec(sqlPre, valuesExec...)
		if err != nil {
			sq.tx.Rollback()
			return fmt.Errorf("error sql %s: %s", t.action, err.Error())
		}
	}

	return nil
}

func (sq *SqlExecMultiple) Commit() error {
	err := sq.tx.Commit()
	if err != nil {
		return fmt.Errorf("error sql ping: %s", err.Error())
	}
	return nil
}

func _insert(table string, data []map[string]interface{}, schema []Fields) ([]map[string]interface{}, []map[string]interface{}, error) {
	length := len(data)
	if length > 0 {
		var sqlExec = make([]map[string]interface{}, 0)
		var data_insert []map[string]interface{}

		for _, item := range data {
			preArray, err := _checkInsertSchema(schema, item)
			if err == nil {
				data_insert = append(data_insert, preArray)
				var column []string
				var values []string
				var i int
				var valuesExec []interface{}
				char := "$"
				for k, v := range preArray {
					i++
					column = append(column, k)
					values = append(values, fmt.Sprintf("%s%d", char, i))
					valuesExec = append(valuesExec, v)
				}

				sqlPreparate := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, strings.Join(column, ", "), strings.Join(values, ", "))
				sqlExec = append(sqlExec, map[string]interface{}{
					"sqlPreparate": sqlPreparate,
					"valuesExec":   valuesExec,
				})
			} else {
				return nil, nil, err
			}
		}
		return sqlExec, data_insert, nil
	} else {
		return nil, nil, errors.New("no existen datos para insertar")
	}
}

func _update(table string, data []map[string]interface{}, schema []Fields) ([]map[string]interface{}, []map[string]interface{}, error) {
	length := len(data)

	if length > 0 {
		var sqlExec = make([]map[string]interface{}, 0)
		var data_update []map[string]interface{}
		for _, item := range data {
			where := make(map[string]interface{})

			length_where := 0
			if item["where"] != nil {
				where = item["where"].(map[string]interface{})
				length_where = len(where)
				delete(item, "where")
			}
			preArray, err := _checkUpdate(schema, item)
			if err != nil {
				return nil, nil, err
			}
			preArray_where := make(map[string]interface{})
			if length_where > 0 {
				preArray, err := _checkWhere(schema, where)
				if err != nil {
					return nil, nil, err
				}
				preArray_where = preArray
			}

			data_update = append(data_update, preArray)
			var setters []string

			sqlWherePreparateUpdate := ""
			var i uint64
			var valuesExec []interface{}
			char := "$"
			for k, v := range preArray {
				i++
				setters = append(setters, fmt.Sprintf("%s= %s%d", k, char, i))
				valuesExec = append(valuesExec, v)
			}

			if length_where > 0 {
				length_newMapWhere := len(preArray_where)
				var wheres []string
				for k, v := range preArray_where {
					i++
					wheres = append(wheres, k+" = "+char+strconv.FormatUint(i, 10))
					valuesExec = append(valuesExec, v)
				}
				if length_newMapWhere > 0 {
					sqlWherePreparateUpdate = "WHERE " + strings.Join(wheres, " AND ")
				}
			}
			sqlPreparate := fmt.Sprintf("UPDATE %s SET %s %s", table, strings.Join(setters, ", "), sqlWherePreparateUpdate)
			sqlExec = append(sqlExec, map[string]interface{}{
				"sqlPreparate": sqlPreparate,
				"valuesExec":   valuesExec,
			})

		}
		return sqlExec, data_update, nil
	} else {
		return nil, nil, errors.New("no existen datos para actualizar")
	}
}

func _delete(table string, data []map[string]interface{}, schema []Fields) ([]map[string]interface{}, []map[string]interface{}, error) {
	length := len(data)

	if length > 0 {
		var sqlExec = make([]map[string]interface{}, 0)
		var data_delete []map[string]interface{}
		for _, item := range data {

			preArray, err := _checkWhere(schema, item)
			if err != nil {
				return nil, nil, err
			}

			data_delete = append(data_delete, preArray)
			var lineSqlExec = make(map[string]interface{}, 2)
			sqlWherePreparateDelete := ""
			var i int
			var p uint64
			length_newMap := len(preArray)
			var valuesExec []interface{}
			if length_newMap > 0 {
				sqlWherePreparateDelete += " WHERE "
			}
			char := "$"
			for k, v := range preArray {
				p++
				if i+1 < length_newMap {
					// sqlWherePreparateUpdate += fmt.Sprintf("%s = '%s' AND ", ke, va)
					sqlWherePreparateDelete += k + " = " + char + strconv.FormatUint(p, 10) + " AND "
				} else {
					//sqlWherePreparateUpdate += fmt.Sprintf("%s = '%s'", ke, va)
					sqlWherePreparateDelete += k + " = " + char + strconv.FormatUint(p, 10)
				}
				valuesExec = append(valuesExec, v)
				i++
			}

			sqlPreparate := fmt.Sprintf("DELETE FROM %s %s", table, sqlWherePreparateDelete)
			lineSqlExec["sqlPreparate"] = sqlPreparate
			lineSqlExec["valuesExec"] = valuesExec
			sqlExec = append(sqlExec, lineSqlExec)
		}
		return sqlExec, data_delete, nil
	} else {
		return nil, nil, errors.New("no existen datos para actualizar")
	}
}

func _checkInsertSchema(schema []Fields, tabla_map map[string]interface{}) (map[string]interface{}, error) {

	// var err_cont uint64 = 0
	var err_cont uint
	var error string

	data := make(map[string]interface{})

	for _, item := range schema {
		isNil := tabla_map[item.Name] == nil
		defaultIsNil := item.Default == nil
		if !isNil {
			value := tabla_map[item.Name]
			new_value, err := strconvDataType(string(item.Type), value)
			if err != nil {
				err_cont++
				error += fmt.Sprintf("%d.- El campo %s %s", err_cont, item.Description, err.Error())
			}
			var val interface{}
			switch item.Type {
			case "string":
				val, err = caseString(new_value.(string), item.ValidateType.(TypeStrings))
			case "float64":
				val, err = caseFloat(new_value.(float64), item.ValidateType.(TypeFloat64))
			case "uint64":
				val, err = caseUint(new_value.(uint64), item.ValidateType.(TypeUint64))
			case "int64":
				val, err = caseInt(new_value.(int64), item.ValidateType.(TypeInt64))
			default:
				val, err = nil, errors.New("tipo de dato no asignado")
			}
			if err == nil {
				data[item.Name] = val
			} else {
				err_cont++
				error += fmt.Sprintf("%d.- Se encontró fallas al validar el campo %s \n %s\n", err_cont, item.Description, err.Error())
			}
		} else {
			if !defaultIsNil {
				data[item.Name] = item.Default
			} else {
				if item.Required {
					err_cont++
					error += fmt.Sprintf("%d.- El campo %s es Requerido\n", err_cont, item.Description)
				}
			}
		}

	}
	if err_cont > 0 {
		return nil, errors.New(error)
	} else {
		return data, nil
	}
}

func _checkUpdate(schema []Fields, tabla_map map[string]interface{}) (map[string]interface{}, error) {
	var err_cont uint
	var error string
	data := make(map[string]interface{})
	for _, item := range schema {
		isNil := tabla_map[item.Name] == nil
		if !isNil {
			if item.Update {
				value := tabla_map[item.Name]
				new_value, err := strconvDataType(string(item.Type), value)
				if err != nil {
					err_cont++
					error += fmt.Sprintf("%d.- El campo %s %s", err_cont, item.Description, err.Error())
				}
				var val interface{}
				switch item.Type {
				case "string":
					if new_value.(string) == "" {
						if !item.Empty {
							err_cont++
							error += fmt.Sprintf("%d.- El campo %s no puede estar vació\n", err_cont, item.Description)
						}
					} else {
						val, err = caseString(new_value.(string), item.ValidateType.(TypeStrings))
					}
				case "float64":
					val, err = caseFloat(new_value.(float64), item.ValidateType.(TypeFloat64))
				case "uint64":
					val, err = caseUint(new_value.(uint64), item.ValidateType.(TypeUint64))
				case "int64":
					val, err = caseInt(new_value.(int64), item.ValidateType.(TypeInt64))
				default:
					val, err = nil, errors.New("tipo de dato no asignado")
				}
				if err == nil {
					data[item.Name] = val
				} else {
					err_cont++
					error += fmt.Sprintf("%d.- Se encontró fallas al validar el campo %s \n %s\n", err_cont, item.Description, err.Error())
				}
			} else {
				err_cont++
				error += fmt.Sprintf("%d.- El campo %s no puede ser modificado\n", err_cont, item.Description)
			}
		}
	}
	if err_cont > 0 {
		return nil, errors.New(error)
	} else {
		return data, nil
	}
}

func _checkWhere(schema []Fields, table_where map[string]interface{}) (map[string]interface{}, error) {
	var err_cont uint
	var error string
	data := make(map[string]interface{})
	for _, item := range schema {
		isNil := table_where[item.Name] == nil
		if !isNil {
			value := table_where[item.Name]
			if !item.Where && !item.PrimaryKey {
				err_cont++
				error += fmt.Sprintf("%d.- El campo %s no puede ser utilizado de esta forma\n", err_cont, item.Description)
			} else {
				if value.(string) == "" {
					err_cont++
					error += fmt.Sprintf("%d.- El campo %s esta vació verificar\n", err_cont, item.Description)
				} else {
					data[item.Name] = value
				}

			}
		} else {
			if item.PrimaryKey {
				err_cont++
				error += fmt.Sprintf("%d.- El campo %s es obligatorio\n", err_cont, item.Description)
			}
		}
	}
	if err_cont > 0 {
		return nil, errors.New(error)
	} else {
		return data, nil
	}
}

func caseString(value string, schema TypeStrings) (string, error) {
	value = strings.TrimSpace(value)
	if schema.Expr != nil {
		if !schema.Expr.MatchString(value) {
			return "", errors.New("- no cumple con las características\n")
		}
	}

	if schema.Date {
		err := CheckDate(value)
		if err != nil {
			return "", fmt.Errorf("- %s\n", err.Error())
		} else {
			return value, nil
		}
	}

	if schema.Encriptar {
		result, _ := bcrypt.GenerateFromPassword([]byte(value), 13)
		value = string(result)
		return value, nil
	}

	if schema.Cifrar {
		hash, _ := AesEncrypt_PHP([]byte(value), GetKey_PrivateCrypto())
		value = hash
		return value, nil
	}

	if schema.Min > 0 {
		if len(value) < schema.Min {
			return "", fmt.Errorf("- No Cumple los caracteres mínimos que debe tener (%v)\n", schema.Min)
		}
	}

	if schema.Max > 0 {
		if len(value) > schema.Max {
			return "", fmt.Errorf("- No Cumple los caracteres máximos que debe tener (%v)\n", schema.Max)
		}
	}

	if schema.UpperCase {
		value = strings.ToUpper(value)
	} else if schema.LowerCase {
		value = strings.ToLower(value)
	}
	return value, nil
}

func caseFloat(value float64, schema TypeFloat64) (float64, error) {
	error := ""
	err_cont := 0
	if schema.Menor != 0 {
		if value <= schema.Menor {
			err_cont++
			error += fmt.Sprintf("- No puede se menor a %f", schema.Menor)
		}
	}
	if schema.Mayor != 0 {
		if value >= schema.Mayor {
			err_cont++
			error += fmt.Sprintf("- No puede se mayor a %f", schema.Mayor)
		}
	}
	if !schema.Negativo {
		if value < float64(0) {
			err_cont++
			error += fmt.Sprintf("- no puede ser negativo")
		}
	}
	if schema.Porcentaje {
		value = value / float64(100)
	}
	if err_cont > 0 {
		return 0, errors.New(error)
	} else {
		return value, nil
	}
}

func caseInt(value int64, schema TypeInt64) (int64, error) {
	error := ""
	err_cont := 0
	if !schema.Negativo {
		if value < int64(0) {
			err_cont++
			error += fmt.Sprintf("- No puede ser negativo")
		}
	}
	if schema.Min != 0 {
		if value < schema.Min {
			err_cont++
			error += fmt.Sprintf("- No puede se menor a %d", schema.Min)
		}
	}
	if schema.Max != 0 {
		if value > schema.Max {
			err_cont++
			error += fmt.Sprintf("- No puede se mayor a %d", schema.Max)
		}
	}
	if err_cont > 0 {
		return int64(0), errors.New(error)
	} else {
		return value, nil
	}
}

func caseUint(value uint64, schema TypeUint64) (uint64, error) {
	if schema.Max > 0 {
		if value > schema.Max {
			return 0, errors.New("- no esta en el rango permitido")
		}
	}
	return value, nil
}

func convertStringToType(types string, value_undefined interface{}) (val interface{}, err error) {
	value := fmt.Sprintf("%v", value_undefined)
	switch types {
	case "uint64":
		val, err = strconv.ParseUint(value, 10, 64)
		return
	case "int64":
		val, err = strconv.ParseInt(value, 10, 64)
		return
	case "float64":
		val, err = strconv.ParseFloat(value, 64)
		return
	default:
		return nil, errors.New("No se puede convertir el tipo de dato")
	}
}

func strconvDataType(types string, values interface{}) (interface{}, error) {
	type_value := reflect.TypeOf(values).String()
	switch types {
	case "string":
		if types == type_value {
			return values, nil
		}
		return nil, errors.New("tipo de dato incorrecto")
	case "float64":
		if types == type_value {
			return values, nil
		}

		if type_value == "string" {
			new_value, err := convertStringToType("float64", values)
			if err != nil {
				return nil, err
			}
			return new_value, nil
		} else if type_value == "int64" {
			new_value := float64(values.(int64))
			return new_value, nil
		}

		return nil, errors.New("tipo de dato incorrecto")
	case "uint64":
		if types == type_value {
			return values, nil
		}

		if type_value == "float64" {
			new_value := uint64(values.(float64))
			return new_value, nil
		} else if type_value == "string" {
			new_value, err := convertStringToType("uint64", values)
			if err != nil {
				return nil, err
			}
			return new_value, nil
		} else if type_value == "int64" {
			new_value := uint64(values.(int64))
			return new_value, nil
		}

		return nil, errors.New("tipo de dato incorrecto")
	case "int64":
		if types == type_value {
			return values, nil
		}

		if type_value == "float64" {
			new_value := int64(values.(float64))
			return new_value, nil
		} else if type_value == "string" {
			new_value, err := convertStringToType("int64", values)
			if err != nil {
				return nil, err
			}
			return new_value, nil
		} else if type_value == "uint64" {
			new_value := int64(values.(uint64))
			return new_value, nil
		}

		return nil, errors.New("tipo de dato incorrecto")
	default:
		return nil, errors.New("No se puede convertir el tipo de dato")
	}
}

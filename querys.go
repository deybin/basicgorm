package basicgorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type QConfig struct {
	Cloud     bool
	Database  string
	Procedure bool
}

type Querys struct {
	Table   string   /** nombre de la tabla*/
	query   sintaxis /** guarda la estructura sql  de la consulta que se va contrayendo para luego ser formateada y mostrada en un string */
	rowSql  *sql.Rows
	colSql  []string
	db      *sql.DB
	tx      *sql.Tx
	ctx     context.Context
	err     error
	argsLen int           /** lleva en cuenta la cantidad de argumentos que tiene la consulta*/
	args    []interface{} /** almacena los argumentos que se le esta pasando ala consulta, el len de esta variable debe de ser igual al argsLen */
}

/** guarda la estructura de consulta sql, aparir de aquí se generar la consulta sql */
type sintaxis struct {
	Select        string
	Where         string
	Join          []string
	Top           string
	OrderBy       string
	GroupBy       string
	queryFull     string /** guarda la consulta sql directa en string */
	workQueryFull bool   /** establece si se va a utilizar una consulta directa mediante queryFull o mediante la estructura true:= se considerara queryFull false:= se considerara  estructura para formar la consulta sql*/
}

/** operaciones utilizadas con la sentencia WHERE*/
type OperatorWhere string

const (
	I           OperatorWhere = "="
	D           OperatorWhere = "<>"
	MY          OperatorWhere = ">"
	MYI         OperatorWhere = ">="
	MN          OperatorWhere = "<"
	MNI         OperatorWhere = "<="
	LIKE        OperatorWhere = "LIKE"
	IN          OperatorWhere = "IN"
	NOT_IN      OperatorWhere = "NOT IN"
	BETWEEN     OperatorWhere = "BETWEEN"
	NOT_BETWEEN OperatorWhere = "NOT BETWEEN"
)

/** Tipos de Join a utilizar en la consulta*/

type TypeJoin string

const (
	INNER TypeJoin = "INNER JOIN"
	LEFT  TypeJoin = "LEFT JOIN"
	RIGHT TypeJoin = "RIGHT JOIN"
	FULL  TypeJoin = "FULL OUTER JOIN"
)

func (q *Querys) Connect(config QConfig) *Querys {
	cloud := config.Cloud
	var errs error
	if cloud {
		q.db, errs = ConnectionCloud()
		if errs != nil {
			q.err = errs
			fmt.Println("Error SQL:", errs.Error())
			return q
		}
	} else {
		q.db, errs = Connection(config.Database)
		if errs != nil {
			q.err = errs
			fmt.Println("Error SQL:", errs.Error())
			return q
		}
	}
	q.ctx = context.Background()
	q.tx, errs = q.db.BeginTx(q.ctx, nil)

	if errs != nil {
		q.err = errs
		fmt.Println("Error SQL:", errs.Error())
	}

	return q
}

/*
*
SetQueryString establece una consulta SQL completa y sus argumentos en el struct Querys.

Esta función se utiliza para establecer una consulta SQL completa y sus argumentos en el struct Querys.
Esto permite ejecutar consultas SQL personalizadas que no se construyeron utilizando los métodos
de construcción de consultas normales del struct Querys.

Parámetros:
  - query: La consulta SQL completa que se va a establecer en el struct Querys.
  - arg: Los argumentos de la consulta SQL, que pueden ser un solo valor o un slice de valores.

Devuelve:
  - Una referencia al struct Querys actualizado con la consulta SQL y sus argumentos.
*/
func (q *Querys) SetQueryString(query string, arg interface{}) *Querys {
	q.query.workQueryFull = true
	q.query.queryFull = query
	if arg == nil {
		return q
	}
	if reflect.TypeOf(arg).String() == "[]interface {}" {
		q.args = append(q.args, arg.([]interface{})...)
	} else {
		q.args = append(q.args, arg)
	}

	return q
}

func (q *Querys) SetTable(table string) *Querys {
	q.Table = table
	return q
}

/*
*
Select establece la cláusula SELECT de la consulta SQL.
Puede especificar una lista de campos como argumentos.
Si no se proporcionan campos, se seleccionan todos (*).

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select().Exec("mi_database").All()
	result,err:=queryBuilder.Select("campo1,campo2").Exec("mi_database").All()
	result,err:=queryBuilder.Select("campo1", "campo2").Exec("mi_database").All()
	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - campos: Lista de nombres de campos a seleccionar. Si está vacío, se seleccionan todos los campos.

Devuelve:
  - Un puntero al struct Querys actualizado para permitir el encadenamiento de métodos.
*/
func (q *Querys) Select(campos ...string) *Querys {
	if len(campos) == 0 {
		q.query.Select = "SELECT * FROM " + q.Table
	} else {
		q.query.Select = "SELECT " + strings.Join(campos, ",") + " FROM " + q.Table
	}

	return q
}

/*
*
Join añade una cláusula JOIN a la consulta SQL.
El tipo de unión (INNER, LEFT, RIGHT, etc.), la tabla y la condición ON se especifican como parámetros.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "tabla_principal"}
	result,err:=queryBuilder.Select("tabla_principal.columna1, tabla_secundaria.columna2").
	    Join("INNER_JOIN", "tabla_secundaria", "tabla_principal.id = tabla_secundaria.id").
	    Where("tabla_principal.columna3", "=", valor).Exec("mi_database").All()
	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - tp: Tipo de unión (INNER, LEFT, RIGHT, etc.).
  - table: Nombre de la tabla a unir.
  - on: Condición ON para la unión.

Devuelve:
  - Un puntero al struct Querys actualizado para permitir el encadenamiento de métodos.
*/
func (q *Querys) Join(tp TypeJoin, table string, on string) *Querys {
	q.query.Join = append(q.query.Join, fmt.Sprintf(" %s  %s ON %s", tp, table, on))
	return q
}

/*
*
Where establece la cláusula WHERE de la consulta SQL con una condición y un operador.
La condición puede contener placeholders ($) para argumentos de la consulta.
El operador se utiliza para comparar valores en la condición.
El argumento es el valor que se comparará en la condición.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select("campo1, campo2").Where("campo3", "=", valor).Exec("mi_database").All()
	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - where: Condición para la cláusula WHERE.
  - op: Operador para comparar valores en la condición (por ejemplo, "=", "<>", ">", "<","<=", ">=", "LIKE", "IN", "NOT IN", "BETWEEN" "NOT BETWEEN").
  - arg: Valor que se compara en la condición.

Devuelve:
  - Un puntero al struct Querys actualizado para permitir el encadenamiento de métodos.
*/
func (q *Querys) Where(where string, op OperatorWhere, arg interface{}) *Querys {
	q.argsLen++
	argString, err := getSintaxisFilter(q, op, arg)
	if err != nil {
		fmt.Println(err)
		return q
	}
	q.query.Where = fmt.Sprintf(" WHERE %s %s %s", where, op, argString)
	return q
}

/*
*
And añade una cláusula AND adicional a la cláusula WHERE existente de la consulta SQL.
La condición, el operador y el argumento se especifican como parámetros.

Esta función se utiliza para agregar condiciones adicionales a la cláusula WHERE de la consulta SQL.
Si la cláusula WHERE aún no está especificada en la consulta, esta función no hace nada.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select("campo1, campo2").Where("campo3", "=", valor).And("campo4", ">", otroValor).Exec("mi_database").All()
	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - and: Condición adicional para agregar a la cláusula WHERE existente.
  - op: Operador para comparar valores en la condición (por ejemplo, "=", "<>", ">", "<","<=", ">=", "LIKE", "IN", "NOT IN", "BETWEEN" "NOT BETWEEN").
  - arg: Valor que se compara en la condición.

Devuelve:
  - Un puntero al struct Querys actualizado para permitir el encadenamiento de métodos.
*/
func (q *Querys) And(and string, op OperatorWhere, arg interface{}) *Querys {
	if q.query.Where == "" {
		return q
	}
	argString, err := getSintaxisFilter(q, op, arg)
	if err != nil {
		fmt.Println(err)
		return q
	}
	q.query.Where += fmt.Sprintf(" AND %s %s %s", and, op, argString)
	return q
}

/*
*
Or añade una cláusula OR adicional a la cláusula WHERE existente de la consulta SQL.
La condición, el operador y el argumento se especifican como parámetros.

Esta función se utiliza para agregar condiciones adicionales a la cláusula WHERE de la consulta SQL.
Si la cláusula WHERE aún no está especificada en la consulta, esta función no hace nada.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select("campo1, campo2").Where("campo3", "=", valor).Or("campo4", ">", otroValor).Exec("mi_database").All()
	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - and: Condición adicional para agregar a la cláusula WHERE existente.
  - op: Operador para comparar valores en la condición (por ejemplo, "=", "<>", ">", "<","<=", ">=", "LIKE", "IN", "NOT IN", "BETWEEN" "NOT BETWEEN").
  - arg: Valor que se compara en la condición.

Devuelve:
  - Un puntero al struct Querys actualizado para permitir el encadenamiento de métodos.
*/
func (q *Querys) Or(or string, op OperatorWhere, arg interface{}) *Querys {
	if q.query.Where == "" {
		return q
	}
	argString, err := getSintaxisFilter(q, op, arg)
	if err != nil {
		fmt.Println(err)
		return q
	}
	q.query.Where += fmt.Sprintf(" OR %s %s %s", or, op, argString)
	return q
}

/*
*
OrderBy establece la cláusula ORDER BY de la consulta SQL.
Puede especificar una lista de campos como argumentos.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select("campo1, campo2").Where("campo3","=", valor).OrderBy("campo4 DESC", "campo5 ASC").Exec("mi_database").All()
	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - campos: Lista de nombres de campos a utilizar en la cláusula ORDER BY.

Devuelve:
  - Un puntero al struct Querys actualizado para permitir el encadenamiento de métodos.
*/
func (q *Querys) OrderBy(campos ...string) *Querys {
	q.query.OrderBy = " ORDER BY " + strings.Join(campos, ",")
	return q
}

/*
*
GroupBy establece la cláusula GROUP BY de la consulta SQL.
Puede especificar una lista de campos como argumentos.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select("campo1", "campo2").GroupBy("campo4", "campo5").Exec("mi_database").All()
	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - group: Lista de nombres de campos a utilizar en la cláusula GROUP BY.

Devuelve:
  - Un puntero al struct Querys actualizado para permitir el encadenamiento de métodos.
*/
func (q *Querys) GroupBy(group ...string) *Querys {
	if len(group) <= 0 {
		return q
	}
	q.query.GroupBy = fmt.Sprintf(" GROUP BY %s", strings.Join(group, ","))
	return q
}

/*
*
Top establece la cláusula LIMIT de la consulta SQL para seleccionar un número específico de filas.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select("campo1", "campo2").Top(10).Exec("mi_database").All()
	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - top: Número de filas a seleccionar.

Devuelve:
  - Un puntero al struct Querys actualizado para permitir el encadenamiento de métodos.
*/
func (q *Querys) Top(top int) *Querys {
	q.query.Top = fmt.Sprintf(" LIMIT %d", top)
	return q
}

/*
*
Limit establece la cláusula LIMIT de la consulta SQL para seleccionar un número específico de filas,
con opción para especificar un offset.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select().Where("campo3","=", valor).Limit(10).Exec("mi_database").All() // Limita la consulta a 10 filas
	result,err:=queryBuilder.Select("campo1", "campo2").Where("campo3","=", valor).Limit(10,3).Exec("mi_database").All() // Limita la consulta a 10 filas omitiendo 3 filas en el conjunto de resultados de una consulta
	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - limit: Lista de uno o dos enteros. Si se proporciona un solo entero, establece el límite de filas.
    Si se proporcionan dos enteros, el primero especifica el límite de filas y el segundo especifica el offset.

Devuelve:
  - Un puntero al struct Querys actualizado para permitir el encadenamiento de métodos.
*/
func (q *Querys) Limit(limit ...int) *Querys {
	if len(limit) == 2 {
		q.query.Top = fmt.Sprintf(" LIMIT %d OFFSET %d", limit[0], limit[1])
	} else if len(limit) == 1 {
		q.query.Top = fmt.Sprintf(" LIMIT %d", limit[0])
	} else {
		q.query.Top = " LIMIT 1"
	}

	return q
}

/*
*
Exec ejecuta la consulta SQL en la base de datos y almacena los resultados en la estructura Querys.

Esta función se utiliza para ejecutar consultas SQL en la base de datos y almacenar los resultados
en la estructura Querys para su posterior procesamiento.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select("campo1", "campo2").Where("campo3", "=", valor).Exec("mi_database").All() // Ejecuta la consulta SQL utilizando la configuración proporcionada.

Parámetros:
  - config: Configuración para la conexión a la base de datos, incluyendo detalles como la nube, el nombre de la base de datos, etc.

Devuelve:
  - Un puntero al struct Querys actualizado con los resultados de la consulta ejecutada.
*/
func (q *Querys) Exec(config QConfig) *Querys {
	cloud := config.Cloud
	var db *sql.DB
	if cloud {
		var errs error
		db, errs = ConnectionCloud()
		if errs != nil {
			q.err = errs
			fmt.Println("Error SQL:", errs.Error())
			return q
		}

	} else {
		var errs error
		db, errs = Connection(config.Database)
		if errs != nil {
			q.err = errs
			fmt.Println("Error SQL:", errs.Error())
			return q
		}

	}

	ctx := context.Background()
	err := db.PingContext(ctx)
	defer db.Close()
	if err != nil {
		q.err = err
		fmt.Println("Error SQL ping:", err.Error())
		return q
	}
	queryString := q.GetQuery()
	// fmt.Println("query:", queryString)
	if !config.Procedure {
		rows, err := db.QueryContext(ctx, queryString, q.args...)
		if err != nil {
			q.err = err
			fmt.Println("Error SQL exec:", err.Error())
			return q
		}

		cols, _ := rows.Columns()

		q.rowSql = rows
		q.colSql = cols

		return q
	} else {
		_, err = db.ExecContext(ctx, queryString, q.args...)
		if err != nil {
			q.err = err
			fmt.Println("Error SQL exec:", err.Error())
		}
		return q
	}

}

func (q *Querys) ExecTx() *Querys {
	if q.err != nil {
		return q
	}

	err := q.db.PingContext(q.ctx)
	defer q.db.Close()
	if err != nil {
		q.err = err
		fmt.Println("Error SQL ping:", err.Error())
		return q
	}

	ctx := context.Background()
	queryString := q.GetQuery()
	rows, err := q.tx.QueryContext(ctx, queryString, q.args...)
	if err != nil {
		fmt.Println("Error SQL exec tx:", err.Error())
		q.err = err
		return q
	}
	cols, _ := rows.Columns()

	q.rowSql = rows
	q.colSql = cols
	return q
}

func (q *Querys) ResetQuery() {
	q.query = sintaxis{}
}

func (q *Querys) SetNewArgs(arg ...interface{}) error {
	if len(arg) != q.argsLen {
		return fmt.Errorf("parámetros enviados (%d), requeridos (%d)", len(arg), q.argsLen)
	}
	q.args = arg
	return nil
}

func (q *Querys) Close() {
	q.tx.Commit()
	q.rowSql.Close()
	q.db.Close()
}

/*
*
One recupera un solo resultado de la consulta SQL y lo devuelve como un mapa.

Esta función se utiliza para obtener un solo resultado de una consulta SQL ejecutada previamente
utilizando el método Exec, y devuelve el resultado como un mapa donde las claves son los nombres de las columnas
y los valores son los valores de las columnas correspondientes para la primera fila de resultados.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select("campo1, campo2").Where("campo3", "=", valor).Exec("mi_database").One() // Obtiene un solo resultado de la consulta ejecutada.

Devuelve:
  - Un mapa donde las claves son los nombres de las columnas y los valores son los valores de las columnas correspondientes
    para la primera fila de resultados.
  - Un error, si ocurre alguno durante el proceso de recuperación del resultado.
*/
func (q *Querys) One() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if q.err != nil {
		return m, q.err
	}
	for q.rowSql.Next() {
		columns := make([]interface{}, len(q.colSql))
		columnPointers := make([]interface{}, len(q.colSql))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := q.rowSql.Scan(columnPointers...); err != nil {
			fmt.Println(err)
			return map[string]interface{}{}, err
		}
		for i, colName := range q.colSql {
			val := columnPointers[i].(*interface{})
			l := *val
			if l != nil {
				m[colName] = l
			} else {
				m[colName] = l
			}
		}
		break
	}

	defer q.rowSql.Close()
	return m, nil
}

/*
*
Text recupera el valor de una columna específica de la primera fila de resultados de la consulta SQL.

Esta función se utiliza para obtener el valor de una columna específica de la primera fila de resultados
de una consulta SQL ejecutada previamente utilizando el método Exec.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	queryBuilder.Select("campo1, campo2").Where("campo3", "=", valor).Exec("mi_database") // Ejecuta la consulta SQL utilizando la configuración proporcionada.
	valor, err := queryBuilder.Text("nombreColumna") // Obtiene el valor de la columna "nombreColumna".

Parámetros:
  - columna: Nombre de la columna de la cual se desea obtener el valor.

Devuelve:
  - El valor de la columna especificada.
  - Un error, si ocurre alguno durante el proceso de recuperación del valor.
*/
func (q *Querys) Text(columna string) (interface{}, error) {
	m := make(map[string]interface{})
	if q.err != nil {
		return nil, q.err
	}
	for q.rowSql.Next() {
		columns := make([]interface{}, len(q.colSql))
		columnPointers := make([]interface{}, len(q.colSql))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := q.rowSql.Scan(columnPointers...); err != nil {
			fmt.Println(err)
			return nil, err
		}

		for i, colName := range q.colSql {
			val := columnPointers[i].(*interface{})

			l := *val
			if l != nil {
				m[colName] = l
			} else {
				m[colName] = l
			}

		}

		break
	}
	defer q.rowSql.Close()
	return m[columna], nil
}

/*
*
All recupera todos los resultados de la consulta SQL y los devuelve como una lista de mapas.

Esta función se utiliza para obtener todos los resultados de una consulta SQL ejecutada previamente
utilizando el método Exec, y devuelve los resultados como una lista de mapas donde cada mapa
representa una fila de resultados, con los nombres de las columnas como claves y los valores de las columnas como valores.

Ejemplo de uso:

	queryBuilder := &Querys{Table: "mi_tabla"}
	result,err:=queryBuilder.Select("campo1, campo2").Where("campo3", "=", valor).Exec("mi_database").All() // Obtiene todos los resultados de la consulta ejecutada.

Devuelve:
  - Una lista de mapas, donde cada mapa representa una fila de resultados con los nombres de las columnas como claves y los valores de las columnas como valores.
  - Un error, si ocurre alguno durante el proceso de recuperación de resultados.
*/
func (q *Querys) All() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	if q.err != nil {
		return result, q.err
	}

	for q.rowSql.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.

		columns := make([]interface{}, len(q.colSql))
		columnPointers := make([]interface{}, len(q.colSql))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := q.rowSql.Scan(columnPointers...); err != nil {
			fmt.Println(err)
			return []map[string]interface{}{}, err
		}

		//Crea nuestro mapa y recupera el valor de cada columna del segmento de punteros, almacenándolo en el mapa con el nombre de la columna como clave.
		m := make(map[string]interface{})
		for i, colName := range q.colSql {
			val := columnPointers[i].(*interface{})
			l := *val
			if l != nil {

				m[colName] = l

			} else {
				m[colName] = l
			}
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		result = append(result, m)
	}
	defer q.rowSql.Close()
	return result, nil
}

/*
*
GetQuery devuelve la consulta SQL completa construida utilizando los métodos del struct Querys.

Esta función se utiliza para obtener la consulta SQL completa que se ha construido utilizando
los métodos del struct Querys, incluyendo la selección de columnas, cláusulas JOIN, condiciones WHERE,
agrupamiento GROUP BY, ordenamiento ORDER BY, y limitación de resultados con TOP o LIMIT.

Devuelve:
  - Una cadena que representa la consulta SQL completa construida.
*/
func (q *Querys) GetQuery() string {
	var queryString string
	if !q.query.workQueryFull {
		queryString = q.query.Select
		/** aplicando los join  inner join, left join y right join*/
		if len(q.query.Join) > 0 {
			for _, v := range q.query.Join {
				queryString += v
			}
		}
		/** aplicando Where : where ,and ,or ,in, between ,not in ,not between*/
		queryString += q.query.Where

		/** aplicando Group by*/
		queryString += q.query.GroupBy

		/** aplicando order by  */
		queryString += q.query.OrderBy
		/** aplicando Top y LImit  */
		queryString += q.query.Top
	} else {
		queryString = q.query.queryFull
	}

	return queryString
}

func (q *Querys) GetErrors() error {

	return q.err
}

/*
*
getSintaxisFilter genera la sintaxis adecuada para los filtros de las consultas SQL y maneja los argumentos correspondientes.

Esta función se utiliza para generar la sintaxis adecuada para los filtros de las consultas SQL, como IN, NOT IN, BETWEEN, y NOT BETWEEN,
y maneja los argumentos correspondientes, agregándolos a la lista de argumentos de la consulta SQL.

Parámetros:
  - q: Una instancia del struct Querys que contiene información sobre los argumentos de la consulta.
  - op: El operador de filtro (OperatorWhere) que se va a aplicar en la consulta.
  - arg: El valor o valores del filtro.

Devuelve:
  - Una cadena que representa la sintaxis adecuada para el filtro en la consulta SQL.
  - Un error, si ocurre alguno durante el proceso de generación de la sintaxis o manejo de los argumentos.
*/
func getSintaxisFilter(q *Querys, op OperatorWhere, arg interface{}) (string, error) {
	var argString string

	if op == IN || op == NOT_IN {
		if reflect.TypeOf(arg).String() != "[]interface {}" {
			return "", errors.New("tipo de dato incorrecto para filtrado IN")
		}
		if len(arg.([]interface{})) <= 0 {
			return "", errors.New("valor vació para filtrado IN")
		}
		arrayArgsSql := make([]string, 0)
		for _, v := range arg.([]interface{}) {
			arrayArgsSql = append(arrayArgsSql, fmt.Sprintf("$%d", q.argsLen))
			q.args = append(q.args, v)
			q.argsLen++
		}
		argString = fmt.Sprintf("(%s)", strings.Join(arrayArgsSql, ","))
	} else if op == BETWEEN || op == NOT_BETWEEN {
		if reflect.TypeOf(arg).String() != "[]interface {}" {
			return "", errors.New("tipo de dato incorrecto para filtrado BETWEEN")
		}
		if len(arg.([]interface{})) < 2 {
			return "", errors.New("valor vació o bien valores incompletos para filtrado BETWEEN")
		}
		argString = fmt.Sprintf("$%d AND ", q.argsLen)
		q.args = append(q.args, arg.([]interface{})[0])
		q.argsLen++
		argString += fmt.Sprintf("$%d", q.argsLen)
		q.args = append(q.args, arg.([]interface{})[1])
		q.argsLen++
	} else {
		argString = fmt.Sprintf("$%d", q.argsLen)
		q.args = append(q.args, arg)
		q.argsLen++
	}

	return argString, nil
}

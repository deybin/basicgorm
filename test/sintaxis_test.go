package test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/deybin/basicgorm"
)

func TestSintaxis(t *testing.T) {

	// Abrir conexión a la base de datos
	db, err := basicgorm.Connection("new_capital")
	if err != nil {
		fmt.Println(err)
	}
	// Consulta SQL con placeholders (?) para parámetros
	query := "SELECT l_clie FROM requ_clientes WHERE n_docu = $1"

	ctx := context.Background()
	err = db.PingContext(ctx)
	defer db.Close()
	if err != nil {
		fmt.Println("Error SQL:", err.Error())
	}

	// // Crear una sentencia preparada
	// stmt, err := db.PrepareContext(ctx, query)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer stmt.Close()

	// ID proporcionado por el usuario (debería ser validado y sanitizado)
	userID := "' or '1'='1"

	// Ejecutar la consulta utilizando la sentencia preparada
	rows, err := db.QueryContext(ctx, query, userID)
	// rows, err := stmt.Query(userID)
	if err != nil {
		fmt.Println("34:", err)
	}
	defer rows.Close()

	// Iterar sobre los resultados
	for rows.Next() {
		var l_clie string
		if err := rows.Scan(&l_clie); err != nil {
			fmt.Println("err", err)
		}
		fmt.Printf("Nombre: %s", l_clie)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, rows)
	}
}

/** Test aplicando SELECT con WHERE*/
func TestQueryWhere(t *testing.T) {
	Query := basicgorm.Querys{
		Table: "requ_clientes",
	}
	r, err := Query.Select().Where(" to_date(f_naci,'DD/MM/YYYY')", basicgorm.I, "1994-04-04").Exec(basicgorm.QConfig{Database: "new_capital"}).All()
	fmt.Println("test query:", Query.GetQuery())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

/** Test aplicando SELECT con WHERE y AND*/
func TestQueryWhereWhitAnd(t *testing.T) {
	Query := basicgorm.Querys{
		Table: "requ_clientes",
	}
	r, err := Query.Select().Where(" to_date(f_naci,'DD/MM/YYYY')", basicgorm.MYI, "1994-04-04").And("c_ubig", basicgorm.I, "120107").Exec(basicgorm.QConfig{Database: "new_capital"}).All()
	fmt.Println("test query:", Query.GetQuery())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

func TestQueryWhereIN(t *testing.T) {
	Query := basicgorm.Querys{
		Table: "requ_clientes",
	}
	r, err := Query.Select().Where(" to_date(f_naci,'DD/MM/YYYY')", basicgorm.IN, []interface{}{"1994-04-04", "1994-04-04", "1994-04-04"}).Exec(basicgorm.QConfig{Database: "new_capital"}).All()
	fmt.Println("test query:", Query.GetQuery())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

func TestQueryWhereBETWEEN(t *testing.T) {
	Query := basicgorm.Querys{
		Table: "requ_clientes",
	}
	r, err := Query.Select().Where(" to_date(f_naci,'DD/MM/YYYY')", basicgorm.BETWEEN, []interface{}{"1994-04-04", "1994-05-04"}).Exec(basicgorm.QConfig{Database: "new_capital"}).All()
	fmt.Println("test query:", Query.GetQuery())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

func TestQueryWhereWithAndIN(t *testing.T) {
	Query := basicgorm.Querys{
		Table: "requ_clientes",
	}
	r, err := Query.Select("n_docu").Where("c_ubig", basicgorm.I, "120119").And("n_docu", basicgorm.IN, []interface{}{"47727049", "20060977", "43198110"}).Exec(basicgorm.QConfig{Database: "new_capital"}).OrderBy("n_docu").All()
	fmt.Println("test query:", Query.GetQuery())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

func TestQueryWhereWithAndINOrderBy(t *testing.T) {
	Query := basicgorm.Querys{
		Table: "requ_clientes",
	}
	r, err := Query.Select("n_docu").Where("c_ubig", basicgorm.I, "120119").And("n_docu", basicgorm.IN, []interface{}{"47727049", "20060977", "43198110"}).OrderBy("n_docu desc").Top(2).Exec(basicgorm.QConfig{Database: "new_capital"}).All()
	fmt.Println("test query:", Query.GetQuery())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

func TestQueryInnerJoin(t *testing.T) {
	Query := basicgorm.Querys{
		Table: "requ_clientes as a",
	}
	r, err := Query.Select("a.n_docu,b.l_nomb").Join(basicgorm.INNER, "fina_clientes as b", "a.n_docu=b.n_docu").Where("a.c_ubig", basicgorm.I, "120119").And("a.n_docu", basicgorm.IN, []interface{}{"47727049", "20060977", "43198110"}).OrderBy("a.n_docu desc").Top(2).Exec(basicgorm.QConfig{Database: "new_capital"}).All()
	fmt.Println("test query:", Query.GetQuery())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

func TestQueryStringFull(t *testing.T) {
	Query := basicgorm.Querys{
		Table: "requ_clientes as a",
	}
	r, err := Query.SetQueryString("SELECT a.n_docu,b.l_nomb FROM requ_clientes as a INNER JOIN  fina_clientes as b ON a.n_docu=b.n_docu WHERE a.c_ubig = $1 AND a.n_docu IN ($2,$3,$4) ORDER BY a.n_docu desc LIMIT 2", "120119", "47727049", "20060977", "43198110").Exec(basicgorm.QConfig{Database: "new_capital"}).All()
	fmt.Println("test query:", Query.GetQuery())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

func TestQueryStringFull_plpgsql(t *testing.T) {
	Query := basicgorm.Querys{}
	sql_string := fmt.Sprintf(`
		DO $$
			DECLARE
				_s_capi float8;
				_s_inte float8;
				_s_mora float8;
				_s_desc float8;
				_s_tota float8;
				_n_days int;
				_k_stad int;
			BEGIN
				SELECT INTO _s_capi,_s_inte,_s_mora,_s_desc,_s_tota,_n_days
				sum(s_acapi), sum(s_ainte), sum(s_mora), sum(s_desc), sum(s_amor), sum(n_dias) 
				from  Fina_CreditsDetalle where id_cred='%s';

				SELECT INTO _k_stad count(k_stad) from  Fina_CreditsDetalle where id_cred='%s' and k_stad=0;
				if _k_stad>0 then _k_stad=0; else _k_stad=1; end if;
				
				Update Fina_Credits set s_acapi=_s_capi, s_ainte=_s_inte, s_mora=_s_mora, s_desc=_s_desc, s_amor=_s_tota, n_days=_n_days, s_prom=_n_days::float8/%d, k_stad=_k_stad
				where id_cred='%s';
				
			END;
			$$ LANGUAGE plpgsql;
		`, "0010000000001", "0010000000001", 30, "0010000000001")

	procedure := Query.SetQueryString(sql_string, nil).Exec(basicgorm.QConfig{Database: "new_capital", Procedure: true})
	// fmt.Println("test query:", Query.GetQuery())
	err := procedure.GetErrors()
	if err != nil {
		fmt.Println("eerr:", err)
	}

	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, "")
	}
}

func TestSintaxisContexto(t *testing.T) {
	query := new(basicgorm.Querys).SetTable("requ_clientes").Connect(basicgorm.QConfig{Database: "documentos"})
	query.Select("docs")
	rs := []map[string]interface{}{}
	arr := []string{"00001", "00002", "00003", "00004", "00005", "00006", "00007"}
	for i := 0; i < len(arr); i++ {
		query.ResetQuery()
		r, err := query.Select("docs").Where("docs", basicgorm.I, arr[i]).ExecTx().All()
		fmt.Println("test query:", query.GetQuery())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(r)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("no se cierra aun la conexión")
	time.Sleep(10 * time.Second)
	query.Close()
	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, rs)
	}
}

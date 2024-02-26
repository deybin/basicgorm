package basicgorm

import (
	"regexp"
)

// Schema es una interfaz que define métodos para obtener información sobre el esquema de una tabla en una base de datos.
type Schema interface {
	GetTableName() string
	GetSchemaInsert() []Fields
	GetSchemaUpdate() []Fields
	GetSchemaDelete() []Fields
}

type (
	// DataType basicGORM data type
	DataType string
)

const (
	Bool   DataType = "bool"
	Int    DataType = "int64"
	Uint   DataType = "uint64"
	Float  DataType = "float64"
	String DataType = "string"
	Time   DataType = "time"
	Bytes  DataType = "bytes"
)

/*
Las etiquetas de structure o también llamado etiquetas de campo estos metadatos serán los siguientes según el tipo de dato
*/
type Fields struct {
	Name         string      //Nombre del campo
	Description  string      //Descripción del campo
	Type         DataType    //A bajo nivel es un string donde se especifica de que tipo sera el campo
	Required     bool        //Si el valor para inserción de este campo es requerido o no
	PrimaryKey   bool        //Si el campo es primary key entonces es obligatorio este campo para insert,update y delete
	Where        bool        //El campo puede ser utilizado para filtrar al utilizar el update y delete
	Update       bool        //El campo puede ser modificado
	Default      interface{} //Valor por defecto que se tomara si no se le valor al campo, el tipo del valor debe de ser igual al Type del campo
	Empty        bool        //El campo aceptara valor vació si se realiza la actualización
	ValidateType interface{} //Los datos serán validados mas a fondo mediante esta opción para eso se le debe de asignar los siguientes typo de struct: TypeStrings, TypeFloat64, TypeUint64 yTypeInt64
}

type TypeStrings struct {
	LowerCase bool           //Convierte en minúscula el valor del campo
	UpperCase bool           //Convierte en mayúscula el valor del campo
	Encriptar bool           //Crea un hash del valor del campo
	Cifrar    bool           //Cifra el valor del campo
	Date      bool           //Verifica que el valor del campo sea una fecha valida con formato dd/mm/yyyy
	Min       int            //Cuantos caracteres como mínimo debe de tener el valor del campo
	Max       int            //Cuantos caracteres como máximo debe de tener el valor del campo
	Expr      *regexp.Regexp //Expresión regular que debe cumplir el valor que almacenara el campo

}

type TypeFloat64 struct {
	Porcentaje bool    //Convierte el valor del campo en porcentaje
	Negativo   bool    //El campo aceptara valores negativos
	Menor      float64 //Los valores que aceptaran tienen que ser menores o igual que este metadato
	Mayor      float64 //Los valores que aceptaran tienen que ser mayores o igual  que este metadato
}

type TypeUint64 struct {
	Max uint64 //Hasta que valor aceptara que se almacene en este campo
}

type TypeInt64 struct {
	Max      int64 // Hasta que valor aceptara que se almacene en este campo
	Min      int64 // Valor como mínimo  aceptara que se almacene en este campo
	Negativo bool  // Rl campo aceptara valores negativos
}

type Regex interface {
	Letras(start int8, end int16) *regexp.Regexp
	Float() *regexp.Regexp
}

func Null() *regexp.Regexp {
	return regexp.MustCompile(``)
}

func Number() *regexp.Regexp {
	return regexp.MustCompile(`[0-9]{0,}$`)
}

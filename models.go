package basicgorm

import (
	"regexp"
)

/*
*
Las etiquetas de structure o también llamado etiquetas de campo estos metadatos serán los siguientes según el tipo de dato:

	Type data All
	Name: string => nombre del campo
	Description: string => descripción del campo
	Type: DataType => a bajo nivel es un string donde se especifica de que tipo sera el campo

Required: bool => si el valor para inserción de este campo es requerido o no
PrimaryKey: bool => si el campo es primary key entonces es obligatorio este campo para insert,update y delete
Where: bool => si el campo puede ser utilizado para filtrar al utilizar el update y delete
Update: bool => si el campo puede ser modificado.
Default: interface{} => calor por defecto que se tomara si no se le valor al campo, el tipo del valor debe de ser igual al Type del campo.
Empty: bool => el campo aceptara valor vació si se realiza la actualización.
ValidateType: interface{} => los datos seran validados mas a fondo mediante esta opción para eso se le debe de asignar los sguientes typo de struct:

  - TypeStrings

  - TypeFloat64

  - TypeUint64

  - TypeInt64

    Type data string
    LowerCase: bool => convierte en minúscula el valor del campo
    UpperCase: bool => convierte en mayúscula el valor del campo
    Encriptar: bool => crea un hash del valor del campo
    Cifrar: bool => cifra el valor del campo
    Date: bool => verifica que el valor del campo sea una fecha valida con formato dd/mm/yyyy
    Min: int => cuantos caracteres como mínimo debe de tener el valor del campo
    Max: int => cuantos caracteres como máximo debe de tener el valor del campo
    Expr: *regexp.Regexp => expresión regular que debe cumplir el valor que almacenara el campo

    Type data float64
    Porcentaje: bool => convierte el valor del campo en porcentaje
    Negativo: bool => el campo aceptara valores negativos
    Menor: float64 => los valores que aceptaran tienen que ser menores o igual que este metadato
    Mayor: float64 => los valores que aceptaran tienen que ser mayores o igual  que este metadato

    Type data uint64
    Max: uint64 => hasta que valor aceptara que se almacene en este campo

    Type data int64
    Max: uint64 => hasta que valor aceptara que se almacene en este campo
    Min: uint64 => valor como mínimo  aceptara que se almacene en este campo
    Negativo: bool => el campo aceptara valores negativos
*/
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

type Fields struct {
	Name         string
	Description  string
	Type         DataType
	Required     bool
	PrimaryKey   bool
	Where        bool
	Update       bool
	Default      interface{}
	Empty        bool
	ValidateType interface{}
}

type TypeStrings struct {
	LowerCase bool
	UpperCase bool
	Encriptar bool
	Cifrar    bool
	Date      bool
	Min       int
	Max       int
	Expr      *regexp.Regexp
}

type TypeFloat64 struct {
	Porcentaje bool
	Negativo   bool
	Menor      float64
	Mayor      float64
}

type TypeUint64 struct {
	Max uint64
}

type TypeInt64 struct {
	Max      int64
	Min      int64
	Negativo bool
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

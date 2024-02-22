package basicgorm

import (
	"regexp"
)

/*
*
Las etiquetas de structure o también llamado etiquetas de campo estos metadatos serán los siguientes según el tipo de dato:

	Type data All
	json: string => nombre de la key que se pondrá al generar el json
	description: string => descripción del campo

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
	// DataType GORM data type
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

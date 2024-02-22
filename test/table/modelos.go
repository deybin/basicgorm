package table

import (
	"github.com/deybin/basicgorm"
)

func GetSucursal() ([]basicgorm.Fields, string) {
	tableName := "requ_" + "sucursal"
	var sucursal []basicgorm.Fields
	sucursal = append(sucursal, basicgorm.Fields{ //c_sucu
		Name:        "c_sucu",
		Description: "c_sucu",
		Required:    true,
		PrimaryKey:  true,
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min: 3,
			Max: 3,
		},
	})
	sucursal = append(sucursal, basicgorm.Fields{ //l_sucu
		Name:        "l_sucu",
		Description: "l_sucu",
		Required:    true,
		Update:      true,
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min:       3,
			Max:       100,
			LowerCase: true,
		},
	})
	return sucursal, tableName

}

func GetAlmacen() ([]basicgorm.Fields, string) {
	tableName := "requ_" + "almacen"
	var sucursal []basicgorm.Fields
	sucursal = append(sucursal, basicgorm.Fields{ //c_sucu
		Name:        "c_sucu",
		Description: "c_sucu",
		Required:    true,
		PrimaryKey:  true,
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min: 3,
			Max: 3,
		},
	})
	sucursal = append(sucursal, basicgorm.Fields{ //c_sucu
		Name:        "c_alma",
		Description: "c_alma",
		Required:    true,
		PrimaryKey:  true,
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min: 3,
			Max: 3,
		},
	})
	sucursal = append(sucursal, basicgorm.Fields{ //l_sucu
		Name:        "l_alma",
		Description: "l_alma",
		Required:    true,
		Update:      true,
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min:       3,
			Max:       100,
			LowerCase: true,
		},
	})
	return sucursal, tableName

}

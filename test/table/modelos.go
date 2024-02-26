package table

import "github.com/deybin/basicgorm"

type Sucursal struct {
	name string
}

func (s *Sucursal) New() *Sucursal {
	s.name = "requ_" + "sucursal"
	return s
}

func (s *Sucursal) getSchema() []basicgorm.Fields {
	var schema []basicgorm.Fields

	schema = append(schema, basicgorm.Fields{ //c_sucu
		Name:        "c_sucu",
		Description: "c_sucu",
		Required:    true,
		Where:       true,
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min: 3,
			Max: 3,
		},
	})

	schema = append(schema, basicgorm.Fields{ //l_sucu
		Name:        "l_sucu",
		Description: "l_sucu",
		Required:    true,
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min:       10,
			Max:       100,
			LowerCase: true,
		},
	})

	schema = append(schema, basicgorm.Fields{ //l_dire
		Name:        "l_dire",
		Description: "l_dire",
		Required:    true,
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min:       5,
			Max:       200,
			LowerCase: true,
		},
	})

	schema = append(schema, basicgorm.Fields{ //c_ubig
		Name:        "c_ubig",
		Description: "c_ubig",
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min: 6,
			Max: 6,
		},
	})

	schema = append(schema, basicgorm.Fields{ //n_celu
		Name:        "n_celu",
		Description: "n_celu",
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min: 9,
			Max: 24,
		},
	})

	schema = append(schema, basicgorm.Fields{ //n_tele
		Name:        "n_tele",
		Description: "n_tele",
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Min: 9,
			Max: 24,
		},
	})

	return schema
}

func (s *Sucursal) GetSchemaInsert() []basicgorm.Fields {
	return s.getSchema()
}

func (s *Sucursal) GetTableName() string {
	return s.name
}

func (s *Sucursal) GetSchemaUpdate() []basicgorm.Fields {
	var update []basicgorm.Fields
	tmp := s.getSchema()
	for _, v := range tmp {
		if v.Update {
			update = append(update, v)
		}
	}
	return update
}

func (s *Sucursal) GetSchemaDelete() []basicgorm.Fields {
	var delete []basicgorm.Fields
	tmp := s.getSchema()
	for _, v := range tmp {
		if v.Where || v.PrimaryKey {
			delete = append(delete, v)
		}
	}
	return delete
}

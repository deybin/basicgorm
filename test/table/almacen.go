package table

import (
	"github.com/deybin/basicgorm"
)

type Store struct {
	name string
}

func (s *Store) New() *Store {
	s.name = "requ_" + "almacen"
	return s
}
func (s *Store) getSchema() []basicgorm.Fields {
	var schema []basicgorm.Fields
	schema = append(schema, basicgorm.Fields{ //c_sucu
		Name:        "c_sucu",
		Description: "c_sucu",
		Required:    true,
		Where:       true,
		Type:        basicgorm.String,
		ValidateType: basicgorm.TypeStrings{
			Expr: basicgorm.Number(),
			Min:  3,
			Max:  3,
		},
	})
	schema = append(schema, basicgorm.Fields{ //c_alma
		Name:        "c_alma",
		Description: "c_alma",
		Required:    true,
		PrimaryKey:  true,
		Type:        "string",
		ValidateType: basicgorm.TypeStrings{
			Expr: basicgorm.Number(),
			Min:  3,
			Max:  3,
		},
	})
	schema = append(schema, basicgorm.Fields{ //l_alma
		Name:        "l_alma",
		Description: "l_alma",
		Required:    true,
		Update:      true,
		Type:        "string",
		ValidateType: basicgorm.TypeStrings{
			Min:       3,
			Max:       50,
			LowerCase: true,
		},
	})
	return schema
}

func (s *Store) GetSchemaInsert() []basicgorm.Fields {
	return s.getSchema()
}

func (s *Store) GetTableName() string {
	return s.name
}

func (s *Store) GetSchemaUpdate() []basicgorm.Fields {
	var update []basicgorm.Fields
	tmp := s.getSchema()
	for _, v := range tmp {
		if v.Update || v.PrimaryKey || v.Where {
			update = append(update, v)
		}
	}
	return update
}

func (s *Store) GetSchemaDelete() []basicgorm.Fields {
	var delete []basicgorm.Fields
	tmp := s.getSchema()
	for _, v := range tmp {
		if v.Where || v.PrimaryKey {
			delete = append(delete, v)
		}
	}
	return delete
}

func (s *Store) GetId() string {
	return "id_name"
}

func (s *Store) SetTable(name string) {
	s.name = name
}

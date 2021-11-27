package db

type Entity interface {
	AsMap() map[string]interface{}
	GetEntityId() EntityId
	FromMap(map[string]interface{}) error
}

type EntityId string

type Filter struct {
	filter map[string]interface{}
}

type DbEntity struct {
	Entity
}

func NewFilter(fieldName string, value interface{}) *Filter {
	bsonM := make(map[string]interface{})
	filter := &Filter{
		filter: bsonM,
	}
	return filter
}

func (f *Filter) Filter(fieldName string, value interface{}) *Filter {
	f.filter[fieldName] = value
	return f
}

package helpers

import (
	"github.com/lohithk3345/voting_system/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Filter struct {
	doc bson.M
}

var Filters = Filter{}

func (f *Filter) ByID(id types.ID) Filter {
	return Filter{doc: bson.M{"_id": id}}
}

func (f *Filter) ByEmail(email string) Filter {
	return Filter{doc: bson.M{"email": email}}
}

func (f *Filter) ByName(name string) Filter {
	return Filter{doc: bson.M{"name": name}}
}

func (f *Filter) ByUserID(id types.ID) Filter {
	return Filter{doc: bson.M{"userId": id}}
}

func (f *Filter) Doc() bson.M {
	return f.doc
}

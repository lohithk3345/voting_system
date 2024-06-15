package helpers

import "go.mongodb.org/mongo-driver/bson"

type UserUpdate struct {
	doc bson.M
}

func (f *UserUpdate) Name(name string) UserUpdate {
	return UserUpdate{doc: bson.M{"$set": bson.M{"name": name}}}
}

func (f *UserUpdate) Age(age int) UserUpdate {
	return UserUpdate{doc: bson.M{"$set": bson.M{"age": age}}}
}

func (f *UserUpdate) Doc() bson.M {
	return f.doc
}

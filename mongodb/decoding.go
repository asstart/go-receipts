package mongodb

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Configure custom type mapping for map[string]interface{}
//
// Problem description:
// You have a struct with the following fields
//
// type CustomNestedMapStruct struct {
// 		ID string
// 		Data map[string]interface{}
// }
//
// You put a time.Time to the Data
//
// During the saving to a mongodb collection, there's no problem
// This will be saved the same way as if you have structure like this
//
// type CustomFlatStructure struct {
// 		ID string
// 		Date time.Time
// }
//
// But during the decoding this back, you will face with problem
// that CustomNestedMapStruct.Data will contains primitive.DateTime instead time.Time
//
// To fix this, you need to register TypeMapEntry
//
// Look at customRegistry

type CustomNestedMapStruct struct {
	ID   string
	Data map[string]interface{}
}

type CustomFlatStructure struct {
	ID   string
	Date time.Time
}

func ExecWithNestedMapDefaultMapType(conStr string, db string, coll string) (CustomNestedMapStruct, error) {
	con, err := getConnection(conStr)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	c := con.Database(db).Collection(coll)

	err = insertNested(c)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	res, err := readNestedDefault(c)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	return res, nil
}

func ExecWithNestedMapCustomMapType(conStr string, db string, coll string) (CustomNestedMapStruct, error) {
	con, err := getConnection(conStr)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	c := con.Database(db).Collection(coll)

	err = insertNested(c)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	reg := customRegistry()

	res, err := readNestedWithCustomMapType(c, reg)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	return res, nil
}

func ExecWithFlat(conStr, db string, coll string) (CustomFlatStructure, error) {
	con, err := getConnection(conStr)
	if err != nil {
		return CustomFlatStructure{}, err
	}

	c := con.Database(db).Collection(coll)

	err = insertFlat(c)
	if err != nil {
		return CustomFlatStructure{}, err
	}
	res, err := readFlat(c)
	if err != nil {
		return CustomFlatStructure{}, err
	}
	return res, nil
}

func insertNested(c *mongo.Collection) error {
	doc := CustomNestedMapStruct{
		ID: "nested",
		Data: map[string]interface{}{
			"createdAt": time.Now(),
		},
	}

	_, err := c.InsertOne(
		context.Background(),
		doc,
		nil,
	)
	return err
}

func readNestedDefault(c *mongo.Collection) (CustomNestedMapStruct, error) {
	var res CustomNestedMapStruct

	err := c.FindOne(
		context.Background(),
		primitive.M{
			"id": "nested",
		},
		nil,
	).Decode(&res)

	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	return res, nil
}

func customRegistry() *bsoncodec.Registry {
	rb := bsoncodec.NewRegistryBuilder()

	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)

	rb.RegisterTypeMapEntry(bsontype.DateTime, reflect.TypeOf(time.Time{}))
	return rb.Build()
}

func readNestedWithCustomMapType(c *mongo.Collection, registry *bsoncodec.Registry) (CustomNestedMapStruct, error) {
	sr := c.FindOne(
		context.Background(),
		primitive.M{
			"id": "nested",
		},
		nil,
	)

	if sr.Err() != nil {
		return CustomNestedMapStruct{}, sr.Err()
	}

	raw, err := sr.DecodeBytes()
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	var res CustomNestedMapStruct
	err = bson.UnmarshalWithRegistry(registry, raw, &res)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	return res, nil
}

func insertFlat(c *mongo.Collection) error {
	doc := CustomFlatStructure{
		ID:   "flat",
		Date: time.Now(),
	}

	_, err := c.InsertOne(
		context.Background(),
		doc,
		nil,
	)
	return err
}

func readFlat(c *mongo.Collection) (CustomFlatStructure, error) {
	var res CustomFlatStructure

	err := c.FindOne(
		context.Background(),
		primitive.M{
			"id": "flat",
		},
		nil,
	).Decode(&res)

	if err != nil {
		return CustomFlatStructure{}, err
	}

	return res, nil
}

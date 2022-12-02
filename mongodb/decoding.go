package mongodb

import (
	"context"
	"fmt"
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

func ExecWithNestedMapDefaultMapType(conStr string, db string, coll string, id_postfix string) (CustomNestedMapStruct, error) {
	con, err := getConnection(conStr)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	c := con.Database(db).Collection(coll)

	id := fmt.Sprintf("nested_%v", id_postfix)

	err = insertNested(c, id)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	res, err := readNestedDefault(c, id)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	return res, nil
}

func ExecWithNestedMapAllTypes(conStr string, db string, coll string, id_postfix string) (CustomNestedMapStruct, error) {
	con, err := getConnection(conStr)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	id := fmt.Sprintf("rich_nested_%v", id_postfix)

	c := con.Database(db).Collection(coll)

	err = insertNestedAllTypes(c, id)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	res, err := readNestedDefault(c, id)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	return res, nil
}

func ExecWithNestedMapAllTypesCustomRegister(conStr string, db string, coll string, id_postfix string) (CustomNestedMapStruct, error) {
	con, err := getConnection(conStr)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	id := fmt.Sprintf("rich_nested_%v", id_postfix)

	c := con.Database(db).Collection(coll)

	err = insertNestedAllTypes(c, id)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	reg := customRegistry()

	res, err := readNestedWithCustomMapType(c, reg, id)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	return res, nil
}

func ExecWithNestedMapCustomMapType(conStr string, db string, coll string, id_postfix string) (CustomNestedMapStruct, error) {
	con, err := getConnection(conStr)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	c := con.Database(db).Collection(coll)

	id := fmt.Sprintf("nested_%v", id_postfix)

	err = insertNested(c, id)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}

	reg := customRegistry()

	res, err := readNestedWithCustomMapType(c, reg, id)
	if err != nil {
		return CustomNestedMapStruct{}, err
	}
	return res, nil
}

func ExecWithFlat(conStr, db string, coll string, id_postfix string) (CustomFlatStructure, error) {
	con, err := getConnection(conStr)
	if err != nil {
		return CustomFlatStructure{}, err
	}

	c := con.Database(db).Collection(coll)

	id := fmt.Sprintf("flat_%v", id_postfix)

	err = insertFlat(c, id)
	if err != nil {
		return CustomFlatStructure{}, err
	}
	res, err := readFlat(c, id)
	if err != nil {
		return CustomFlatStructure{}, err
	}
	return res, nil
}

func insertNested(c *mongo.Collection, id string) error {
	doc := CustomNestedMapStruct{
		ID: id,
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

func insertNestedAllTypes(c *mongo.Collection, id string) error {
	now := time.Now()
	doc := CustomNestedMapStruct{
		ID: id,
		Data: map[string]interface{}{
			"int8":    int8(1),
			"int16":   int16(2),
			"int32":   int32(3),
			"int64":   int64(4),
			"int":     int(5),
			"uint8":   uint8(6),
			"uint16":  uint16(7),
			"uint32":  uint32(8),
			"uint64":  uint64(9),
			"uint":    uint(10),
			"float32": float32(1.4),
			"float64": float64(2.3),
			"bool":    true,
			// "complex64":                        complex64(complex(1, 1)),
			// "complex128":                       complex128(complex(1, 1)),
			"string":                           "some string",
			"byte":                             byte(11),
			"rune":                             rune(12),
			"array":                            [3]int{1, 2, 3},
			"slice":                            []int{1, 2, 3},
			"time.Time":                        time.Now(),
			"map[string]string":                map[string]string{"1": "11", "2": "22"},
			"map[string]interface - primitive": map[string]interface{}{"1": 1, "2": 3.5, "3": true},
			"struct":                           CustomFlatStructure{"1", time.Now()},
			"map[string]interface - with nested types": map[string]interface{}{"1": CustomFlatStructure{"1", time.Now()}},
			"time.Time pointer":                        &now,
		},
	}

	_, err := c.InsertOne(
		context.Background(),
		doc,
		nil,
	)

	return err
}

func readNestedDefault(c *mongo.Collection, id string) (CustomNestedMapStruct, error) {
	var res CustomNestedMapStruct

	err := c.FindOne(
		context.Background(),
		primitive.M{
			"id": id,
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
	rb.RegisterTypeMapEntry(bson.TypeArray, reflect.TypeOf([]interface{}{}))

	return rb.Build()
}

func readNestedWithCustomMapType(c *mongo.Collection, registry *bsoncodec.Registry, id string) (CustomNestedMapStruct, error) {
	sr := c.FindOne(
		context.Background(),
		primitive.M{
			"id": id,
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

func insertFlat(c *mongo.Collection, id string) error {
	doc := CustomFlatStructure{
		ID:   id,
		Date: time.Now(),
	}

	_, err := c.InsertOne(
		context.Background(),
		doc,
		nil,
	)
	return err
}

func readFlat(c *mongo.Collection, id string) (CustomFlatStructure, error) {
	var res CustomFlatStructure

	err := c.FindOne(
		context.Background(),
		primitive.M{
			"id": id,
		},
		nil,
	).Decode(&res)

	if err != nil {
		return CustomFlatStructure{}, err
	}

	return res, nil
}

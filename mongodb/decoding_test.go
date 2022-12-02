package mongodb_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/asstart/go-receipts/mongodb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCustomMapType(t *testing.T) {
	conStr := "mongodb://test:test@localhost:27017/"
	db := "test"
	coll := "test_decoding"

	rand.Seed(time.Now().UnixMilli())
	id_postfix := fmt.Sprintf("%v", rand.Int())

	flat, err := mongodb.ExecWithFlat(conStr, db, coll, id_postfix)
	assert.Nil(t, err)
	assert.NotEqual(t, flat.Date, time.Time{}) // check that it's not default value(empty value)
	assert.IsType(t, flat.Date, time.Time{})

	nestedDefault, err := mongodb.ExecWithNestedMapDefaultMapType(conStr, db, coll, id_postfix)
	assert.Nil(t, err)
	assert.IsType(t, nestedDefault.Data["createdAt"], primitive.DateTime(0))

	nestedCstm, err := mongodb.ExecWithNestedMapCustomMapType(conStr, db, coll, id_postfix)
	assert.Nil(t, err)
	assert.IsType(t, nestedCstm.Data["createdAt"], time.Time{})
}

func TestReadAllTypesMapDefault(t *testing.T) {
	conStr := "mongodb://test:test@localhost:27017/"
	db := "test"
	coll := "test_decoding"
	rand.Seed(time.Now().UnixMilli())
	id_postfix := fmt.Sprintf("%v", rand.Int())

	mp, err := mongodb.ExecWithNestedMapAllTypes(conStr, db, coll, id_postfix)
	assert.Nil(t, err)
	for k, v := range mp.Data {
		fmt.Printf("value type: %T\n", v)
		fmt.Printf("k v: %v:%v\n", k, v)
		fmt.Println("---")
	}
}

func TestReadAllTypesMapCustomRegister(t *testing.T) {
	conStr := "mongodb://test:test@localhost:27017/"
	db := "test"
	coll := "test_decoding"
	rand.Seed(time.Now().UnixMilli())
	id_postfix := fmt.Sprintf("%v", rand.Int())

	mp, err := mongodb.ExecWithNestedMapAllTypesCustomRegister(conStr, db, coll, id_postfix)
	assert.Nil(t, err)
	for k, v := range mp.Data {
		fmt.Printf("value type: %T\n", v)
		fmt.Printf("k v: %v:%v\n", k, v)
		fmt.Println("---")
	}
}

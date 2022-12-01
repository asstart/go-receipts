package mongodb_test

import (
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

	flat, err := mongodb.ExecWithFlat(conStr, db, coll)
	assert.Nil(t, err)
	assert.NotEqual(t, flat.Date, time.Time{}) // check that it's not default value(empty value)
	assert.IsType(t, flat.Date, time.Time{})

	nestedDefault, err := mongodb.ExecWithNestedMapDefaultMapType(conStr, db, coll)
	assert.Nil(t, err)
	assert.IsType(t, nestedDefault.Data["createdAt"], primitive.DateTime(0))

	nestedCstm, err := mongodb.ExecWithNestedMapCustomMapType(conStr, db, coll)
	assert.Nil(t, err)
	assert.IsType(t, nestedCstm.Data["createdAt"], time.Time{})
}

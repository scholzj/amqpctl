package tests

import (
	"testing"
	"../utils"
	"../operations"
	"qpid.apache.org/amqp"
)

func TestQuery(t *testing.T) {
	mgmtLink := MockMgmtLink{}
	args := []string{"query"}

	mgmtLink.ResponseProperties = map[string]interface{}{"statusCode": int32(200)}
	mgmtLink.ResponseBody = amqp.Map(map[interface{}]interface{}{"attributeNames": amqp.List([]interface{}{"name", "type", "identity"}), "results": amqp.List([]interface{}{amqp.List([]interface{}{"objectA", "typeX", "id1"}), amqp.List([]interface{}{"objectA", "typeX", "id1"})})})

	operations.Query(args, utils.MgmtLink(&mgmtLink))

	if mgmtLink.RequestProperties["operation"] != "QUERY" {
		t.Error("Request propert contains wrong operation")
	}

	if !(contains(mgmtLink.RequestBody["attributeNames"].([]interface{}), "name") && contains(mgmtLink.RequestBody["attributeNames"].([]interface{}), "type") && contains(mgmtLink.RequestBody["attributeNames"].([]interface{}), "identity")) {
		t.Error("Request body contains wrong attributes")
	}
}

func contains(list []interface{}, search string) bool {
	for _, s := range list {
		if s == search {
			return true
		}
	}
	return false
}
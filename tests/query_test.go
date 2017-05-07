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

	if !(contains(mgmtLink.RequestBody["attributeNames"].([]string), "name") && contains(mgmtLink.RequestBody["attributeNames"].([]string), "type") && contains(mgmtLink.RequestBody["attributeNames"].([]string), "identity")) {
		t.Error("Request body contains wrong attributes")
	}
}

func TestQueryWithEntity(t *testing.T) {
	mgmtLink := MockMgmtLink{}
	args := []string{"query", "typeX"}

	mgmtLink.ResponseProperties = map[string]interface{}{"statusCode": int32(200)}
	mgmtLink.ResponseBody = amqp.Map(map[interface{}]interface{}{"attributeNames": amqp.List([]interface{}{"name", "type", "identity"}), "results": amqp.List([]interface{}{amqp.List([]interface{}{"objectA", "typeX", "id1"}), amqp.List([]interface{}{"objectA", "typeX", "id1"})})})

	operations.Query(args, utils.MgmtLink(&mgmtLink))

	if mgmtLink.RequestProperties["operation"] != "QUERY" || mgmtLink.RequestProperties["entityType"] != "typeX" {
		t.Error("Request propert contains wrong operation")
	}

	if !(contains(mgmtLink.RequestBody["attributeNames"].([]string), "name") && contains(mgmtLink.RequestBody["attributeNames"].([]string), "type") && contains(mgmtLink.RequestBody["attributeNames"].([]string), "identity")) {
		t.Error("Request body contains wrong attributes")
	}
}

func TestQueryWithEntityAndAttributes(t *testing.T) {
	mgmtLink := MockMgmtLink{}
	args := []string{"query", "typeX", "att1", "att2"}

	mgmtLink.ResponseProperties = map[string]interface{}{"statusCode": int32(200)}
	mgmtLink.ResponseBody = amqp.Map(map[interface{}]interface{}{"attributeNames": amqp.List([]interface{}{"att1", "att2"}), "results": amqp.List([]interface{}{amqp.List([]interface{}{"objectA", "typeX"}), amqp.List([]interface{}{"objectA", "typeX"})})})

	operations.Query(args, utils.MgmtLink(&mgmtLink))

	if mgmtLink.RequestProperties["operation"] != "QUERY" || mgmtLink.RequestProperties["entityType"] != "typeX" {
		t.Error("Request propert contains wrong operation")
	}

	if !(contains(mgmtLink.RequestBody["attributeNames"].([]string), "att1") && contains(mgmtLink.RequestBody["attributeNames"].([]string), "att2")) {
		t.Error("Request body contains wrong attributes")
	}
}

/*func TestQueryWithStatus(t *testing.T) {
	mgmtLink := MockMgmtLink{}
	args := []string{"query"}

	mgmtLink.ResponseProperties = map[string]interface{}{"statusCode": int32(404)}
	mgmtLink.ResponseBody = amqp.Map(map[interface{}]interface{}{"attributeNames": amqp.List([]interface{}{"name", "type", "identity"}), "results": amqp.List([]interface{}{amqp.List([]interface{}{"objectA", "typeX", "id1"}), amqp.List([]interface{}{"objectA", "typeX", "id1"})})})

	operations.Query(args, utils.MgmtLink(&mgmtLink))

	if mgmtLink.RequestProperties["operation"] != "QUERY" {
		t.Error("Request propert contains wrong operation")
	}

	if !(contains(mgmtLink.RequestBody["attributeNames"].([]string), "name") && contains(mgmtLink.RequestBody["attributeNames"].([]string), "type") && contains(mgmtLink.RequestBody["attributeNames"].([]string), "identity")) {
		t.Error("Request body contains wrong attributes")
	}
}*/

func contains(list []string, search string) bool {
	for _, s := range list {
		if s == search {
			return true
		}
	}
	return false
}
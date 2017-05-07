package tests

type MockMgmtLink struct {
	RequestProperties map[string]interface{}
	RequestBody map[interface{}]interface{}
	ResponseProperties map[string]interface{}
	ResponseBody interface{}
	Err error
}

func (l *MockMgmtLink) Connect() (err error) {
	err = nil

	return
}

func (l *MockMgmtLink) Close() {
	return
}

func (l *MockMgmtLink) Operation(reqProperties map[string]interface{}, reqBody map[interface{}]interface{}) (respProperties map[string]interface{}, respBody interface{}, err error) {
	err = l.Err
	respProperties = l.ResponseProperties
	respBody = l.ResponseBody

	l.RequestProperties = reqProperties
	l.RequestBody = reqBody

	return
}
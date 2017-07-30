package mgmtlink

type MgmtLink interface {
	Connect() error
	Close()
	Operation(map[string]interface{}, map[interface{}]interface{}) (map[string]interface{}, interface{}, error)
}
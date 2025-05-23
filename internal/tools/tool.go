package tools

type Tool interface {
	Name() string
	Description() string
	InputSchema() map[string]interface{}
	Execute(args map[string]interface{}) (string, error)
}
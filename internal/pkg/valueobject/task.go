package valueobject

type TaskType string

const (
	TaskTypeHandlePoolCreated = "handle_pool_created"
)

type TaskHandlePoolCreatedPayload struct {
	PoolAddress string `json:"poolAddress"`
}

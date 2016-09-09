package dds

// A Storage is a place to store entities
type Storage interface {
	Get(id string) (interface{}, error)

	Create(i interface{}) (string, error)

	Delete(id string) error

	Update(i interface{}) error
}

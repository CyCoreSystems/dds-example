package dds

// A Storage is a place to store entities
type Storage interface {
	Get(id string)

	List()

	Delete(id string) error

	Update(id string, i interface{}) error
}

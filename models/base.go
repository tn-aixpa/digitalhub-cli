package models

type BaseModel interface {
	GetSpec() Spec
	//GetMetadata() Metadata example
}

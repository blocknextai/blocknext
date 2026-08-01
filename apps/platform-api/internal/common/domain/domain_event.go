package domain

type DomainEvent interface {
	EventName() string
}

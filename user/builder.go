package user

import "log"

type ServiceBuilder interface {
	Repository(repository Repository) ServiceBuilder
	Client(client Client) ServiceBuilder
	PasswordManager(passwordManager PasswordManager) ServiceBuilder
	Validator(validator Validator) ServiceBuilder
	Build() Service
}

func NewServiceBuilder() ServiceBuilder {
	return &serviceBuilder{}
}

type serviceBuilder struct {
	repository      Repository
	client          Client
	passwordManager PasswordManager
	validator       Validator
}

func (sb *serviceBuilder) Repository(repository Repository) ServiceBuilder {
	sb.repository = repository
	return sb
}

func (sb *serviceBuilder) Client(client Client) ServiceBuilder {
	sb.client = client
	return sb
}

func (sb *serviceBuilder) PasswordManager(passwordManager PasswordManager) ServiceBuilder {
	sb.passwordManager = passwordManager
	return sb
}

func (sb *serviceBuilder) Validator(validator Validator) ServiceBuilder {
	sb.validator = validator
	return sb
}

func (sb *serviceBuilder) Build() Service {
	tc, err := NewTokenClient("5003")
	if err != nil {
		log.Fatalln(err)
	}

	return &service{
		repository:      sb.repository,
		client:          sb.client,
		passwordManager: sb.passwordManager,
		v:               sb.validator,
		tc:              tc,
	}
}

package user

type ServiceBuilder interface {
	Repository(repository Repository) ServiceBuilder
	Client(client Client) ServiceBuilder
	TokenClient(token TokenClient) ServiceBuilder
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
	token           TokenClient
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

func (sb *serviceBuilder) TokenClient(token TokenClient) ServiceBuilder {
	sb.token = token
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
	return &service{
		repository: sb.repository,
		client:     sb.client,
		password:   sb.passwordManager,
		validate:   sb.validator,
		token:      sb.token,
	}
}

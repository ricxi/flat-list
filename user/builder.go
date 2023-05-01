package user

// ! ServiceBuilder may be deprecated in favour of the
// ! idiomatic functional options pattern. I'm going to
// ! keep this here until I refactor some of the tests.
// ! I also want to see how it stacks up to the fo pattern
// ! when I do more benchmarking and profiling.
type ServiceBuilder interface {
	Repository(repository Repository) ServiceBuilder
	MailerClient(client MailerClient) ServiceBuilder
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
	mailer          MailerClient
	token           TokenClient
	passwordManager PasswordManager
	validator       Validator
}

func (sb *serviceBuilder) Repository(repository Repository) ServiceBuilder {
	sb.repository = repository
	return sb
}

func (sb *serviceBuilder) MailerClient(client MailerClient) ServiceBuilder {
	sb.mailer = client
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
		mailer:     sb.mailer,
		password:   sb.passwordManager,
		validate:   sb.validator,
		token:      sb.token,
	}
}

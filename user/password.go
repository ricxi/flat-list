package user

import "golang.org/x/crypto/bcrypt"

type PasswordService interface {
	GenerateHash(password string) (string, error)
	CompareHashWith(hashedPassword, password string) error
}

type passwordService struct {
	cost int
}

func NewPasswordService(cost int) *passwordService {
	ps := passwordService{}

	if cost == 0 {
		ps.cost = bcrypt.DefaultCost
	}

	return &ps
}

// GenerateHash creates a hash for a given password
func (ps *passwordService) GenerateHash(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), ps.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (ps *passwordService) CompareHashWith(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	// this should not happen
	// if errors.Is(err, bcrypt.ErrHashTooShort) {
	// 	return nil, err
	// }
}

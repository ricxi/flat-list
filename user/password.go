package user

import "golang.org/x/crypto/bcrypt"

type PasswordManager interface {
	GenerateHash(password string) (string, error)
	CompareHashWith(hashedPassword, password string) error
}

type passwordManager struct {
	cost int
}

func NewPasswordManager(cost int) PasswordManager {
	pm := passwordManager{}

	if cost == 0 {
		pm.cost = bcrypt.DefaultCost
	}

	return &pm
}

// GenerateHash creates a hash for a given password
func (pm *passwordManager) GenerateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (pm *passwordManager) CompareHashWith(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	// this should not happen
	// if errors.Is(err, bcrypt.ErrHashTooShort) {
	// 	return nil, err
	// }
}

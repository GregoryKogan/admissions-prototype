package crypto

import (
	"crypto/rand"
	"crypto/subtle"

	"golang.org/x/crypto/argon2"
)

type CryptoService interface {
	GenerateHash(password []byte) (*HashedPassword, error)
	Compare(password []byte, hashedPassword *HashedPassword) (bool, error)
}

type HashedPassword struct {
	Hash      []byte
	Salt      []byte
	Algorithm string
}

type CryptoServiceImpl struct {
	params *Argon2idParams
}

func NewCryptoService() CryptoService {
	return &CryptoServiceImpl{
		params: defaultArgon2idParams(),
	}
}

func (s *CryptoServiceImpl) GenerateHash(password []byte) (*HashedPassword, error) {
	salt, err := randomSecret(s.params.SaltLen)
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey(password, salt, s.params.Time, s.params.Memory, s.params.Threads, s.params.KeyLen)

	return &HashedPassword{
		Hash:      hash,
		Salt:      salt,
		Algorithm: "argon2id",
	}, nil
}

func (s *CryptoServiceImpl) Compare(password []byte, hashedPassword *HashedPassword) (bool, error) {
	hash := argon2.IDKey(password, hashedPassword.Salt, s.params.Time, s.params.Memory, s.params.Threads, s.params.KeyLen)
	equal := subtle.ConstantTimeCompare(hash, hashedPassword.Hash) == 1
	return equal, nil
}

type Argon2idParams struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

func defaultArgon2idParams() *Argon2idParams {
	// Source: https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#password-hashing-algorithms
	// Recommendations for the minimum memory size (m), the minimum number of iterations (t), and the degree of parallelism (p) as follows:
	// m=47104 (46 MiB), t=1, p=1 (Do not use with Argon2i)
	// m=19456 (19 MiB), t=2, p=1 (Do not use with Argon2i)
	// m=12288 (12 MiB), t=3, p=1
	// m=9216 (9 MiB), t=4, p=1
	// m=7168 (7 MiB), t=5, p=1
	return newArgon2idParams(3, 12288, 1, 32, 16)
}

func newArgon2idParams(time, memory uint32, threads uint8, keyLen, saltLen uint32) *Argon2idParams {
	return &Argon2idParams{
		Time:    time,
		Memory:  memory,
		Threads: threads,
		KeyLen:  keyLen,
		SaltLen: saltLen,
	}
}

func randomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}

	return secret, nil
}

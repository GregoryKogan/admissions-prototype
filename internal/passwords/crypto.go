package passwords

import (
	"crypto/rand"
	"crypto/subtle"

	"golang.org/x/crypto/argon2"
)

type Argon2idParams struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

type HashedPassword struct {
	Hash      []byte
	Salt      []byte
	Algorithm string
}

func NewArgon2idParams(time, memory uint32, threads uint8, keyLen, saltLen uint32) *Argon2idParams {
	return &Argon2idParams{
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
		saltLen: saltLen,
	}
}

func DefaultArgon2idParams() *Argon2idParams {
	// Source: https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#password-hashing-algorithms
	// Recommendations for the minimum memory size (m), the minimum number of iterations (t), and the degree of parallelism (p) as follows:
	// m=47104 (46 MiB), t=1, p=1 (Do not use with Argon2i)
	// m=19456 (19 MiB), t=2, p=1 (Do not use with Argon2i)
	// m=12288 (12 MiB), t=3, p=1
	// m=9216 (9 MiB), t=4, p=1
	// m=7168 (7 MiB), t=5, p=1
	return NewArgon2idParams(3, 12288, 1, 32, 16)
}

func randomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (p *Argon2idParams) GenerateHash(password []byte, salt []byte) (*HashedPassword, error) {
	var err error
	if len(salt) == 0 {
		salt, err = randomSecret(p.saltLen)
	}
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey(password, salt, p.time, p.memory, p.threads, p.keyLen)

	return &HashedPassword{
		Hash:      hash,
		Salt:      salt,
		Algorithm: "argon2id",
	}, nil
}

func (p *Argon2idParams) Compare(password []byte, hashedPassword *HashedPassword) (bool, error) {
	newHash, err := p.GenerateHash(password, hashedPassword.Salt)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(newHash.Hash, hashedPassword.Hash) == 1, nil
}

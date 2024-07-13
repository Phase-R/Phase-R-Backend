package utils

/*func ComparePasswords(hashedPassword, password string) (bool, error) {
	pos := -1
	for i := 0; i < len(hashedPassword); i++ {
		if hashedPassword[i] == '.' {
			pos = i
			break
		}
	}

	if pos == -1 {
		return false, fmt.Errorf("invalid hashed password format")
	}

	b64Salt := hashedPassword[:pos]
	b64Hash := hashedPassword[pos+1:]

	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false, err
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(b64Hash)
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	if byteSliceComp(hash, storedHash) {
		return true, nil
	}

	return false, nil
}

func byteSliceComp(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

func PwdSaltAndHash(password string) (string, error) {
	hashedPassword, err := hashy(password)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func hashy(password string) (string, error) {
	salt, err := salted(16)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s.%s", b64salt, b64hash), nil
}

func salted(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
} */

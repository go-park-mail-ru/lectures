package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

var (
	plainPassword = []byte("qwerty123456")

	// соль должна быть для каждого юзера своя, не используте в таком виде
	salt = []byte{0xd7, 0xc2, 0xf2, 0x51, 0xaa, 0x6a, 0x4e, 0x7b}
)

// md5 - плохой вариант, подвержен брутфорс-атаке
func PasswordMD5(plainPassword []byte) []byte {
	return md5.New().Sum(plainPassword)
}

// https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Password_Storage_Cheat_Sheet.md

// bcrypt where PBKDF2 or scrypt support is not available.
func PasswordBcrypt(plainPassword []byte) []byte {
	passBcrypt, _ := bcrypt.GenerateFromPassword(plainPassword, 14)
	return passBcrypt
}

// PBKDF2 when FIPS certification or enterprise support on many platforms is required;
func PasswordPBKDF2(plainPassword []byte) []byte {
	return pbkdf2.Key(plainPassword, salt, 4096, 32, sha1.New)
}

// scrypt where resisting any/all hardware accelerated attacks is necessary but support isn’t.
func PasswordScrypt(plainPassword []byte) []byte {
	passScrypt, _ := scrypt.Key(plainPassword, salt, 1<<15, 8, 1, 32)
	return passScrypt
}

// Argon2 is the winner of the password hashing competition and should be considered as your first choice for new applications;
func PasswordArgon2(plainPassword []byte) []byte {
	return argon2.IDKey(plainPassword, salt, 1, 64*1024, 4, 32)
}

func main() {
	fmt.Printf("PasswordMD5: %x\n", PasswordMD5(plainPassword))
	fmt.Printf("PasswordBcrypt: %x\n", PasswordBcrypt(plainPassword))
	fmt.Printf("PasswordPBKDF2: %x\n", PasswordPBKDF2(plainPassword))
	fmt.Printf("PasswordScrypt: %x\n", PasswordScrypt(plainPassword))
	fmt.Printf("PasswordArgon2: %x\n", PasswordArgon2(plainPassword))
}

package main

import (
	"testing"
)

func BenchmarkMD5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PasswordMD5(plainPassword)
	}
}

func BenchmarkBcrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PasswordBcrypt(plainPassword)
	}
}

func BenchmarkPBKDF2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PasswordPBKDF2(plainPassword)
	}
}

func BenchmarkScrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PasswordScrypt(plainPassword)
	}
}

func BenchmarkArgon2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PasswordArgon2(plainPassword)
	}
}

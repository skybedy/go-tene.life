// news-password reads a password from stdin and prints only its bcrypt hash.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	data, err := io.ReadAll(io.LimitReader(os.Stdin, 1024))
	if err != nil {
		log.Fatal(err)
	}
	password := strings.TrimRight(string(data), "\r\n")
	if len(password) < 6 || len(password) > 72 {
		log.Fatal("Heslo musí mít 6–72 bajtů.")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(hash))
}

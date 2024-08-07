package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

// Остальные поля кроме Email в структуре User нас не интересуют (лишние аллокации).
type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	var user User

	for scanner.Scan() {
		// Пропускаем пользователей без совпадения домена
		if !strings.Contains(scanner.Text(), domain) {
			continue
		}

		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return result, err
		}

		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}

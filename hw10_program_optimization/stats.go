package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var errTooMuchUsers = errors.New("too much users")

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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]string

// getUsers считывает построчно информацию о пользователях.
func getUsers(r io.Reader) (result users, err error) {
	usersCount := 0

	// Считаем построчно в массив
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// Если пользователей больше 100_000 - ошибка
		if usersCount >= len(result) {
			return result, errTooMuchUsers
		}

		result[usersCount] = scanner.Text()
		usersCount++
	}

	return
}

// countDomains подсчитывает количество доменов 1 уровня в электронной почте пользователей.
func countDomains(u users, domain string) (result DomainStat, err error) {
	result = make(DomainStat)
	var user User

	for idx := range u {
		// Пропускаем пользователей без совпадения домена
		if !strings.Contains(u[idx], domain) {
			continue
		}

		if err = json.Unmarshal([]byte(u[idx]), &user); err != nil {
			return
		}

		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}

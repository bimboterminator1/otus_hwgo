package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain), nil
}

type users []User

func getUsers(r io.Reader, domain string) (result users, err error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		str := scanner.Text()
		if strings.Contains(str, domain) {
			lines = append(lines, scanner.Text())
		}
	}
	result = make(users, len(lines))
	for i, line := range lines {
		u := User{}
		if err := json.Unmarshal([]byte(line), &u); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}
		result[i] = u
	}

	return result, nil
}

func countDomains(u users, domain string) DomainStat {
	result := make(DomainStat)
	for _, user := range u {
		matched := strings.HasSuffix(user.Email, "."+domain)
		if matched && strings.Contains(user.Email, "@") {
			str := strings.ToLower(user.Email)
			domain := strings.SplitN(str, "@", 2)[1]
			num := result[domain]
			num++
			result[domain] = num
		}
	}
	return result
}

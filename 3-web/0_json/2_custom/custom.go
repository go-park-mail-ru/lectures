package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Company string

// MarshalJSON удовлетворяет интерфейсу json.Marshaler
func (c Company) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("ООО «%s»", c))
}

// UnmarshalJSON удовлетворяет интерфейсу json.Unmarshaler
func (c *Company) UnmarshalJSON(data []byte) error {
	var company string
	if err := json.Unmarshal(data, &company); err != nil {
		return err
	}

	*c = Company(strings.TrimSuffix(strings.TrimPrefix(company, "ООО «"), "»"))
	return nil
}

type User struct {
	ID             int
	PassportSeries int
	PassportNumber int
	Company        Company
}

// MarshalJSON удовлетворяет интерфейсу json.Marshaler
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID       int     `json:"user_id"`
		Passport string  `json:"passport"`
		Company  Company `json:"company"`
	}{
		ID:       u.ID,
		Passport: fmt.Sprintf("%d %d", u.PassportSeries, u.PassportNumber),
		Company:  u.Company,
	})
}

// UnmarshalJSON удовлетворяет интерфейсу json.Unmarshaler
func (u *User) UnmarshalJSON(data []byte) error {
	var temp struct {
		ID       int     `json:"user_id"`
		Passport string  `json:"passport"`
		Company  Company `json:"company"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	u.ID = temp.ID
	u.Company = temp.Company

	passportParts := strings.SplitN(temp.Passport, " ", 2)
	u.PassportSeries, _ = strconv.Atoi(passportParts[0])
	u.PassportNumber, _ = strconv.Atoi(passportParts[1])

	return nil
}

func main() {
	u := &User{
		ID:             42,
		PassportSeries: 4512,
		PassportNumber: 123456,
		Company:        "VK",
	}
	result, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json string: %s\n", string(result))

	var unmarshalledUser User
	err = json.Unmarshal(result, &unmarshalledUser)
	if err != nil {
		panic(err)
	}

	fmt.Printf("unmarshalledUser: %+v\n", unmarshalledUser)
}

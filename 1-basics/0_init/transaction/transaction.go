package transaction

import "fmt"

type Account struct {
	Title string
}

type Transaction struct {
	ID string
	Amount float64
	Account Account
}

func (t Transaction) String() string {
	return fmt.Sprintf("id: %s, account: %s, amount: %f", t.ID, t.Account.Title, t.Amount)
}
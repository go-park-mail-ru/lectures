package transaction

func NewList() List {
	return List{
		transactions: map[string]Transaction{},
	}
}

type List struct {
	transactions map[string]Transaction
}

func (l *List) Add(transaction Transaction) {
	l.transactions[transaction.ID] = transaction
}

func (l List) Get() []Transaction  {
	list := make([]Transaction, 0, len(l.transactions))
	for _, transaction := range l.transactions {
		list = append(list, transaction)
	}

	return list
}

package main

// Ledger contains the leger of all the transations that have
// happened in the system. So we have a complete list of transactions
// as they have happened.
type Ledger struct {
	UserId int64 `sql:",pk"`
	OrgId  int64 `sql:",pk"`
	Amount float64

	CreatedAt time.Time
	UpdatedAt time.Time
}

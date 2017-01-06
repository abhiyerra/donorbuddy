package main

import (
	"gopkg.in/pg.v5"
)

type Org struct {
	Name     string
	EIN      string
	City     string
	State    string
	Country  string
	Category string
	StripeCustomerID
}

type User struct {
}

type UserOrg struct {
}

type PaymentLedger struct {
}

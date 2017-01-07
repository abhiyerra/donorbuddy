package main

// Org is a table of the organizations that we will be supporting for
// donations..
type Org struct {
	Id int64

	Name string
	EIN  string

	Address string
	City    string
	State   string
	Country string

	Category string

	Verified bool

	StripeCustomerID string `json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Users   []User   `pg:",many2many:user_orgs"`
	Ledgers []Ledger `pg:",many2many:ledgers"`
}

func showOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		orgId = vars["orgId"]
		org   Org
		err   error
	)

	if org.Id, err = strconv.ParseInt(orgId, 10, 64); err != nil {
		respondJson(w, r, err)
		return
	}

	if err = config.DB.Select(&org); err != nil {
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, org)
}

func searchOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		orgs []Org
	)

	err := config.DB.Model(&orgs).Where().Limit(50).Select()
	if err != nil {
		respondJson(w, r, err)
	}

	respondJson(w, r, orgs)
}

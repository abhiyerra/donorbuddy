package main

// UserOrg creates a map of the org that the user is donating to. A
// user can have more than one of the same user and org map. That is
// how we make donating more to a certain cause easier.
type UserOrg struct {
	UserId int64 `sql:",pk"`
	OrgId  int64 `sql:",pk"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func putUserOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars    = mux.Vars(r)
		orgId   = vars["orgId"]
		userOrg UserOrg
		err     error
	)

	if userOrg.OrgId, err = strconv.ParseInt(orgId, 10, 64); err != nil {
		respondJson(w, r, err)
		return
	}

	if err = config.DB.Insert(&userOrg); err != nil {
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, userOrg)
}

func deleteUserOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars    = mux.Vars(r)
		orgId   = vars["orgId"]
		userOrg UserOrg
		err     error
	)

	if userOrg.OrgId, err = strconv.ParseInt(orgId, 10, 64); err != nil {
		respondJson(w, r, err)
		return
	}

	_, err = config.DB.Model(&userOrg).Where("user_id = ?user_id and org_id = ?org_id").Limit(1).Delete()
	if err != nil {
		respondJson(w, r, err)
		return

	}

	respondJson(w, r, struct{}{})
}

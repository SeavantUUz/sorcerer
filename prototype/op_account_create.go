package prototype

func (a *AccountCreateOperation) GetSigner(auths *map[string]bool) {
	(*auths)[a.Creator.Value] = true
}

func (a *AccountCreateOperation) Validate() error {
	return nil
}

func (a *AccountCreateOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("account_create", (*Operation_Op1)(nil), (*AccountCreateOperation)(nil))
}

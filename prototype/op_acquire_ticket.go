package prototype

func (a *AcquireTicketOperation) GetSigner(auths *map[string]bool) {
	(*auths)[a.Account.Value] = true
}

func (a *AcquireTicketOperation) Validate() error {
	return nil
}

func (a *AcquireTicketOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("acquire_ticket", (*Operation_Op21)(nil), (*AcquireTicketOperation)(nil))
}

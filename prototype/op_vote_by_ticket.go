package prototype

func (a *VoteByTicketOperation) GetSigner(auths *map[string]bool) {
	(*auths)[a.Account.Value] = true
}

func (a *VoteByTicketOperation) Validate() error {
	return nil
}

func (a *VoteByTicketOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("vote_by_ticket", (*Operation_Op22)(nil), (*VoteByTicketOperation)(nil))
}

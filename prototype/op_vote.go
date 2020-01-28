package prototype

func (m *VoteOperation) GetSigner(auths *map[string]bool) {
	(*auths)[m.Voter.Value] = true
}

func (m *VoteOperation) Validate() error {
	return nil
}

func (m *VoteOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("vote", (*Operation_Op9)(nil), (*VoteOperation)(nil))
}

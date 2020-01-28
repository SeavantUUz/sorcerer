package prototype

func (m *BpVoteOperation) GetSigner(auths *map[string]bool) {
	(*auths)[m.Voter.Value] = true
}

func (m *BpVoteOperation) Validate() error {
	return nil
}

func (m *BpVoteOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("bp_vote", (*Operation_Op5)(nil), (*BpVoteOperation)(nil))
}

package prototype

func (m *StakeOperation) GetSigner(auths *map[string]bool) {
	(*auths)[m.From.Value] = true
}

func (m *StakeOperation) Validate() error {
	return nil
}

func (m *StakeOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("stake", (*Operation_Op17)(nil), (*StakeOperation)(nil))
}

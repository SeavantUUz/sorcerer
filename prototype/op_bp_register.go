package prototype

func (m *BpRegisterOperation) GetSigner(auths *map[string]bool) {
	(*auths)[m.Owner.Value] = true
}

func (m *BpRegisterOperation) Validate() error {
	return nil
}

func (m *BpRegisterOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("bp_register", (*Operation_Op3)(nil), (*BpRegisterOperation)(nil))
}

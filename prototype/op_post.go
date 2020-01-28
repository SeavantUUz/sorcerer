package prototype

func (m *PostOperation) GetSigner(auths *map[string]bool) {
	(*auths)[m.Owner.Value] = true
}

func (m *PostOperation) Validate() error {
	return nil
}

func (m *PostOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("post", (*Operation_Op6)(nil), (*PostOperation)(nil))
}

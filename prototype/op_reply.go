package prototype

func (m *ReplyOperation) GetSigner(auths *map[string]bool) {
	(*auths)[m.Owner.Value] = true
}

func (m *ReplyOperation) Validate() error {
	return nil
}

func (m *ReplyOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("reply", (*Operation_Op7)(nil), (*ReplyOperation)(nil))
}

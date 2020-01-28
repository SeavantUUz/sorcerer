package prototype

func (t *ConvertVestOperation) GetSigner(auths *map[string]bool) {
	(*auths)[t.From.Value] = true
}

func (t *ConvertVestOperation) Validate() error {
	return nil
}

func (m *ConvertVestOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("convert_vest", (*Operation_Op16)(nil), (*ConvertVestOperation)(nil))
}

package prototype

func (t *TransferOperation) GetSigner(auths *map[string]bool) {
	(*auths)[t.From.Value] = true
}

func (t *TransferOperation) Validate() error {
	return nil
}

func (m *TransferOperation) GetAffectedProps(props *map[string]bool) {
	(*props)[m.GetFrom().GetValue()] = true
	(*props)[m.GetTo().GetValue()] = true
}

func init() {
	registerOperation("transfer", (*Operation_Op2)(nil), (*TransferOperation)(nil))
}

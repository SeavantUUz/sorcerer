package prototype

func (m *TransferToVestOperation) GetSigner(auths *map[string]bool) {
	(*auths)[m.From.Value] = true
}

func (m *TransferToVestOperation) Validate() error {
	return nil
}

func (m *TransferToVestOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("transfer_to_vest", (*Operation_Op10)(nil), (*TransferToVestOperation)(nil))
}

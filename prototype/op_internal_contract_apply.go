package prototype

func (op *InternalContractApplyOperation) GetSigner(auths *map[string]bool) {
	(*auths)[op.FromCaller.Value] = true
}

func (op *InternalContractApplyOperation) Validate() error {
	return nil
}

func (m *InternalContractApplyOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

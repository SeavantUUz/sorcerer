package prototype

func (m *ContractDeployOperation) GetSigner(auths *map[string]bool) {
	(*auths)[m.Owner.Value] = true
}

func (m *ContractDeployOperation) Validate() error {
	return nil
}

func (m *ContractDeployOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

func init() {
	registerOperation("contract_deploy", (*Operation_Op13)(nil), (*ContractDeployOperation)(nil))
}

package prototype

func (m *InternalContractApplyOperation) GetAffectedProps(props *map[string]bool) {
	(*props)["*"] = true
}

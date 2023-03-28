package mocks

import "github.com/datapace/datapace/transactions"

var _ transactions.ContractLedger = (*mockContractLedger)(nil)

type mockContractLedger struct{}

// NewContractLedger returns contract ledger mock instance.
func NewContractLedger() transactions.ContractLedger {
	return mockContractLedger{}
}

func (mcl mockContractLedger) Create(...transactions.Contract) error {
	return nil
}

func (mcl mockContractLedger) Sign(transactions.Contract) error {
	return nil
}

var _ transactions.ContractRepository = (*mockContractRepository)(nil)

type mockContractRepository struct {
	mockContractLedger
}

// NewContractRepository returns contract repository mock instance.
func NewContractRepository() transactions.ContractRepository {
	return mockContractRepository{}
}

func (mcr mockContractRepository) Activate(...transactions.Contract) error {
	return nil
}

func (mcr mockContractRepository) Remove(...transactions.Contract) error {
	return nil
}

func (mcr mockContractRepository) List(string, uint64, uint64, transactions.Role) transactions.ContractPage {
	return transactions.ContractPage{}
}

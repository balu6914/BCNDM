package transactions

import "time"

// ContractLedger contains contract specific writer API definition.
type ContractLedger interface {
	// Creates contract and stores it into datastore.
	Create(...Contract) error

	// Signs contract.
	Sign(Contract) error
}

// ContractRepository defines API for contract persistance (read and write).
type ContractRepository interface {
	ContractLedger

	// Activates contract.
	Activate(...Contract) error

	// Removes contract.
	Remove(...Contract) error

	// List finds and returns page of contracts by contract owner or partner,
	// depending on role.
	List(string, uint64, uint64, Role) ContractPage
}

// ContractPage represents contract page with page data and metadata.
type ContractPage struct {
	Page      uint64
	Limit     uint64
	Total     uint64
	Contracts []Contract
}

// Contract contains contract data.
type Contract struct {
	StreamID   string
	StreamName string
	StartTime  time.Time
	EndTime    time.Time
	OwnerID    string
	PartnerID  string
	Share      uint64
	Signed     bool
}

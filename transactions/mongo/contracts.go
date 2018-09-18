package mongo

import (
	"monetasa/transactions"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ transactions.ContractRepository = (*contractRepository)(nil)

type contractRepository struct {
	db *mgo.Session
}

// NewContractRepository returns mongoDB specific instance of contract repository.
func NewContractRepository(db *mgo.Session) transactions.ContractRepository {
	return contractRepository{db: db}
}

func (cr contractRepository) Create(contracts ...transactions.Contract) error {
	mongoContracts := []interface{}{}
	for _, contract := range contracts {
		mc := toMongoContract(contract)
		mongoContracts = append(mongoContracts, mc)
	}

	session := cr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(contractsCollection)

	query := bson.M{
		"stream_id": contracts[0].StreamID,
		"end_time": bson.M{
			"$gt": time.Now(),
		},
		"active": true,
	}
	mcs := []mongoContract{}
	if count, err := collection.Find(query).Count(); err != nil {
		return err
	}
	if count != 0 {
		return transactions.ErrConflict
	}

	if err := collection.Insert(mongoContracts...); err != nil {
		if mgo.IsDup(err) {
			return transactions.ErrConflict
		}

		return err
	}

	return nil
}

func (cr contractRepository) Sign(contract transactions.Contract) error {
	mc := toMongoContract(contract)
	mc.Active = true

	session := cr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(contractsCollection)

	update := bson.M{
		"$set": mongoContract{
			Signed: true,
			Active: true,
		},
	}
	if err := collection.Update(mc, update); err != nil {
		if err == mgo.ErrNotFound {
			return transactions.ErrNotFound
		}
		if mgo.IsDup(err) {
			return transactions.ErrConflict
		}

		return err
	}

	return nil
}

func (cr contractRepository) Activate(contracts ...transactions.Contract) error {
	for _, contract := range contracts {
		if err := cr.setActive(contract, true); err != nil {
			return err
		}
	}

	return nil
}

func (cr contractRepository) Remove(contracts ...transactions.Contract) error {
	for _, contract := range contracts {
		if err := cr.setActive(contract, false); err != nil {
			return err
		}
	}

	return nil
}

func (cr contractRepository) List(userID string, pageNo, limit uint64, role transactions.Role) transactions.ContractPage {
	session := cr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(contractsCollection)

	query := bson.M{"active": true}
	switch role {
	case transactions.Owner:
		query["owner_id"] = userID
	case transactions.Partner:
		query["partner_id"] = userID
	case transactions.AllRoles:
		query["$or"] = []bson.M{
			{"owner_id": userID},
			{"partner_id": userID},
		}
	}

	count, err := collection.Find(query).Count()
	if err != nil {
		return transactions.ContractPage{}
	}
	total := uint64(count)

	offset := pageNo * limit
	if offset > total {
		return transactions.ContractPage{}
	}

	mcs := []mongoContract{}
	if err := collection.Find(query).Skip(int(offset)).Limit(int(limit)).All(&mcs); err != nil {
		return transactions.ContractPage{}
	}

	contracts := []transactions.Contract{}
	for _, mc := range mcs {
		contracts = append(contracts, toContract(mc))
	}

	return transactions.ContractPage{
		Page:      pageNo,
		Limit:     limit,
		Total:     total,
		Contracts: contracts,
	}
}

func (cr contractRepository) setActive(contract transactions.Contract, active bool) error {
	mc := toMongoContract(contract)
	mc.Active = !active

	session := cr.db.Copy()
	defer session.Close()
	collection := session.DB(dbName).C(contractsCollection)

	update := bson.M{
		"$set": mongoContract{
			Active: active,
		},
	}
	if err := collection.Update(mc, update); err != nil {
		if err == mgo.ErrNotFound {
			return transactions.ErrNotFound
		}
		if mgo.IsDup(err) {
			return transactions.ErrConflict
		}

		return err
	}

	return nil
}

type mongoContract struct {
	StreamID  string    `bson:"stream_id,omitempty"`
	StartTime time.Time `bson:"start_time,omitempty"`
	EndTime   time.Time `bson:"end_time,omitempty"`
	OwnerID   string    `bson:"owner_id,omitempty"`
	PartnerID string    `bson:"partner_id,omitempty"`
	Share     uint64    `bson:"share,omitempty"`
	Signed    bool      `bson:"signed,omitempty"`
	Active    bool      `bson:"active"`
}

func toMongoContract(contract transactions.Contract) mongoContract {
	return mongoContract{
		StreamID:  contract.StreamID,
		StartTime: contract.StartTime,
		EndTime:   contract.EndTime,
		OwnerID:   contract.OwnerID,
		PartnerID: contract.PartnerID,
		Share:     contract.Share,
		Signed:    contract.Signed,
	}
}

func toContract(contract mongoContract) transactions.Contract {
	return transactions.Contract{
		StreamID:  contract.StreamID,
		StartTime: contract.StartTime,
		EndTime:   contract.EndTime,
		OwnerID:   contract.OwnerID,
		PartnerID: contract.PartnerID,
		Share:     contract.Share,
		Signed:    contract.Signed,
	}
}

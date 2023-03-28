package mongo

import (
	trs "github.com/datapace/datapace/terms"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ trs.TermsRepository = (*termsRepository)(nil)

type termsRepository struct {
	db *mgo.Session
}

// NewTermsRepository returns new Terms repository.
func NewTermsRepository(db *mgo.Session) trs.TermsRepository {
	return &termsRepository{db: db}
}

func (tr termsRepository) Save(t trs.Terms) (string, error) {
	s := tr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	t.ID = bson.NewObjectId()
	dbSub := toDBTerms(t)

	if err := c.Insert(dbSub); err != nil {
		if mgo.IsDup(err) {
			return "", trs.ErrConflict
		}

		return "", err
	}

	return dbSub.ID.Hex(), nil
}

func (tr termsRepository) One(id string) (trs.Terms, error) {
	s := tr.db.Copy()
	defer s.Close()
	c := s.DB(dbName).C(collectionName)

	var dbSub terms
	query := bson.M{
		"_id":    bson.ObjectIdHex(id),
		"active": true,
	}
	if err := c.Find(query).One(&dbSub); err != nil {
		if err == mgo.ErrNotFound {
			return trs.Terms{}, trs.ErrNotFound
		}

		return trs.Terms{}, err
	}

	sub := toTerms(dbSub)

	return sub, nil
}

// Terms is terms representation in DB.
type terms struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	StreamID  string        `bson:"stream_id,omitempty"`
	TermsURL  string        `bson:"terms_url,omitempty"`
	TermsHash string        `bson:"terms_hash,omitempty"`
}

func toDBTerms(t trs.Terms) terms {
	return terms{
		ID:        t.ID,
		StreamID:  t.StreamID,
		TermsURL:  t.TermsURL,
		TermsHash: t.TermsHash,
	}
}

func toTerms(dbTerms terms) trs.Terms {
	return trs.Terms{
		ID:        dbTerms.ID,
		StreamID:  dbTerms.StreamID,
		TermsURL:  dbTerms.TermsURL,
		TermsHash: dbTerms.TermsHash,
	}
}

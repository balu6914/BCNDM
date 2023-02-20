package mocks

import (
	"github.com/datapace/datapace/streams"
	"gopkg.in/mgo.v2/bson"
)

var _ streams.CategoryRepository = (*categoryRepositoryMock)(nil)

type categoryRepositoryMock struct{}

func NewCategoryRepository() streams.CategoryRepository {
	return &categoryRepositoryMock{}
}

func (crm *categoryRepositoryMock) Save(categoryName string, subCategoryName []string) (string, error) {
	id := bson.NewObjectId().Hex()
	return id, nil
}

func (crm *categoryRepositoryMock) List(key string) ([]streams.Category, error) {
	list := []streams.Category{}
	return list, nil
}

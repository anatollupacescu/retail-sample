package warehouse_test

import (
	"errors"
	"testing"

	itemTypeMocks "github.com/anatollupacescu/retail-sample/pkg/retail-sample/itemtype/mocks"
	"github.com/anatollupacescu/retail-sample/pkg/retail-sample/warehouse/mocks"

	"github.com/anatollupacescu/retail-sample/pkg/retail-sample/warehouse"

	qt "github.com/frankban/quicktest"
	"github.com/golang/mock/gomock"
)

/*
testing
Adding:
	type doesn't exist ->
		no side effects, returns error
	item type found
		item present in the store ->
			update successful -> ok
			update fails -> return
		item not present ->
			add item

*/
func TestWarehouse(t *testing.T) {

	t.Run("should reject non existent item type ids", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := itemTypeMocks.NewMockRepository(mockCtrl)
		// type not found
		itemTypeRepository.EXPECT().Get(uint64(1)).Return("")

		itemStore := mocks.NewMockStore(mockCtrl)
		// will not try to add item
		itemStore.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

		w := warehouse.Repository{
			ItemStore:          itemStore,
			ItemTypeRepository: itemTypeRepository,
		}

		err := w.Add(1, 23)

		c.Assert(err, qt.Equals, warehouse.ErrItemTypeNotFound)
	})

	t.Run("should update successfully when type exists and item already in store", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := itemTypeMocks.NewMockRepository(mockCtrl)
		itemTypeRepository.EXPECT().Get(uint64(1)).Return("butter")

		itemStore := mocks.NewMockStore(mockCtrl)
		itemStore.EXPECT().Get(uint64(1)).Return(9, nil)
		itemStore.EXPECT().Update(uint64(1), gomock.Eq(10)).Return(nil)

		itemStore.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

		w := warehouse.Repository{
			ItemStore:          itemStore,
			ItemTypeRepository: itemTypeRepository,
		}

		err := w.Add(1, 10)

		c.Assert(err, qt.IsNil)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := itemTypeMocks.NewMockRepository(mockCtrl)
		itemTypeRepository.EXPECT().Get(uint64(1)).Return("butter")

		itemStore := mocks.NewMockStore(mockCtrl)
		itemStore.EXPECT().Get(uint64(1)).Return(0, nil)
		itemStore.EXPECT().Update(uint64(1), gomock.Eq(10)).Return(errors.New("update no go"))

		itemStore.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

		w := warehouse.Repository{
			ItemStore:          itemStore,
			ItemTypeRepository: itemTypeRepository,
		}

		err := w.Add(1, 10)

		c.Assert(err, qt.Equals, warehouse.ErrUpdate)
	})

	t.Run("should add item if type exists and item is not in store", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := itemTypeMocks.NewMockRepository(mockCtrl)
		// type exists
		itemTypeRepository.EXPECT().Get(uint64(1)).Return("butter")

		itemStore := mocks.NewMockStore(mockCtrl)
		itemStore.EXPECT().Get(uint64(1)).Return(0, warehouse.ErrItemNotFound)
		itemStore.EXPECT().Add(uint64(1), 10).Times(1)

		itemStore.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		w := warehouse.Repository{
			ItemStore:          itemStore,
			ItemTypeRepository: itemTypeRepository,
		}

		err := w.Add(1, 10)

		c.Assert(err, qt.IsNil)
	})
}

package warehouse_test

import (
	"errors"
	"testing"

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
Quantity:
	0 when empty
	sum when has items
*/
func TestWarehouse(t *testing.T) {

	t.Run("can not Add a type that does not exists", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := mocks.NewMockItemTypeRepository(mockCtrl)
		// type not found
		itemTypeRepository.EXPECT().Get(uint64(1)).Return("")

		itemRepository := mocks.NewMockItemRepository(mockCtrl)
		// will not try to add item
		itemRepository.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

		w := warehouse.Repository{
			ItemRepository:     itemRepository,
			ItemTypeRepository: itemTypeRepository,
		}

		err := w.Add(1, 23)

		c.Assert(err, qt.Equals, warehouse.ErrItemTypeNotFound)
	})

	t.Run("should update successfully when type exists and item already in store", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := mocks.NewMockItemTypeRepository(mockCtrl)
		itemTypeRepository.EXPECT().Get(uint64(1)).Return("butter")

		itemRepository := mocks.NewMockItemRepository(mockCtrl)
		itemRepository.EXPECT().Get(uint64(1)).Return(9, nil)
		itemRepository.EXPECT().Update(uint64(1), gomock.Eq(10)).Return(nil)

		itemRepository.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

		w := warehouse.Repository{
			ItemRepository:     itemRepository,
			ItemTypeRepository: itemTypeRepository,
		}

		err := w.Add(1, 10)

		c.Assert(err, qt.IsNil)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := mocks.NewMockItemTypeRepository(mockCtrl)
		itemTypeRepository.EXPECT().Get(uint64(1)).Return("butter")

		itemRepository := mocks.NewMockItemRepository(mockCtrl)
		itemRepository.EXPECT().Get(uint64(1)).Return(0, nil)
		itemRepository.EXPECT().Update(uint64(1), gomock.Eq(10)).Return(errors.New("update no go"))

		itemRepository.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

		w := warehouse.Repository{
			ItemRepository:     itemRepository,
			ItemTypeRepository: itemTypeRepository,
		}

		err := w.Add(1, 10)

		c.Assert(err, qt.Equals, warehouse.ErrUpdate)
	})

	t.Run("should add item if type exists and item is not in store", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := mocks.NewMockItemTypeRepository(mockCtrl)
		// type exists
		itemTypeRepository.EXPECT().Get(uint64(1)).Return("butter")

		itemRepository := mocks.NewMockItemRepository(mockCtrl)
		itemRepository.EXPECT().Get(uint64(1)).Return(0, warehouse.ErrItemNotFound)
		itemRepository.EXPECT().Add(uint64(1), 10).Times(1)

		itemRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		w := warehouse.Repository{
			ItemRepository:     itemRepository,
			ItemTypeRepository: itemTypeRepository,
		}

		err := w.Add(1, 10)

		c.Assert(err, qt.IsNil)
	})

	t.Run("quantity should be zero in empty store for existing type", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := mocks.NewMockItemTypeRepository(mockCtrl)
		itemRepository := mocks.NewMockItemRepository(mockCtrl)

		itemRepository.EXPECT().Get(uint64(1)).Return(99, nil)

		w := warehouse.Repository{
			ItemRepository:     itemRepository,
			ItemTypeRepository: itemTypeRepository,
		}

		qty, err := w.Get(uint64(1))

		c.Assert(err, qt.IsNil)
		c.Assert(qty, qt.Equals, 99)
	})

	t.Run("should return error when requesting quantity for non-existing type", func(t *testing.T) {
		c := qt.New(t)
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		itemTypeRepository := mocks.NewMockItemTypeRepository(mockCtrl)
		itemRepository := mocks.NewMockItemRepository(mockCtrl)

		itemRepository.EXPECT().Get(uint64(1)).Return(0, warehouse.ErrItemNotFound)

		w := warehouse.Repository{
			ItemRepository:     itemRepository,
			ItemTypeRepository: itemTypeRepository,
		}

		_, err := w.Get(uint64(1))

		c.Assert(err, qt.Equals, warehouse.ErrItemNotFound)
	})

}

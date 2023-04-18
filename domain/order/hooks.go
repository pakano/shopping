package order

import (
	"shopping/domain/cart"
	"shopping/domain/product"

	"gorm.io/gorm"
)

// 创建之前，查找购物车并删除
func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	// 查找当前用户的购物车
	var currentCart cart.Cart
	if err := tx.Where("UserID = ?", order.UserID).First(&currentCart).Error; err != nil {
		return err
	}

	// 删除购物车中的所有商品项
	if err := tx.Where("CartID = ?", currentCart.ID).Unscoped().Delete(&cart.Item{}).Error; err != nil {
		return err
	}

	// 删除当前购物车
	if err := tx.Unscoped().Delete(&currentCart).Error; err != nil {
		return err
	}

	return nil
}

// 保存之前，更新产品库存
func (orderedItem *OrderedItem) BeforeSave(tx *gorm.DB) (err error) {
	// 查找当前订单项的商品
	var currentProduct product.Product
	var currentOrderedItem OrderedItem
	if err := tx.Where("ID = ?", orderedItem.ProductID).First(&currentProduct).Error; err != nil {
		return err
	}

	// 查找当前订单项的数量
	reservedStockCount := 0
	if err := tx.Where("ID = ?", orderedItem.ID).First(&currentOrderedItem).Error; err == nil {
		reservedStockCount = currentOrderedItem.Count
	}

	// 计算新的库存数量
	newStockCount := currentProduct.StockCount + reservedStockCount - orderedItem.Count
	if newStockCount < 0 {
		return ErrNotEnoughStock
	}

	// 更新产品库存
	if err := tx.Model(&currentProduct).Update("StockCount", newStockCount).Error; err != nil {
		return err
	}

	// 如果订单项的数量为0，则从数据库中删除该订单项
	if orderedItem.Count == 0 {
		err := tx.Unscoped().Delete(currentOrderedItem).Error
		return err
	}

	return
}

// 如果订单被取消，金额将返回产品库存
func (order *Order) BeforeUpdate(tx *gorm.DB) (err error) {
	// 如果订单被取消，则将金额返回到产品库存中
	if order.IsCanceled {
		var orderedItems []OrderedItem
		if err := tx.Where("OrderID = ?", order.ID).Find(&orderedItems).Error; err != nil {
			return err
		}
		for _, item := range orderedItems {
			var currentProduct product.Product
			if err := tx.Where("ID = ?", item.ProductID).First(&currentProduct).Error; err != nil {
				return err
			}
			newStockCount := currentProduct.StockCount + item.Count
			if err := tx.Model(&currentProduct).Update(
				"StockCount", newStockCount).Error; err != nil {
				return err
			}
			if err := tx.Model(&item).Update(
				"IsCanceled", true).Error; err != nil {
				return err
			}
		}
	}

	return
}

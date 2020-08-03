package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	appApiHelper "github.com/xifengzhu/eshop/routers/app_api/api_helpers"
)

type CartItemParams struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
	GoodsID  int `json:"goods_id" validate:"required"`
}

type CartItemIDParams struct {
	ItemIDs []int `json:"item_ids" validate:"required,gt=0"`
}

type CartItemQtyParams struct {
	ItemID   int `json:"item_id" validate:"required"`
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

// @Summary 获取购物车列表
// @Produce  json
// @Tags 购物车
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart/my [get]
// @Security ApiKeyAuth
func GetCartItems(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)

	var items []models.CarItem
	parmMap := map[string]interface{}{"user_id": user.ID}

	models.Where(Query{Conditions: parmMap, Preloads: []string{"Goods"}}).Find(&items)

	apiHelpers.ResponseSuccess(c, items)

}

// @Summary 添加商品到获取购物车
// @Produce  json
// @Tags 购物车
// @Param params body CartItemParams true "shopping_cart"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart/add [post]
// @Security ApiKeyAuth
func AddCartItem(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)

	// check params
	var cartItem CartItemParams
	if err := apiHelpers.ValidateParams(c, &cartItem, "json"); err != nil {
		return
	}

	// check goods exist?
	var goods models.Goods
	goods.ID = cartItem.GoodsID
	exist := goods.IsExist()
	if !exist {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "商品不存在或者已下架")
		return
	}

	// 如果用户的购物车中没有该商品，则加进去； 否则购物车数量相加
	item, err := user.FindShoppingCartItemByGoodsID(cartItem.GoodsID)
	if err != nil {
		copier.Copy(&item, &cartItem)
		item.UserID = user.ID
		err := models.Create(&item)
		if err != nil {
			apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
			return
		}
	} else {
		item.Quantity += cartItem.Quantity
		models.Save(&item)
	}

	apiHelpers.ResponseOK(c)
}

// @Summary 勾选购物车项
// @Produce  json
// @Tags 购物车
// @Param params body CartItemIDParams true "items_ids"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart/check [put]
// @Security ApiKeyAuth
func CheckCartItem(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	// check params
	var itemID CartItemIDParams
	if err := apiHelpers.ValidateParams(c, &itemID, "json"); err != nil {
		return
	}

	items, _ := user.GetShoppingCartItemByIDs(itemID.ItemIDs)
	changedAttrs := map[string]interface{}{"checked": true}
	for _, item := range items {
		models.Update(&item, changedAttrs)
	}
	apiHelpers.ResponseOK(c)
}

// @Summary 取消购物车项
// @Produce  json
// @Tags 购物车
// @Param params body CartItemIDParams true "items_ids"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart/uncheck [put]
// @Security ApiKeyAuth
func UnCheckCartItem(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	var itemID CartItemIDParams
	if err := apiHelpers.ValidateParams(c, &itemID, "json"); err != nil {
		return
	}

	items, _ := user.GetShoppingCartItemByIDs(itemID.ItemIDs)
	changedAttrs := map[string]interface{}{"checked": false}
	for _, item := range items {
		models.Update(&item, changedAttrs)
	}
	apiHelpers.ResponseOK(c)
}

// @Summary 删除商品
// @Produce  json
// @Tags 购物车
// @Param params body CartItemIDParams true "items_ids"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart/delete [delete]
// @Security ApiKeyAuth
func DeleteCartItem(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	// check params
	var itemID CartItemIDParams
	if err := apiHelpers.ValidateParams(c, &itemID, "json"); err != nil {
		return
	}

	items, _ := user.GetShoppingCartItemByIDs(itemID.ItemIDs)
	var ids []int
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	err := models.DestroyAll(&[]models.CarItem{}, Query{Conditions: ids})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}
	apiHelpers.ResponseOK(c)

}

// @Summary 更新商品数量
// @Produce  json
// @Tags 购物车
// @Param params body CartItemQtyParams true "商品id与数量"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart/qty [put]
// @Security ApiKeyAuth
func UpdateCartItemQty(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)

	var itemParams CartItemQtyParams
	if err := apiHelpers.ValidateParams(c, &itemParams, "json"); err != nil {
		return
	}

	item, err := user.GetShoppingCartItemByID(itemParams.ItemID)
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	changedAttrs := models.CarItem{Quantity: itemParams.Quantity}

	err = models.Update(&item, &changedAttrs)
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}
	apiHelpers.ResponseOK(c)
}

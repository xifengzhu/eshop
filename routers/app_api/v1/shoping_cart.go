package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	appApiHelper "github.com/xifengzhu/eshop/routers/app_api/api_helpers"
)

type CartItemParams struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
	GoodsID  int `json:"goods_id" binding:"required"`
}

type CartItemIDParams struct {
	ItemIDs []int `json:"item_id" binding:"required,gt=0"`
}

type CartItemQtyParams struct {
	ItemID   int `json:"item_id" binding:"required"`
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

// @Summary 获取购物车列表
// @Produce  json
// @Tags 购物车
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart [get]
// @Security ApiKeyAuth
func GetCartItems(c *gin.Context) {

	user := appApiHelper.CurrentUser(c)

	pagination := apiHelpers.SetDefaultPagination(c)

	var items []models.CarItem
	parmMap := map[string]interface{}{"user_id": user.ID}

	models.AllResource(&items, Query{Conditions: parmMap, Preloads: []string{"Goods"}})

	response := apiHelpers.Collection{Pagination: pagination, List: items}

	apiHelpers.ResponseSuccess(c, response)

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
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	// check goods exist?
	var goods models.Goods
	goods.ID = cartItem.GoodsID
	exist := goods.IsExist()
	if !exist {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("商品不存在或者已下架"))
		return
	}

	// 如果用户的购物车中没有该商品，则加进去； 否则购物车数量相加
	item, err := user.FindShoppingCartItemByGoodsID(cartItem.GoodsID)
	if err != nil {
		copier.Copy(&item, &cartItem)
		item.UserID = user.ID
		err := models.CreateResource(&item)
		if err != nil {
			apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
			return
		}
	} else {
		item.Quantity += cartItem.Quantity
		models.SaveResource(&item)
	}

	apiHelpers.ResponseOK(c)
}

// @Summary 勾选购物车项
// @Produce  json
// @Tags 购物车
// @Param item_ids body CartItemIDParams true "items_ids"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart/check [put]
// @Security ApiKeyAuth
func CheckCartItem(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	// check params
	var itemID CartItemIDParams
	if err := c.ShouldBindJSON(&itemID); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	items, _ := user.GetShoppingCartItemByIDs(itemID.ItemIDs)
	for _, item := range items {
		item.Checked = true
		models.SaveResource(&item)
	}
	apiHelpers.ResponseOK(c)
}

// @Summary 取消购物车项
// @Produce  json
// @Tags 购物车
// @Param item_ids body CartItemIDParams true "items_ids"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart/uncheck [put]
// @Security ApiKeyAuth
func UnCheckCartItem(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	// check params
	var itemID CartItemIDParams
	if err := c.ShouldBindJSON(&itemID); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	items, _ := user.GetShoppingCartItemByIDs(itemID.ItemIDs)
	for _, item := range items {
		item.Checked = false
		models.SaveResource(&item)
	}
	apiHelpers.ResponseOK(c)
}

// @Summary 删除商品
// @Produce  json
// @Tags 购物车
// @Param item_ids body CartItemIDParams true "items_ids"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/shopping_cart/remove [delete]
// @Security ApiKeyAuth
func DeleteCartItem(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	// check params
	var itemID CartItemIDParams
	if err := c.ShouldBindJSON(&itemID); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	items, _ := user.GetShoppingCartItemByIDs(itemID.ItemIDs)
	var ids []int
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	err := models.DestroyAll(&[]models.CarItem{}, Query{Conditions: ids})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
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
	if err := c.ShouldBindJSON(&itemParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	item, err := user.GetShoppingCartItemByID(itemParams.ItemID)
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	item.Quantity = itemParams.Quantity
	err = models.SaveResource(&item)
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}
	apiHelpers.ResponseOK(c)
}

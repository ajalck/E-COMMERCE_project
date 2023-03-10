package handler

import (
	"ajalck/e_commerce/domain"
	services "ajalck/e_commerce/usecase/interface"
	"ajalck/e_commerce/utils"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(userUseCase services.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var newUser domain.User
	if err := c.Bind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uh.userUseCase.CreateUser(c, newUser)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusFound, gin.H{"message": "New user created successfully"})
	}
}
func (uh *UserHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	users, metaData, err := uh.userUseCase.ListProducts(page, perPage)
	type Page struct {
		users    []domain.ProductResponse
		metaData utils.MetaData
	}
	result := Page{
		users:    users,
		metaData: metaData,
	}
	if err == nil {
		c.JSON(200, result.users)
		c.JSON(http.StatusFound, result.metaData)
		return
	}
	if err != nil {
		c.JSON(400, err.Error())
	}
}

//Wishlist

func (uh *UserHandler) AddWishlist(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	err := uh.userUseCase.AddWishlist(user_id, product_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't add new item to wishlist", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("New item added to wishlist", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}
func (uh *UserHandler) ViewWishList(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	wishList, metaData, err := uh.userUseCase.ViewWishList(user_id, page, perPage)
	var results struct {
		WishLists []domain.WishListResponse
		MetaData  utils.MetaData
	}
	results.WishLists = wishList
	results.MetaData = metaData
	if err != nil {
		response := utils.ErrorResponse("couldn't reach to your wishlist", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is your wishlist...", results)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}
func (uh *UserHandler) DeleteWishList(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	err := uh.userUseCase.DeleteWishList(user_id, product_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't remove the item from your wishlist", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("one item removed successfully from your wishlist", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

//Cart

func (uh *UserHandler) AddCart(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	err := uh.userUseCase.AddCart(user_id, product_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't add new item to cart", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("New item added to cart", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}
func (uh *UserHandler) ViewCart(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	cart, metaData, err := uh.userUseCase.ViewCart(user_id, page, perPage)
	results := struct {
		cart     []domain.CartResponse
		MetaData utils.MetaData
	}{
		cart:     cart,
		MetaData: metaData,
	}
	if err != nil {
		response := utils.ErrorResponse("couldn't reach to your cart", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is your cart...", results.cart)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
	c.JSON(200, results.MetaData)
}
func (uh *UserHandler) DeleteCart(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	err := uh.userUseCase.DeleteCart(user_id, product_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't remove an item from your cart", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("One item removed successfully from your cart", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

//Shipping

func (uh *UserHandler) AddShippingDetails(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	var newAddress domain.ShippingDetails
	if err := c.Bind(&newAddress); err != nil {
		response := utils.ErrorResponse("Invalid inputs", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	err := uh.userUseCase.AddShippingDetails(user_id, newAddress)
	if err != nil {
		response := utils.ErrorResponse("Failed to add new shipping details", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusConflict)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("New shipping details added successfully", nil)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)
}
func (uh *UserHandler) ListShippingDetails(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	shippingDetails, err := uh.userUseCase.ListShippingDetails(user_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't list shipping details", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is your shipping details", shippingDetails)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)
}

func (uh *UserHandler) DeleteShippingDetails(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	address_id, _ := strconv.Atoi(c.Query("address_id"))
	err := uh.userUseCase.DeleteShippingDetails(user_id, address_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't delete shipping details", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("successfully removed selected shipping details", nil)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)
}
func (uh *UserHandler) PlaceOrder(c *gin.Context) {
<<<<<<< HEAD
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	address_id, _ := strconv.Atoi(c.Query("address_id"))
	coupon_id, _ := strconv.Atoi(c.Query("coupon_id"))

	err := uh.userUseCase.PlaceOrder(c,user_id, product_id, address_id, coupon_id)
	if err != nil {
		response := utils.ErrorResponse("Order placement failed !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Order placed successfully", nil)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)
=======
	// user_id,_:=strconv.Atoi(c.Writer.Header().Get("id"))
	// product_id,_:=strconv.Atoi(c.Query("product_id"))
	// address_id,_:=strconv.Atoi(c.Query("address_id"))
>>>>>>> 1df3e9a651c7205727c237e280b2685ffd593ddb

}

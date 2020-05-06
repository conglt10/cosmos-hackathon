package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MinNamePrice is Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// MinProductPrice is Initial Starting Price for a product that was never previously owned
var MinProductPrice = sdk.Coins{sdk.NewInt64Coin("producttoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value       string         `json:"value"`
	Owner       sdk.AccAddress `json:"owner"`
	Price       sdk.Coins      `json:"price"`
	Description string         `json:"description"`
}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s
Description: %s`, w.Owner, w.Value, w.Price, w.Description))
}

// Product is struct that contains all the metadata of a product
type Product struct {
	ProductID   string         `json:"productID"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Owner       sdk.AccAddress `json:"owner"`
}

//NewProduct returns a new product
func NewProduct() Product {
	return Product{}
}

// implement fmt.Stringer
func (product Product) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
	ProductID: %s
	Title: %s
	Description: %s
	Owner: %s`, product.ProductID, product.Title, product.Description, product.Owner))
}

// Sell is a struct contains all the metadata of a sell
type Sell struct {
	SellID    string         `json:"sellID"`
	ProductID string         `json:"productID"`
	Seller    sdk.AccAddress `json:"seller"`
	MinPrice  sdk.Coins      `json:"minPrice"`
}

//NewSell returns a new sell
func NewSell() Sell {
	return Sell{
		MinPrice: MinProductPrice,
	}
}

// implement fmt.Stringer
func (sell Sell) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
	SellID: %s
	ProductID: %s
	Seller: %s
	MinPrice: %s`, sell.SellID, sell.ProductID, sell.Seller, sell.MinPrice))
}

// Reservation is a struct contains all the metadata of a reservation
type Reservation struct {
	ReservationID string         `json:"reservationID"`
	ProductID     string         `json:"productID"`
	Buyer         sdk.AccAddress `json:"buyer"`
	Price         sdk.Coins      `json:"price"`
}

//NewReservation returns a new Reservation
func NewReservation() Reservation {
	return Reservation{}
}

// implement fmt.Stringer
func (reservation Reservation) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
	ReservationID: %s
	ProductID: %s
	Buyer: %s
	Price: %s`, reservation.ReservationID, reservation.ProductID, reservation.Buyer, reservation.Price))
}

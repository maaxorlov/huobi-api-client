package huobiapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func (A *HuobiApi) getSpotAcc() (Accounts, error) {
	A.debugData.lastRequestMethod = "getSpotAcc"

	uri := A.buildPrivateURI("GET", "/v1/account/accounts")

	response, err := A.doRequestGET(uri)
	if err != nil {
		return Accounts{}, fmt.Errorf("error in getSpotAcc with response, err := A.doRequestGET(uri): %s", err.Error())
	}

	var accounts Accounts
	err = json.Unmarshal(response, &accounts)
	if err != nil {
		return Accounts{}, fmt.Errorf("error in getSpotAcc with err = json.Unmarshal(response, &accounts): %s", err.Error())
	}

	return accounts, nil
}

func (A *HuobiApi) getBalancesAcc() (Balance, error) {
	A.debugData.lastRequestMethod = "getBalancesAcc"

	uri := A.buildPrivateURI("GET", "/v1/account/accounts/"+A.spotAccID+"/balance")

	response, err := A.doRequestGET(uri)
	if err != nil {
		return Balance{}, fmt.Errorf("error in getBalancesAcc with response, err := A.doRequestGET(uri): %s", err.Error())
	}

	var balance Balance
	err = json.Unmarshal(response, &balance)
	if err != nil {
		return Balance{}, fmt.Errorf("error in getBalancesAcc with err = json.Unmarshal(response, &balance): %s", err.Error())
	}

	if balance.Status != "ok" {
		return Balance{}, errors.New(balance.Message)
	}

	return balance, nil
}

func (A *HuobiApi) getOrderBookOnPair(symbol string) (OrderBook, error) {
	A.debugData.lastRequestMethod = "getOrderBookOnPair"

	additionalURLParams := map[string]string{
		"symbol": symbol,
		"type":   "step0",
	}
	uri := A.buildPublicURI("/market/depth", additionalURLParams)

	response, err := A.doRequestGET(uri)
	if err != nil {
		return OrderBook{}, fmt.Errorf("error in getOrderBookOnPair with response, err := A.doRequestGET(uri): %s", err.Error())
	}

	var orderBook OrderBook
	err = json.Unmarshal(response, &orderBook)
	if err != nil {
		return OrderBook{}, fmt.Errorf("error in getOrderBookOnPair with err = json.Unmarshal(response, &orderBook): %s", err.Error())
	}

	if orderBook.Status != "ok" {
		return OrderBook{}, errors.New(orderBook.Message)
	}

	return orderBook, nil
}

func (A *HuobiApi) placeOrder(side, symbol string, amount, price float64) (string, error) {
	A.debugData.lastRequestMethod = "placeOrder"

	_price := strconv.FormatFloat(price, 'f', -1, 64)
	_amount := strconv.FormatFloat(amount, 'f', -1, 64)
	uri := A.buildPrivateURI("POST", "/v1/order/orders/place")
	request := PlaceOrderRequest{
		AccountId: A.spotAccID,
		Type:      side + "-limit",
		Source:    "spot-api",
		Symbol:    symbol,
		Price:     _price,
		Amount:    _amount,
	}
	jsonBytes, _ := json.Marshal(request)

	response, err := A.doRequestPOST(uri, jsonBytes)
	if err != nil {
		return "", fmt.Errorf("error in placeOrder with response, err := A.doRequestPOST(uri, jsonBytes): %s", err.Error())
	}

	var placeOrder PlaceOrder
	err = json.Unmarshal(response, &placeOrder)
	if err != nil {
		return "", fmt.Errorf("error in placeOrder with err = json.Unmarshal(body, &placeOrder): %s", err.Error())
	}

	if placeOrder.Status != "ok" {
		return "", errors.New(placeOrder.Message)
	}

	return placeOrder.Data, nil
}

func (A *HuobiApi) cancelOrder(id string) error {
	A.debugData.lastRequestMethod = "cancelOrder"

	uri := A.buildPrivateURI("POST", "/v1/order/orders/"+id+"/submitcancel")

	response, err := A.doRequestPOST(uri, nil)
	if err != nil {
		return fmt.Errorf("error in cancelOrder with response, err := A.doRequestPOST(uri, nil): %s", err.Error())
	}

	var cancelOrder CancelOrder
	err = json.Unmarshal(response, &cancelOrder)
	if err != nil {
		return fmt.Errorf("error in cancelOrder with err = json.Unmarshal(response, &cancelOrder): %s", err.Error())
	}

	if cancelOrder.Status != "ok" {
		return errors.New(cancelOrder.Message)
	}

	return nil
}

func (A *HuobiApi) getOrderStatus(id string) (Order, error) {
	A.debugData.lastRequestMethod = "getOrderStatus"

	uri := A.buildPrivateURI("GET", "/v1/order/orders/"+id)

	response, err := A.doRequestGET(uri)
	if err != nil {
		return Order{}, fmt.Errorf("error in getOrderStatus with response, err := A.doRequestGET(uri): %s", err.Error())
	}

	var order Order
	err = json.Unmarshal(response, &order)
	if err != nil {
		return Order{}, fmt.Errorf("error in getOrderStatus with err = json.Unmarshal(response, &order): %s", err.Error())
	}

	if order.Status != "ok" {
		return Order{}, errors.New(order.Message)
	}

	return order, nil
}

func (A *HuobiApi) getMyOpenOrders(symbol string) (Orders, error) {
	A.debugData.lastRequestMethod = "getMyOpenOrder"

	additionalURLParams := map[string]string{"symbol": symbol}
	uri := A.buildPrivateURI("GET", "/v1/order/openOrders", additionalURLParams)

	response, err := A.doRequestGET(uri)
	if err != nil {
		return Orders{}, fmt.Errorf("error in getMyOpenOrders with response, err := A.doRequestGET(uri): %s", err.Error())
	}

	var openOrders Orders
	err = json.Unmarshal(response, &openOrders)
	if err != nil {
		return Orders{}, fmt.Errorf("error in getMyOpenOrders with err = json.Unmarshal(body, &openOrders): %s", err.Error())
	}

	if openOrders.Status != "ok" {
		return Orders{}, errors.New(openOrders.Message)
	}

	return openOrders, nil
}

package huobiapi

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

//  + Init(accountID string, apiKey string, apiSecret string) error
//  + GetBalances() (*map[string]Balance, error)
//  + GetOrderBook(symbol string) (*OrderBook, error)
//  + Buy(symbol string, amount float64, price float64) (*MakedOrder, error)
//  + Sell(symbol string, amount float64, price float64) (*MakedOrder, error)
//  + CancelOrder(symbol string, id string) error
//  + GetOrderStatus(id string, symbol string) (*MakedOrder, error)
//  + GetMyOpenOrders(symbol string) (*[]MakedOrder, error)

func TstInit(a *HuobiApi, apiPassphrase, apiKey, apiSecret string) {
	err := a.Init(apiPassphrase, apiKey, apiSecret)
	if err != nil {
		fmt.Println("\n~~~~~~~~ERROR IN INIT~~~~~~~~")
		fmt.Println(err)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("\nTstInit fail")
		return
	}

	fmt.Println("\nTstInit success")
	time.Sleep(2 * time.Second)
}

func TstGetBalance(a *HuobiApi) {
	_, err := a.GetBalances()
	if err != nil {
		fmt.Println("\n~~~~~~~~ERROR IN GETBALANCE~~~~~~~~")
		fmt.Println(err)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("\nTstGetBalance fail")
		return
	}

	jb, _ := json.Marshal(a.balances)
	fmt.Println("\nBalanceWallet:", string(jb))
	fmt.Println("\nTstGetBalance success")
	time.Sleep(2 * time.Second)
}

func TstGetOrderBook(a *HuobiApi, symbol string) {
	ob, err := a.GetOrderBook(symbol)
	if err != nil {
		fmt.Println("\n~~~~~~~~ERROR IN GETORDERBOOK~~~~~~~~")
		fmt.Println(err)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("\nTstGetOrderBook fail")
		return
	}

	jb, _ := json.Marshal(ob)
	fmt.Println("\nOrderBook:", string(jb))
	fmt.Println("\nTstGetOrderBook success")
	time.Sleep(2 * time.Second)
}

func TstBuy(a *HuobiApi, symbol string, amount float64, price float64) string {
	bo, err := a.Buy(symbol, amount, price)
	if err != nil {
		fmt.Println("\n~~~~~~~~ERROR IN BUY~~~~~~~~")
		fmt.Println(err)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("\nTstBuy fail")
		return ""
	}

	jb, _ := json.Marshal(bo)
	fmt.Println("\nBuy:", string(jb))
	fmt.Println("\nTstBuy success")
	time.Sleep(2 * time.Second)

	return bo.ID
}

func TstSell(a *HuobiApi, symbol string, amount float64, price float64) string {
	so, err := a.Sell(symbol, amount, price)
	if err != nil {
		fmt.Println("\n~~~~~~~~ERROR IN SELL~~~~~~~~")
		fmt.Println(err)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("\nTstSell fail")
		return ""
	}

	jb, _ := json.Marshal(so)
	fmt.Println("\nSell: ", string(jb))
	fmt.Println("\nTstSell success")
	time.Sleep(2 * time.Second)

	return so.ID
}

func TstCancelOrder(a *HuobiApi, symbol string, id string) {
	err := a.CancelOrder(symbol, id)
	if err != nil {
		fmt.Println("\n~~~~~~~~ERROR IN CANCELORDER~~~~~~~~")
		fmt.Println(err)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("\nTstCancelOrder fail")
		return
	}

	fmt.Println("\nTstCancelOrder success")
	time.Sleep(2 * time.Second)
}

func TstGetOrderStatus(a *HuobiApi, id string, symbol string) {
	os, err := a.GetOrderStatus(id, symbol)
	if err != nil {
		fmt.Println("\n~~~~~~~~ERROR IN GETORDERSTATUS~~~~~~~~")
		fmt.Println(err)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("\nTstGetOrderStatus fail")
		return
	}

	jb, _ := json.Marshal(os)
	fmt.Println("\nOrderStatus:", string(jb))
	fmt.Println("\nTstGetOrderStatus success")
	time.Sleep(2 * time.Second)
}

func TstGetMyOpenOrders(a *HuobiApi, symbol string) {
	oo, err := a.GetMyOpenOrders(symbol)
	if err != nil {
		fmt.Println("\n~~~~~~~~ERROR IN GETMYOPENORDER~~~~~~~~")
		fmt.Println(err)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("\nTstGetMyOpenOrders fail")
		return
	}

	jb, _ := json.Marshal(oo)
	fmt.Println("\nOpenOrders:", string(jb))
	fmt.Println("\nTstGetOrderStatus success")
	time.Sleep(2 * time.Second)
}

func TestFunc(t *testing.T) {
	apiPassphrase := ""
	apiKey := "YOUR-API-KEY"
	apiSecret := "YOUR-SECRET-KEY"

	symbol := "USDT_TON"
	price := 1000.
	amount := 0.1
	_, _, _ = amount, price, symbol

	// ============================= TESTS =============================
	log.Println("start")
	a := new(HuobiApi)

	fmt.Println("==========INIT==========")
	TstInit(a, apiPassphrase, apiKey, apiSecret)

	//fmt.Println("==========GetBalance==========")
	//TstGetBalance(a)

	//fmt.Println("==========GetOrderBook==========")
	//TstGetOrderBook(a, symbol)

	//fmt.Println("==========Buy==========")
	//TstBuy(a, symbol, amount, price)

	//fmt.Println("==========Sell==========")
	//TstSell(a, symbol, amount, price)

	//fmt.Println("==========CancelOrder==========")
	//TstCancelOrder(a, symbol, "676893366421671")

	//fmt.Println("==========GetOrderStatus==========")
	//TstGetOrderStatus(a, "676893366421671", symbol)

	//fmt.Println("==========GetMyOpenOrders==========")
	//TstGetMyOpenOrders(a, symbol)

	log.Println("end")
	return
}

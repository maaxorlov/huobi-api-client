package huobiapi

import (
	"apiClient/apiClient"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//HuobiApi struct for api
type HuobiApi struct {
	root      string
	host      string
	akKey     string
	akValue   string
	smKey     string
	smValue   string
	svKey     string
	svValue   string
	sKey      string
	tKey      string
	apiSecret string

	spotAccID       string
	balances        *map[string]apiclient.Balance
	prevBalanceTime *time.Time

	debugData DebugStr
}

//Init implemets apiclient.APIClient interface
func (A *HuobiApi) Init(apiPassphrase string, apiKey string, apiSecret string) error {
	A.root = "https://api.huobi.pro"
	A.akKey = "AccessKeyId"
	A.akValue = apiKey
	A.smKey = "SignatureMethod"
	A.smValue = "HmacSHA256"
	A.svKey = "SignatureVersion"
	A.svValue = "2"
	A.sKey = "Signature"
	A.tKey = "Timestamp"
	A.host = "api.huobi.pro"
	A.apiSecret = apiSecret

	response, err := A.getSpotAcc()
	if err != nil {
		return fmt.Errorf("%s {%s}", err,
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			),
		)
	}

	if response.Status != "ok" {
		return fmt.Errorf("%s {%s}", response.Message,
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			),
		)
	}

	for _, acc := range response.Data {
		if acc.Type == "spot" && acc.State == "working" {
			A.spotAccID = strconv.Itoa(acc.ID)
			break
		}
	}

	if A.spotAccID == "" {
		return fmt.Errorf("can't find spot account for trading {%s}",
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			),
		)
	}

	return nil
}

//GetBalances implemets apiclient.APIClient interface
func (A *HuobiApi) GetBalances() (*map[string]apiclient.Balance, error) {
	if A.prevBalanceTime == nil || A.prevBalanceTime.Before(time.Now().Add(-60*time.Second)) {
		response, err := A.getBalancesAcc()
		if err != nil {
			return nil, fmt.Errorf("%s {%s}", err,
				fmt.Sprintf("method: %s, rd: %s, sd: %s",
					A.debugData.lastRequestMethod,
					A.debugData.lastResponseData,
					A.debugData.lastSentData,
				),
			)
		}

		balances := make(map[string]apiclient.Balance)
		for _, balance := range response.Data.List {
			oldBalance := balances[balance.Currency]
			free := 0.
			locked := 0.

			if balance.Type == "trade" {
				free, _ = strconv.ParseFloat(balance.Balance, 64)
				locked = oldBalance.Locked
			}

			if balance.Type == "frozen" {
				free = oldBalance.Free
				locked, _ = strconv.ParseFloat(balance.Balance, 64)
			}

			if free != 0 || locked != 0 {
				balances[balance.Currency] = apiclient.Balance{
					Free:   free,
					Locked: locked,
				}
			}
		}

		if A.balances == nil {
			A.balances = &balances
		}
		*A.balances = balances

		now := time.Now()
		A.prevBalanceTime = &now
	}

	return A.balances, nil
}

//GetOrderBook implemets apiclient.APIClient interface
func (A *HuobiApi) GetOrderBook(symbol string) (*apiclient.OrderBook, error) {
	symbol = convertSymbol(symbol)

	response, err := A.getOrderBookOnPair(symbol)
	if err != nil {
		return nil, fmt.Errorf("%s {%s}", err,
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			),
		)
	}

	var orderBook apiclient.OrderBook
	for _, ask := range response.Data.Asks {
		appended := false
		price := ask[0]
		amount := ask[1]

		for i, askOB := range orderBook.Asks {
			if askOB.Price == price {
				askOB.Quantity += amount
				appended = true
				break
			} else if price < askOB.Price {
				orderBook.Asks = append(orderBook.Asks, apiclient.Order{})
				copy(orderBook.Asks[i+1:], orderBook.Asks[i:])
				orderBook.Asks[i] = apiclient.Order{Price: price, Quantity: amount}
				appended = true
				break
			}
		}
		if !appended {
			orderBook.Asks = append(orderBook.Asks,
				apiclient.Order{Price: price, Quantity: amount})
		}
	}

	for _, bid := range response.Data.Bids {
		appended := false
		price := bid[0]
		amount := bid[1]

		for i, bidOB := range orderBook.Bids {
			if bidOB.Price == price {
				bidOB.Quantity += amount
				appended = true
				break
			} else if price > bidOB.Price {
				orderBook.Bids = append(orderBook.Bids, apiclient.Order{})
				copy(orderBook.Bids[i+1:], orderBook.Bids[i:])
				orderBook.Bids[i] = apiclient.Order{Price: price, Quantity: amount}
				appended = true
				break
			}
		}

		if !appended {
			orderBook.Bids = append(orderBook.Bids,
				apiclient.Order{Price: price, Quantity: amount})
		}

	}

	return &orderBook, nil
}

//Buy implemets apiclient.APIClient interface
func (A *HuobiApi) Buy(symbol string, amount float64, price float64) (*apiclient.MakedOrder, error) {
	symbol = convertSymbol(symbol)

	response, err := A.placeOrder("buy", symbol, amount, price)
	if err != nil {
		return nil, fmt.Errorf("%s {%s}", err,
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			),
		)
	}

	id := response
	status := apiclient.Filled //?
	side := apiclient.Buy
	rate := price
	leftAmount := amount
	rightAmount := leftAmount * rate

	makedOrder := &apiclient.MakedOrder{
		ID:          id,
		Status:      status,
		LeftAmount:  leftAmount,
		RightAmount: rightAmount,
		Commission:  0.0,
		Rate:        rate,
		Side:        side,
	}

	return makedOrder, nil
}

//Sell implemets apiclient.APIClient interface
func (A *HuobiApi) Sell(symbol string, amount float64, price float64) (*apiclient.MakedOrder, error) {
	symbol = convertSymbol(symbol)

	response, err := A.placeOrder("sell", symbol, amount, price)
	if err != nil {
		return nil, fmt.Errorf("%s {%s}", err,
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			),
		)
	}

	id := response
	status := apiclient.Filled //?
	side := apiclient.Sell
	rate := price
	leftAmount := amount
	rightAmount := leftAmount * rate

	makedOrder := &apiclient.MakedOrder{
		ID:          id,
		Status:      status,
		LeftAmount:  leftAmount,
		RightAmount: rightAmount,
		Commission:  0.0,
		Rate:        rate,
		Side:        side,
	}

	return makedOrder, nil
}

//CancelOrder implemets apiclient.APIClient interface
func (A *HuobiApi) CancelOrder(symbol string, id string) error {
	err := A.cancelOrder(id)
	if err != nil {
		return fmt.Errorf("%s {%s}", err,
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			),
		)
	}

	return nil
}

//GetOrderStatus implemets apiclient.APIClient interface
func (A *HuobiApi) GetOrderStatus(id string, symbol string) (*apiclient.MakedOrder, error) {
	response, err := A.getOrderStatus(id)
	if err != nil {
		return nil, fmt.Errorf("%s {%s}", err,
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			))
	}

	if response.Status != "ok" {
		return nil, fmt.Errorf("%s {%s}", response.Message,
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			),
		)
	}

	side := apiclient.Sell
	_type := strings.Split(response.Data.Type, "-")[0]
	if _type == "buy" {
		side = apiclient.Buy
	}
	status := apiclient.NotFilled
	rate, _ := strconv.ParseFloat(response.Data.Price, 64)
	leftAmount, _ := strconv.ParseFloat(response.Data.Amount, 64)
	rightAmount := leftAmount * rate

	makedOrder := &apiclient.MakedOrder{
		ID:          id,
		Status:      status,
		LeftAmount:  leftAmount,
		RightAmount: rightAmount,
		Commission:  0.0,
		Rate:        rate,
		Side:        side,
	}

	return makedOrder, nil
}

//GetMyOrders implemets apiclient.APIClient interface
func (A *HuobiApi) GetMyOpenOrders(symbol string) (*[]apiclient.MakedOrder, error) {
	symbol = convertSymbol(symbol)

	response, err := A.getMyOpenOrders(symbol)
	if err != nil {
		return nil, fmt.Errorf("%s {%s}", err,
			fmt.Sprintf("method: %s, rd: %s, sd: %s",
				A.debugData.lastRequestMethod,
				A.debugData.lastResponseData,
				A.debugData.lastSentData,
			),
		)
	}

	s := make([]apiclient.MakedOrder, len(response.Data))
	for i, order := range response.Data {
		side := apiclient.Sell
		orderType := strings.Split(order.Type, "-")[0]
		if orderType == "buy" {
			side = apiclient.Buy
		}
		status := apiclient.NotFilled
		rate, _ := strconv.ParseFloat(order.Price, 64)
		leftAmount, _ := strconv.ParseFloat(order.Amount, 64)
		rightAmount := leftAmount * rate
		id := strconv.FormatInt(order.ID, 10)

		s[i] = apiclient.MakedOrder{
			ID:          id,
			Status:      status,
			LeftAmount:  leftAmount,
			RightAmount: rightAmount,
			Commission:  0.0,
			Rate:        rate,
			Side:        side,
		}
	}

	return &s, nil
}

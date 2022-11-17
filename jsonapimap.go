package huobiapi

type DebugStr struct {
	lastRequestMethod string
	lastResponseData  []byte
	lastSentData      []byte
}

type Accounts struct {
	Status  string          `json:"status"`
	Message string          `json:"err-msg"`
	Data    []accountStruct `json:"data"`
}

type accountStruct struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	State string `json:"state"`
}

type Balance struct {
	Status  string            `json:"status"`
	Message string            `json:"err-msg"`
	Data    balanceDataStruct `json:"data"`
}

type balanceDataStruct struct {
	List []balanceStruct `json:"list"`
}

type balanceStruct struct {
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Balance  string `json:"balance"`
}

type OrderBook struct {
	Status  string       `json:"status"`
	Message string       `json:"err-msg"`
	Data    obDataStruct `json:"tick"`
}

type obDataStruct struct {
	Bids [][]float64 `json:"bids"`
	Asks [][]float64 `json:"asks"`
}

type Orders struct {
	Status  string        `json:"status"`
	Message string        `json:"err-msg"`
	Data    []orderStruct `json:"data"`
}

type Order struct {
	Status  string      `json:"status"`
	Message string      `json:"err-msg"`
	Data    orderStruct `json:"data"`
}

type orderStruct struct {
	ID     int64  `json:"id"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
	State  string `json:"state"`
	Type   string `json:"type"`
}
type PlaceOrderRequest struct {
	AccountId string `json:"account-id"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
}

type PlaceOrder struct {
	Status  string `json:"status"`
	Message string `json:"err-msg"`
	Data    string `json:"data"`
}

type CancelOrder struct {
	Status  string `json:"status"`
	Message string `json:"err-msg"`
}

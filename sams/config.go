package sams

// 地址

var AddressListAPI = "https://api-sams.walmartmobile.cn/api/v1/sams/sams-user/receiver_address/address_list"
var SetAddressAPI = "https://api-sams.walmartmobile.cn/api/v1/sams/trade/cart/saveDeliveryAddress"

type SetAddressParam struct {
	AddressId string `json:"addressId"`
	Uid       string `json:"uid"`
	AppId     string `json:"appId"`
	SaasId    string `json:"saasId"`
}

// 商店

var StoreListAPI = "https://api-sams.walmartmobile.cn/api/v1/sams/merchant/storeApi/getRecommendStoreListByLocation"

type StoreListParam struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

// 购物车

var CartAPI = "https://api-sams.walmartmobile.cn/api/v1/sams/trade/cart/getUserCart"

type CartParam struct {
	StoreList []Store `json:"storeList"`
	AddressId string  `json:"addressId"`
	Uid       string  `json:"uid"`
}

// Payment 支付

var PaymentAPI = ""

type PaymentParam struct {
	StoreId    string   `json:"storeId"`
	ClientType string   `json:"client_type"` // APP
	SpuldList  []string `json:"spuldList [String]"`
	OrderType  string   `json:"orderType"`
}

// 商品

var GoodsInfoAPI = "https://api-sams.walmartmobile.cn/api/v1/sams/trade/settlement/checkGoodsInfo"

type GoodsInfoParam struct {
	FloorId   int64   `json:"floorId"`
	StoreId   string  `json:"storeId"`
	GoodsList []Goods `json:"goodsList"`
}

var SettleInfoAPI = "https://api-sams.walmartmobile.cn/api/v1/sams/trade/settlement/getSettleInfo"

// 运力

var CapacityDataAPI = "https://api-sams.walmartmobile.cn/api/v1/sams/delivery/portal/getCapacityData"

// 支付

var CommitPayAPI = "https://api-sams.walmartmobile.cn/api/v1/sams/trade/settlement/commitPay"

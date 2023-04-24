package config

// Client Grpc 客户端配置
type Client struct {
	// auth 鉴权客户端配置
	AuthHost string `json:"authHost" yaml:"authHost"`
	AuthPort string `json:"authPort" yaml:"authPort"`
	// order 订单客户端配置
	OrderHost string `json:"orderHost" yaml:"orderHost"`
	OrderPort string `json:"orderPort" yaml:"orderPort"`
	// product 产品客户端配置
	ProductHost string `json:"productHost" yaml:"productHost"`
	ProductPort string `json:"productPort" yaml:"productPort"`
}

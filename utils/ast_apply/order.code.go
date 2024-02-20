package main

const (

	// OrderNotFoundCode 订单未找到
	OrderNotFoundCode Code = 1 + iota + 1
	// OrderHasResolvedCode 订单已经成交或已经取消
	OrderHasResolvedCode
	// LoOrderCancelFailedCode 市价单不允许手动取消
	LoOrderCancelFailedCode //市价单不允许手动取消
	// NotBidsCode 订单簿没有买单
	NotBidsCode
	// NotAsksCode 订单簿没有卖单
	NotAsksCode
)

type Code int32

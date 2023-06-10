package utils

import (
	"github.com/chenyangguang/woocommerce"
)

var wc *woocommerce.Client

func InitWoomeCommerce(WC_CONSUMER_KEY string, WC_CONSUMER_SECRET string, WC_SHOP_NAME string) {
	// init WC app
	WCApp := woocommerce.App{
		CustomerKey:    WC_CONSUMER_KEY,
		CustomerSecret: WC_CONSUMER_SECRET,
	}
	// init WC client
	wc = woocommerce.NewClient(WCApp, WC_SHOP_NAME)
	// return WC
}
func WCClient() *woocommerce.Client {
	return wc
}

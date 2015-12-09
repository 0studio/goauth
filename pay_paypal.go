package goauth

import (
	"github.com/0studio/paypal"
	"time"
)

// https://developer.paypal.com/webapps/developer/docs/api

var payPalClient *paypal.Client

// = paypal.NewClient(PAYPAL_CLIENTID_DEV, PAYPAL_SECRET_DEV, paypal.APIBaseSandBox)

// = paypal.NewClient(PAYPAL_CLIENTID, PAYPAL_SECRET, paypal.APIBaseLive)
type PayPalParams struct {
	OrderId   string `json:"transaction_id,omitemty"`
	ProductId string `json:"appStoreProductId,omitempty"`
}

func initPayPal(clientId, secret string, isSandBox bool) {
	if payPalClient == nil {
		if isSandBox {
			payPalClient = paypal.NewClient(clientId, secret, paypal.APIBaseSandBox)
		} else {
			payPalClient = paypal.NewClient(clientId, secret, paypal.APIBaseLive)
		}
	} else {
		if (isSandBox && payPalClient.APIBase == paypal.APIBaseSandBox) && clientId == payPalClient.ClientID && secret == payPalClient.Secret {
			// use old
			return
		}
		if (!isSandBox && payPalClient.APIBase == paypal.APIBaseLive) && clientId == payPalClient.ClientID && secret == payPalClient.Secret {
			// use old
			return
		}
		if isSandBox {
			payPalClient = paypal.NewClient(clientId, secret, paypal.APIBaseSandBox)
		} else {
			payPalClient = paypal.NewClient(clientId, secret, paypal.APIBaseLive)
		}
	}
}
func GetPayPalPayment(clientId, secret, orderId string, isSandBox bool) (p *paypal.Payment) {
	initPayPal(clientId, secret, isSandBox)
	return getPayPalGetPaymentTry3(orderId)
}
func getPayPalGetPaymentTry3(orderId string) (p *paypal.Payment) {
	for i := 0; i < 3; i++ {
		p = getPayPalGetPayment(orderId)
		if p != nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func getPayPalGetPayment(orderId string) (p *paypal.Payment) {
	var err error
	p, err = payPalClient.GetPayment(orderId)
	if err == nil {
		return
	}
	return
}

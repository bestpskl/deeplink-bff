package dto

import "time"

type GetDeeplinkListResponse struct {
	Deeplinks []GetDeeplinkResponse `json:"deeplinks"`
}

type GetDeeplinkResponse struct {
	PartnerTxnCreatedDt  time.Time       `json:"partner_txn_created_dt"`
	TxnSessionValidUntil time.Time       `json:"txn_session_valid_until"`
	ProductCode          string          `json:"product_code"`
	ChannelDestination   string          `json:"channel_destination"`
	PartnerTxnRef        string          `json:"partner_txn_ref"`
	PartnerDeeplink      PartnerDeeplink `json:"partner_deeplink"`
	DynamicFields        interface{}     `json:"dynamic_fields"`
	Email                string          `json:"email"`
}

type PartnerDeeplink struct {
	Success string `json:"success"`
	Fail    string `json:"fail"`
}

type GetDeeplinkRequest struct {
	Id string `param:"id"`
}

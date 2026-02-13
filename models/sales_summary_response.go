package models

type SalesSummaryResponse struct {
	TotalRevenue      int               `json:"total_revenue"`
	TotalTransaction  int               `json:"total_transaksi"`
	BestSellerProduct BestSellerProduct `json:"produk_terlaris"`
}

type BestSellerProduct struct {
	Name    string `json:"nama"`
	QntSold int    `json:"qty_terjual"`
}

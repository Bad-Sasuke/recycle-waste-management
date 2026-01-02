export interface Receipt {
  id: string
  shop_id: string
  payment_method: string
  total_amount: number
  vat_rate: number
  vat: number
  net_total: number
  status: string
  customer_request_id?: string
  created_at: string
}

export interface ReceiptItem {
  id: string
  receipt_id: string
  shop_id: string
  waste_id: string
  name: string
  category: string
  weight: number
  unit_price: number
  price: number
}

export interface MonthlyCategoryStat {
  category: string
  weight: number
}

export interface DailyStat {
  date: string
  weight: number
}

export interface UserAnalyticsResponse {
  total_weight_this_month: number
  total_weight_last_month: number
  weight_diff_percent: number
  total_earnings_this_month: number
  monthly_category_stats: MonthlyCategoryStat[]
  daily_stats: DailyStat[]
}

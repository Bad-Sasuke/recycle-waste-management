export interface ReceiptDetailData {
  receipt: {
    id: string
    created_at: string
    payment_method: string
    total_amount: number
    vat_rate: number
    vat: number
    net_total: number
  }
  items: Array<{
    id: string
    name: string
    category: string
    weight: number
    unit_price: number
    price: number
  }>
  shop: {
    name: string
  }
}

export interface ReceiptDetailProps {
  isLoading?: boolean
  receiptData?: ReceiptDetailData | null
}

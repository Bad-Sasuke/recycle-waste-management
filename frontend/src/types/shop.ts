export interface Shop {
  shop_id?: string
  user_id?: string
  shop_code?: string
  name?: string
  description?: string
  address?: string
  phone?: string
  email?: string
  image_url?: string
  opening_time?: string
  closing_time?: string
  latitude?: number
  longitude?: number
  created_at?: string
  updated_at?: string
  average_rating?: number
  total_reviews?: number
}

export interface CreateShopRequest {
  shop_code: string
  name: string
  description?: string
  address: string
  phone?: string
  email?: string
  opening_time?: string
  closing_time?: string
  latitude?: number
  longitude?: number
}

export interface UpdateShopRequest {
  shop_code?: string
  name?: string
  description?: string
  address?: string
  phone?: string
  email?: string
  opening_time?: string
  closing_time?: string
  latitude?: number
  longitude?: number
}

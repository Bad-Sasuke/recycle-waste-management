// Interface for shop information
export interface ShopInfo {
  shop_id: string
  shop_name: string
  shop_image_url: string
}

// Interface for grouped recyclable items (returned by backend grouping)
export interface GroupedRecyclableItem {
  name: string
  category: string
  price: number
  last_update: string
  hours?: string
  url?: string
  shops: ShopInfo[] // Array of shops that have this product
  waste_ids: string[] // Array of waste IDs that have this product name
}

export default interface RecycleWaste {
  waste_id?: string
  name?: string
  category?: string
  price?: number
  last_update?: string
  hours?: string
  url?: string
}

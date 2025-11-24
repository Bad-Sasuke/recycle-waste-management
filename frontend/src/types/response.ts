export interface ResponseAPI<T = unknown> {
  message: string
  data?: T
  page?: number
  limit?: number
  total_pages?: number
  total_items?: number
}

export interface Employee {
  employee_id: string
  shop_id: string
  first_name: string
  last_name: string
  username: string
  created_at: string
  updated_at: string
}

export interface CreateEmployeeRequest {
  shop_id: string
  first_name: string
  last_name: string
  username: string
  password: string
}

export interface UpdateEmployeeRequest {
  first_name?: string
  last_name?: string
  username?: string
  password?: string
}

export interface EmployeeListResponse {
  message: string
  data: Employee[]
  page: number
  limit: number
  total_pages: number
  total_items: number
  status: number
}

export interface EmployeeResponse {
  message: string
  data: Employee
  status: number
}

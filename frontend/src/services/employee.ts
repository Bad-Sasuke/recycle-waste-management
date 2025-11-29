import http from '@/services/http'
import type {
  Employee,
  CreateEmployeeRequest,
  UpdateEmployeeRequest,
  EmployeeListResponse,
  EmployeeResponse,
} from '@/types/employee'

export const getEmployeesByShopID = async (
  shopId: string,
  page: number = 1,
  pageSize: number = 10,
): Promise<EmployeeListResponse> => {
  const response = await http.get<EmployeeListResponse>(`/employees/shop/${shopId}`, {
    params: { page, page_size: pageSize },
  })
  return response.data
}

export const createEmployee = async (data: CreateEmployeeRequest): Promise<EmployeeResponse> => {
  const response = await http.post<EmployeeResponse>('/employees', data)
  return response.data
}

export const updateEmployee = async (
  employeeId: string,
  data: UpdateEmployeeRequest,
): Promise<EmployeeResponse> => {
  const response = await http.put<EmployeeResponse>(`/employees/${employeeId}`, data)
  return response.data
}

export const deleteEmployee = async (employeeId: string): Promise<EmployeeResponse> => {
  const response = await http.delete<EmployeeResponse>(`/employees/${employeeId}`)
  return response.data
}

export const getEmployeeByID = async (employeeId: string): Promise<EmployeeResponse> => {
  const response = await http.get<EmployeeResponse>(`/employees/${employeeId}`)
  return response.data
}

export const employeeLogin = async (
  shopCode: string,
  username: string,
  password: string,
): Promise<{ token: string; employee: Employee }> => {
  const response = await http.post<{ token: string; employee: Employee }>('/employees/login', {
    shop_code: shopCode,
    username,
    password,
  })
  return response.data
}

import http from './http'
import type { UserAnalyticsResponse } from '@/types/receipt'

export const receiptService = {
  async getMyAnalytics(): Promise<UserAnalyticsResponse> {
    const response = await http.get<{ success: boolean; data: UserAnalyticsResponse }>(
      '/receipts/analytics/me',
    )
    return response.data.data
  },
}

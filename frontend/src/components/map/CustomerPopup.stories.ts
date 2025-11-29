import type { Meta, StoryObj } from '@storybook/vue3-vite'
import CustomerPopup from './CustomerPopup.vue'
import type { CustomerRequest } from '@/types/customer_request'

const meta = {
  title: 'Components/Map/CustomerPopup',
  component: CustomerPopup,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component:
          'Customer request popup shown on the map for shop owners to accept or reject waste collection requests.',
      },
    },
  },
  render: (args) => ({
    components: { CustomerPopup },
    setup() {
      return { args }
    },
    template: `
      <div style="padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 600px; display: flex; align-items: center; justify-center;">
        <div style="max-width: 320px;">
          <CustomerPopup v-bind="args" />
        </div>
      </div>
    `,
  }),
} satisfies Meta<typeof CustomerPopup>

export default meta
type Story = StoryObj<typeof meta>

const mockPendingRequest: CustomerRequest = {
  customer_request_id: 'req-123456789',
  latitude: 13.7563,
  longitude: 100.5018,
  description: 'มีขยะพลาสติกและกระดาษเยอะครับ',
  status: 'pending',
  distance: 2.5,
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
}

const mockAcceptedRequest: CustomerRequest = {
  ...mockPendingRequest,
  customer_request_id: 'req-987654321',
  status: 'accepted',
}

export const PendingRequest: Story = {
  args: {
    request: mockPendingRequest,
    distance: 2.5,
  },
}

export const AcceptedRequest: Story = {
  args: {
    request: mockAcceptedRequest,
    distance: 1.8,
  },
}

export const NearbyRequest: Story = {
  args: {
    request: {
      ...mockPendingRequest,
      description: 'อยู่ใกล้ๆ มีของรีไซเคิลให้มารับ',
    },
    distance: 0.5,
  },
}

export const FarRequest: Story = {
  args: {
    request: {
      ...mockPendingRequest,
      description: 'อยู่ไกลหน่อยแต่มีของเยอะมาก',
    },
    distance: 15.3,
  },
}

export const NoDescription: Story = {
  args: {
    request: {
      ...mockPendingRequest,
      description: '',
    },
    distance: 3.2,
  },
}

export const MultipleRequests: Story = {
  args: {
    request: mockPendingRequest, // Provide default args to satisfy type requirement
    distance: 0,
  },
  render: (args) => ({
    components: { CustomerPopup },
    setup() {
      const requests = [
        { request: mockPendingRequest, distance: 2.5 },
        { request: mockAcceptedRequest, distance: 1.8 },
        {
          request: { ...mockPendingRequest, customer_request_id: 'req-333', description: '' },
          distance: 5.1,
        },
      ]
      return { requests, args }
    },
    template: `
      <div style="padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 600px;">
        <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; max-width: 1000px; margin: 0 auto;">
          <CustomerPopup
            v-for="(item, index) in requests"
            :key="index"
            :request="item.request"
            :distance="item.distance"
          />
        </div>
      </div>
    `,
  }),
}

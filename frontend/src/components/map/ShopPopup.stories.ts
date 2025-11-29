import type { Meta, StoryObj } from '@storybook/vue3-vite'
import ShopPopup from './ShopPopup.vue'
import type { Shop } from '@/types/shop'

const meta = {
  title: 'Components/Map/ShopPopup',
  component: ShopPopup,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component:
          'Shop information popup shown on the map with shop details and navigation options.',
      },
    },
  },
  render: (args) => ({
    components: { ShopPopup },
    setup() {
      return { args }
    },
    template: `
      <div style="padding: 40px; background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); min-height: 500px; display: flex; align-items: center; justify-center;">
        <div style="max-width: 320px;">
          <ShopPopup v-bind="args" />
        </div>
      </div>
    `,
  }),
} satisfies Meta<typeof ShopPopup>

export default meta
type Story = StoryObj<typeof meta>

const mockShop: Shop = {
  shop_id: '1',
  name: 'ร้านรีไซเคิลสีเขียว',
  address: '123 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพมหานคร 10110',
  phone: '02-123-4567',
  opening_time: '08:00',
  closing_time: '18:00',
  latitude: 13.7563,
  longitude: 100.5018,
  image_url: 'https://images.unsplash.com/photo-1532996122724-e3c354a0b15b?w=400&h=300',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
}

export const WithImage: Story = {
  args: {
    shop: mockShop,
  },
}

export const WithoutImage: Story = {
  args: {
    shop: {
      ...mockShop,
      image_url: '',
    },
  },
}

export const MinimalInfo: Story = {
  args: {
    shop: {
      shop_id: '2',
      name: 'ร้านรีไซเคิลน้อย',
      address: 'ถนนพระราม 4',
      latitude: 13.7563,
      longitude: 100.5018,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    } as Shop,
  },
}

export const MultipleShops: Story = {
  args: {
    shop: mockShop, // Provide default args to satisfy type requirement
  },
  render: (args) => ({
    components: { ShopPopup },
    setup() {
      const shops = [
        mockShop,
        {
          ...mockShop,
          shop_id: 'shop-2',
          name: 'ร้านกระดาษรีไซเคิล',
          address: '456 ถ.เพชรบุรี',
          image_url: '',
        },
        {
          ...mockShop,
          shop_id: 'shop-3',
          name: 'ร้านรับซื้อของเก่า ลุงแดง',
          opening_time: '09:00',
          closing_time: '18:00',
        },
      ]
      return { shops, args }
    },
    template: `
      <div style="padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 600px;">
        <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; max-width: 1000px; margin: 0 auto;">
          <ShopPopup
            v-for="(shop, index) in shops"
            :key="index"
            :shop="shop"
          />
        </div>
      </div>
    `,
  }),
}

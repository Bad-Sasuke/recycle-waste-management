import type { Meta, StoryObj } from '@storybook/vue3-vite'
import MarketplacePage from './MarketplacePage.vue'

const meta = {
  title: 'Components/MarketplacePage',
  component: MarketplacePage,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
    docs: {
      description: {
        component:
          'Marketplace page component displaying recyclable products with filtering and search capabilities.',
      },
    },
  },
} satisfies Meta<typeof MarketplacePage>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

import type { Meta, StoryObj } from '@storybook/vue3-vite'
import ShopNotification from './ShopNotification.vue'

const meta = {
  title: 'Components/ShopNotification',
  component: ShopNotification,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
    docs: {
      description: {
        component:
          "Notification component that prompts moderators to create a shop if they don't have one.",
      },
    },
  },
} satisfies Meta<typeof ShopNotification>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => ({
    components: { ShopNotification },
    template: `
      <div style="height: 400px; position: relative; background: #f5f5f5;">
        <ShopNotification />
        <div style="padding: 20px;">
          <p>This notification appears in the bottom-right corner for moderators without a shop.</p>
        </div>
      </div>
    `,
  }),
}

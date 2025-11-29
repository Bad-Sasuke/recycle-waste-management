import type { Meta, StoryObj } from '@storybook/vue3-vite'
import IconSupport from './IconSupport.vue'

const meta = {
  title: 'Components/Icons/IconSupport',
  component: IconSupport,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component: 'Heart icon used for support or favorite actions.',
      },
    },
  },
} satisfies Meta<typeof IconSupport>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const WithColor: Story = {
  render: () => ({
    components: { IconSupport },
    template: '<div style="color: red; font-size: 48px;"><IconSupport /></div>',
  }),
}

export const AllSizes: Story = {
  render: () => ({
    components: { IconSupport },
    template: `
      <div style="display: flex; gap: 20px; align-items: center;">
        <div style="font-size: 16px;"><IconSupport /></div>
        <div style="font-size: 24px;"><IconSupport /></div>
        <div style="font-size: 32px;"><IconSupport /></div>
        <div style="font-size: 48px;"><IconSupport /></div>
      </div>
    `,
  }),
}

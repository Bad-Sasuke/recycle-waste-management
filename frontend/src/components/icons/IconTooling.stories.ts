import type { Meta, StoryObj } from '@storybook/vue3-vite'
import IconTooling from './IconTooling.vue'

const meta = {
  title: 'Components/Icons/IconTooling',
  component: IconTooling,
  tags: ['autodocs'],
} satisfies Meta<typeof IconTooling>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

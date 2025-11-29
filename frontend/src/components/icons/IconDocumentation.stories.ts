import type { Meta, StoryObj } from '@storybook/vue3-vite'
import IconDocumentation from './IconDocumentation.vue'

const meta = {
  title: 'Components/Icons/IconDocumentation',
  component: IconDocumentation,
  tags: ['autodocs'],
} satisfies Meta<typeof IconDocumentation>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

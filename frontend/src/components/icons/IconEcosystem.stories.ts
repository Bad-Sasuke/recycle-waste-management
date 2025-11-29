import type { Meta, StoryObj } from '@storybook/vue3-vite'
import IconEcosystem from './IconEcosystem.vue'

const meta = {
  title: 'Components/Icons/IconEcosystem',
  component: IconEcosystem,
  tags: ['autodocs'],
} satisfies Meta<typeof IconEcosystem>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

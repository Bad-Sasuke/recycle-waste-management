import type { Meta, StoryObj } from '@storybook/vue3-vite'
import IconCommunity from './IconCommunity.vue'

// IconCommunity Stories
const metaCommunity = {
  title: 'Components/Icons/IconCommunity',
  component: IconCommunity,
  tags: ['autodocs'],
} satisfies Meta<typeof IconCommunity>

export default metaCommunity

export const Community: StoryObj<typeof metaCommunity> = {}

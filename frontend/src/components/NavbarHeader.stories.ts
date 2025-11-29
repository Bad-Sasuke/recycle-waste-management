import type { Meta, StoryObj } from '@storybook/vue3-vite'
import NavbarHeader from './NavbarHeader.vue'

const meta = {
  title: 'Components/NavbarHeader',
  component: NavbarHeader,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
    docs: {
      description: {
        component:
          'Main navigation header with logo, menu items, language switcher, and user authentication.',
      },
    },
  },
} satisfies Meta<typeof NavbarHeader>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

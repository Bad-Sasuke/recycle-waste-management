import type { Meta, StoryObj } from '@storybook/vue3-vite'
import SwitchLangDesktop from './SwitchLangDesktop.vue'

const meta = {
  title: 'Components/SwitchLang/SwitchLangDesktop',
  component: SwitchLangDesktop,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component: 'Desktop language switcher dropdown component. Hidden on mobile devices.',
      },
    },
  },
} satisfies Meta<typeof SwitchLangDesktop>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

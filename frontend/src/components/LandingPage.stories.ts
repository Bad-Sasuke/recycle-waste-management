import type { Meta, StoryObj } from '@storybook/vue3-vite'
import LandingPage from './LandingPage.vue'

const meta = {
  title: 'Components/LandingPage',
  component: LandingPage,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
    docs: {
      description: {
        component:
          'Landing page component showcasing the recycle waste management platform features and call-to-actions.',
      },
    },
  },
} satisfies Meta<typeof LandingPage>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

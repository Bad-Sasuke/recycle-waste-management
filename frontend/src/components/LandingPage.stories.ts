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
          'Premium SaaS Landing Page for the Recycle Waste Management platform. Features a dark-mode hero, pricing plans, testimonials, and animated elements.',
      },
    },
  },
} satisfies Meta<typeof LandingPage>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

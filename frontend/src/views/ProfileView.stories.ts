import type { Meta, StoryObj } from '@storybook/vue3'
import ProfileView from './ProfileView.vue'

const meta = {
  title: 'Pages/ProfileView',
  component: ProfileView,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
  },
  // Note: This page requires Authentication and API mocking to function fully.
  // In a real Storybook setup, consider using msw-storybook-addon or mocking the store.
} satisfies Meta<typeof ProfileView>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => ({
    components: { ProfileView },
    template: '<ProfileView />',
  }),
}

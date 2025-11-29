import type { Meta, StoryObj } from '@storybook/vue3-vite'
import UserLocationPopup from './UserLocationPopup.vue'

const meta = {
  title: 'Components/Map/UserLocationPopup',
  component: UserLocationPopup,
  tags: ['autodocs'],
  argTypes: {
    status: {
      control: 'select',
      options: ['idle', 'pending', 'accepted'],
    },
  },
  parameters: {
    docs: {
      description: {
        component:
          'User location popup shown on the map with different status states (idle, pending, accepted).',
      },
    },
  },
  render: (args) => ({
    components: { UserLocationPopup },
    setup() {
      return { args }
    },
    template: `
      <div style="padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 400px; display: flex; align-items: center; justify-center;">
        <div style="max-width: 300px;">
          <UserLocationPopup v-bind="args" />
        </div>
      </div>
    `,
  }),
} satisfies Meta<typeof UserLocationPopup>

export default meta
type Story = StoryObj<typeof meta>

export const Idle: Story = {
  args: {
    status: 'idle',
  },
}

export const Pending: Story = {
  args: {
    status: 'pending',
  },
}

export const Accepted: Story = {
  args: {
    status: 'accepted',
  },
}

export const AllStates: Story = {
  render: () => ({
    components: { UserLocationPopup },
    template: `
      <div style="padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 500px;">
        <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 20px;">
          <div>
            <h3 style="color: white; margin-bottom: 10px; text-align: center;">Idle State</h3>
            <UserLocationPopup status="idle" />
          </div>
          <div>
            <h3 style="color: white; margin-bottom: 10px; text-align: center;">Pending State</h3>
            <UserLocationPopup status="pending" />
          </div>
          <div>
            <h3 style="color: white; margin-bottom: 10px; text-align: center;">Accepted State</h3>
            <UserLocationPopup status="accepted" />
          </div>
        </div>
      </div>
    `,
  }),
}

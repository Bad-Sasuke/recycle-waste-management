import type { Meta, StoryObj } from '@storybook/vue3-vite'
import LiquidMarker from './LiquidMarker.vue'

const meta = {
  title: 'Components/Map/LiquidMarker',
  component: LiquidMarker,
  tags: ['autodocs'],
  argTypes: {
    colorClass: { control: 'text' },
  },
  parameters: {
    docs: {
      description: {
        component: 'Animated liquid-style map marker with floating and pulse effects.',
      },
    },
  },
} satisfies Meta<typeof LiquidMarker>

export default meta
type Story = StoryObj<typeof meta>

export const Primary: Story = {
  args: {
    colorClass: 'bg-primary',
  },
}

export const Success: Story = {
  args: {
    colorClass: 'bg-success',
  },
}

export const Warning: Story = {
  args: {
    colorClass: 'bg-warning',
  },
}

export const Error: Story = {
  args: {
    colorClass: 'bg-error',
  },
}

export const Info: Story = {
  args: {
    colorClass: 'bg-info',
  },
}

export const AllColors: Story = {
  render: () => ({
    components: { LiquidMarker },
    template: `
      <div style="display: flex; gap: 60px; padding: 40px; background: #f5f5f5; flex-wrap: wrap;">
        <div style="text-align: center;">
          <LiquidMarker colorClass="bg-primary" />
          <p style="margin-top: 10px;">Primary</p>
        </div>
        <div style="text-align: center;">
          <LiquidMarker colorClass="bg-success" />
          <p style="margin-top: 10px;">Success</p>
        </div>
        <div style="text-align: center;">
          <LiquidMarker colorClass="bg-warning" />
          <p style="margin-top: 10px;">Warning</p>
        </div>
        <div style="text-align: center;">
          <LiquidMarker colorClass="bg-error" />
          <p style="margin-top: 10px;">Error</p>
        </div>
        <div style="text-align: center;">
          <LiquidMarker colorClass="bg-info" />
          <p style="margin-top: 10px;">Info</p>
        </div>
      </div>
    `,
  }),
}

export const WithCustomIcon: Story = {
  render: () => ({
    components: { LiquidMarker },
    template: `
      <div style="padding: 40px; background: #f5f5f5;">
        <LiquidMarker colorClass="bg-primary">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" width="20" height="20" style="position: relative; z-index: 2;">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
          </svg>
        </LiquidMarker>
      </div>
    `,
  }),
}

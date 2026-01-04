import type { Meta, StoryObj } from '@storybook/vue3-vite'
import LoadingOverlay from './LoadingOverlay.vue'

const meta = {
  title: 'Components/Common/LoadingOverlay',
  component: LoadingOverlay,
  tags: ['autodocs'],
  argTypes: {
    loading: { control: 'boolean' },
    text: { control: 'text' },
    fullscreen: { control: 'boolean' },
  },
  parameters: {
    docs: {
      description: {
        component:
          'A loading overlay component with a spinner and text, suitable for blocking UI while waiting for async operations.',
      },
    },
    layout: 'fullscreen',
  },
} satisfies Meta<typeof LoadingOverlay>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    loading: true,
    text: 'Loading...',
    fullscreen: false,
  },
  render: (args) => ({
    components: { LoadingOverlay },
    setup() {
      return { args }
    },
    template: `
      <div class="relative w-full h-96 bg-gray-800 text-white p-4 overflow-hidden">
        <h1 class="text-2xl font-bold mb-4">Content Behind Overlay</h1>
        <p>This is some content that is being covered by the loading overlay.</p>
        <p>You can toggle the loading state in the controls.</p>
        <LoadingOverlay v-bind="args" />
      </div>
    `,
  }),
}

export const CustomText: Story = {
  args: {
    loading: true,
    text: 'Connecting to Server...',
    fullscreen: false,
  },
  render: Default.render,
}

export const Fullscreen: Story = {
  args: {
    loading: true,
    text: 'Loading Application...',
    fullscreen: true,
  },
  render: (args) => ({
    components: { LoadingOverlay },
    setup() {
      return { args }
    },
    template: `
        <div>
            <div class="p-8">
                <h1>This overlay covers the entire screen</h1>
                <p>Scroll down to see it stays fixed.</p>
                <div style="height: 150vh; background: linear-gradient(to bottom, #111827, #374151);"></div>
            </div>
            <LoadingOverlay v-bind="args" />
        </div>
    `,
  }),
}

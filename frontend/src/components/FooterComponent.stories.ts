import type { Meta, StoryObj } from '@storybook/vue3-vite'
import FooterComponent from './FooterComponent.vue'

const meta = {
  title: 'Components/FooterComponent',
  component: FooterComponent,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
    docs: {
      description: {
        component:
          'Footer component displaying copyright and navigation links. Designed to sit at the bottom of the page using mt-auto.',
      },
    },
  },
  render: (args) => ({
    components: { FooterComponent },
    setup() {
      return { args }
    },
    template: `
      <div class="flex flex-col min-h-screen bg-gray-100">
        <div class="p-8">
          <h1 class="text-2xl font-bold mb-4">Page Content</h1>
          <p>This is a dummy content to demonstrate the footer positioning.</p>
          <p>The footer should appear at the bottom of this viewport.</p>
        </div>
        <FooterComponent v-bind="args" />
      </div>
    `,
  }),
} satisfies Meta<typeof FooterComponent>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

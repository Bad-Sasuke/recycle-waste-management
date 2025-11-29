import type { Meta, StoryObj } from '@storybook/vue3-vite'
import GlassCard from './GlassCard.vue'

const meta = {
  title: 'Components/Common/GlassCard',
  component: GlassCard,
  tags: ['autodocs'],
  render: (args) => ({
    components: { GlassCard },
    setup() {
      return { args }
    },
    template: `
      <div style="padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 300px;">
        <GlassCard>
          <h3 style="font-size: 1.5rem; font-weight: bold; margin-bottom: 0.5rem;">Glass Card Title</h3>
          <p>This is a glassmorphism card component with blur effect and shine animation.</p>
          <p>It creates a modern, elegant UI element.</p>
        </GlassCard>
      </div>
    `,
  }),
} satisfies Meta<typeof GlassCard>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const WithList: Story = {
  render: () => ({
    components: { GlassCard },
    template: `
      <div style="padding: 40px; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); min-height: 300px;">
        <GlassCard>
          <h3 style="font-size: 1.5rem; font-weight: bold; margin-bottom: 0.5rem;">Features</h3>
          <ul style="list-style: disc; padding-left: 1.5rem;">
            <li>Glassmorphism design</li>
            <li>Backdrop blur effect</li>
            <li>Hover animation</li>
            <li>Shine effect</li>
          </ul>
        </GlassCard>
      </div>
    `,
  }),
}

export const MultipleCards: Story = {
  render: () => ({
    components: { GlassCard },
    template: `
      <div style="padding: 40px; background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); min-height: 400px; display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px;">
        <GlassCard>
          <h3 style="font-size: 1.25rem; font-weight: bold;">Card 1</h3>
          <p>First glass card</p>
        </GlassCard>
        <GlassCard>
          <h3 style="font-size: 1.25rem; font-weight: bold;">Card 2</h3>
          <p>Second glass card</p>
        </GlassCard>
        <GlassCard>
          <h3 style="font-size: 1.25rem; font-weight: bold;">Card 3</h3>
          <p>Third glass card</p>
        </GlassCard>
      </div>
    `,
  }),
}

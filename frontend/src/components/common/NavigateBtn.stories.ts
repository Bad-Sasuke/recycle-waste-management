import type { Meta, StoryObj } from '@storybook/vue3-vite'
import NavigateBtn from './NavigateBtn.vue'

const meta = {
  title: 'Components/Common/NavigateBtn',
  component: NavigateBtn,
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['primary', 'warning', 'danger', 'info'],
    },
    onClick: { action: 'clicked' },
    href: { control: 'text' },
    target: {
      control: 'select',
      options: ['_self', '_blank', '_parent', '_top'],
    },
  } as any,
  args: {
    href: undefined,
    target: '_self',
    variant: 'primary',
  },
} satisfies Meta<typeof NavigateBtn>

export default meta
type Story = StoryObj<typeof meta>

export const Primary: Story = {
  args: {
    variant: 'primary',
  },
  render: (args) => ({
    components: { NavigateBtn },
    setup() {
      return { args }
    },
    template: '<NavigateBtn v-bind="args">Primary Button</NavigateBtn>',
  }),
}

export const Warning: Story = {
  args: {
    variant: 'warning',
  },
  render: (args) => ({
    components: { NavigateBtn },
    setup() {
      return { args }
    },
    template: '<NavigateBtn v-bind="args">Warning Button</NavigateBtn>',
  }),
}

export const Danger: Story = {
  args: {
    variant: 'danger',
  },
  render: (args) => ({
    components: { NavigateBtn },
    setup() {
      return { args }
    },
    template: '<NavigateBtn v-bind="args">Danger Button</NavigateBtn>',
  }),
}

export const Info: Story = {
  args: {
    variant: 'info',
  },
  render: (args) => ({
    components: { NavigateBtn },
    setup() {
      return { args }
    },
    template: '<NavigateBtn v-bind="args">Info Button</NavigateBtn>',
  }),
}

export const AsLink: Story = {
  args: {
    variant: 'primary',
    href: 'https://example.com',
    target: '_blank',
  },
  render: (args) => ({
    components: { NavigateBtn },
    setup() {
      return { args }
    },
    template: '<NavigateBtn v-bind="args">External Link →</NavigateBtn>',
  }),
}

export const AllVariants: Story = {
  render: () => ({
    components: { NavigateBtn },
    template: `
      <div style="display: grid; gap: 16px; max-width: 300px;">
        <NavigateBtn variant="primary">Primary</NavigateBtn>
        <NavigateBtn variant="warning">Warning</NavigateBtn>
        <NavigateBtn variant="danger">Danger</NavigateBtn>
        <NavigateBtn variant="info">Info</NavigateBtn>
      </div>
    `,
  }),
}

import type { Meta, StoryObj } from '@storybook/vue3-vite'
import ErrorText from './ErrorText.vue'

const meta = {
  title: 'Components/ErrorText',
  component: ErrorText,
  tags: ['autodocs'],
  argTypes: {
    message: { control: 'text' },
  },
} satisfies Meta<typeof ErrorText>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    message: 'This field is required',
  },
}

export const NoMessage: Story = {
  args: {
    message: undefined,
  },
}

export const LongMessage: Story = {
  args: {
    message: 'This is a very long error message that might wrap to multiple lines',
  },
}

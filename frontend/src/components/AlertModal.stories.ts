import type { Meta, StoryObj } from '@storybook/vue3-vite'
import { ref, onMounted } from 'vue'
import AlertModal from './AlertModal.vue'

// Define a custom type for our story args since AlertModal uses methods instead of props
type AlertModalArgs = {
  message: string
  title: string
  type: 'info' | 'success' | 'warning' | 'error'
}

const meta: Meta<typeof AlertModal> = {
  title: 'Components/AlertModal',
  component: AlertModal,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component:
          'Modal dialog for displaying alerts (success, error, info, warning). Controlled via exposed show() method.',
      },
    },
  },

  render: (args: any) => ({
    components: { AlertModal },
    setup() {
      const modalRef = ref<any>(null)

      const showAlert = () => {
        // Use args to control the show method
        const msg = args.message || 'This is an alert message'
        const title = args.title || 'Alert Title'
        const type = args.type || 'info'
        modalRef.value?.show(msg, title, type)
      }

      onMounted(() => {
        setTimeout(showAlert, 500)
      })

      return { modalRef, showAlert, args }
    },
    template: `
      <div>
        <button @click="showAlert" class="btn btn-primary">Show Alert</button>
        <AlertModal ref="modalRef" />
      </div>
    `,
  }),
  // Define argTypes to simulate props for Storybook controls
  argTypes: {
    message: { control: 'text' },
    title: { control: 'text' },
    type: { control: 'select', options: ['info', 'success', 'warning', 'error'] },
  } as any,
}

export default meta
type Story = StoryObj<AlertModalArgs>

export const Default: Story = {
  args: {
    message: 'This is a default alert',
    title: 'Info',
    type: 'info',
  },
}

export const Success: Story = {
  args: {
    message: 'Operation completed successfully!',
    title: 'Success',
    type: 'success',
  },
}

export const Error: Story = {
  args: {
    message: 'Something went wrong.',
    title: 'Error',
    type: 'error',
  },
}

export const Warning: Story = {
  args: {
    message: 'Please be careful.',
    title: 'Warning',
    type: 'warning',
  },
}

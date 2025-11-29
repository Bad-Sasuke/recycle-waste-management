import type { Meta, StoryObj } from '@storybook/vue3-vite'
import PopupLogin from './PopupLogin.vue'

const meta = {
  title: 'Components/PopupLogin',
  component: PopupLogin,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component: 'Login modal with GitHub and Google OAuth options.',
      },
    },
  },
  render: () => ({
    components: { PopupLogin },
    setup() {
      const openLoginModal = () => {
        const modal = document.getElementById('popup-login') as HTMLDialogElement
        modal?.showModal()
      }

      return { openLoginModal }
    },
    template: `
      <div>
        <button @click="openLoginModal" class="btn btn-primary">
          Open Login Modal
        </button>
        <PopupLogin />
      </div>
    `,
  }),
} satisfies Meta<typeof PopupLogin>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

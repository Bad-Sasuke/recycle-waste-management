import type { Meta, StoryObj } from '@storybook/vue3-vite'
import { ref } from 'vue'
import RoleSelectionModal from './RoleSelectionModal.vue'

const meta = {
  title: 'Components/RoleSelectionModal',
  component: RoleSelectionModal,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component:
          'Modal for new users to select their role (Customer or Shop Owner) when first logging in.',
      },
    },
  },
  render: () => ({
    components: { RoleSelectionModal },
    setup() {
      const modalRef = ref<InstanceType<typeof RoleSelectionModal> | null>(null)

      const openModal = () => {
        modalRef.value?.showModal()
      }

      return { modalRef, openModal }
    },
    template: `
      <div>
        <button @click="openModal" class="btn btn-primary">
          Open Role Selection Modal
        </button>
        <RoleSelectionModal ref="modalRef" />
      </div>
    `,
  }),
} satisfies Meta<typeof RoleSelectionModal>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

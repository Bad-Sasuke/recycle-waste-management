import type { Meta, StoryObj } from '@storybook/vue3-vite'
import { onMounted } from 'vue'
import PopupWaste from './PopupWaste.vue'
import { useCategoryWasteStore } from '../stores/category_waste'
import { useWastesStore } from '../stores/wastes'

const meta = {
  title: 'Components/PopupWaste',
  component: PopupWaste,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component: 'Modal popup for adding new waste items.',
      },
    },
  },
  render: (args) => ({
    components: { PopupWaste },
    setup() {
      const categoryStore = useCategoryWasteStore()
      const wastesStore = useWastesStore()

      // Mock category data
      categoryStore.category = [
        { id: '1', name: 'พลาสติก' },
        { id: '2', name: 'กระดาษ' },
        { id: '3', name: 'โลหะ' },
        { id: '4', name: 'แก้ว' },
      ]

      // Mock addWaste action
      wastesStore.addWaste = async (data: any) => {
        console.log('Mock addWaste called with:', data)
        await new Promise((resolve) => setTimeout(resolve, 1000))
        alert(`Added waste: ${data.name}`)
        return Promise.resolve(1) // Return a number as expected
      }

      const openModal = () => {
        const modal = document.getElementById('modal-waste') as HTMLDialogElement
        if (modal) modal.showModal()
      }

      onMounted(() => {
        setTimeout(() => {
          openModal()
        }, 500)
      })

      return { openModal, args }
    },
    template: `
      <div style="padding: 20px;">
        <button @click="openModal" class="btn btn-primary">
          Open Add Waste Modal
        </button>
        <PopupWaste />
      </div>
    `,
  }),
} satisfies Meta<typeof PopupWaste>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

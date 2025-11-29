import type { Meta, StoryObj } from '@storybook/vue3-vite'
import { onMounted } from 'vue'
import EditWasteModal from './EditWasteModal.vue'
import { useCategoryWasteStore } from '../stores/category_waste'
import { useWastesStore } from '../stores/wastes'

const meta = {
  title: 'Components/EditWasteModal',
  component: EditWasteModal,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component:
          'Modal for editing waste item details (name, price, category, image). Available for admin and moderator roles.',
      },
    },
  },
  render: (args) => ({
    components: { EditWasteModal },
    setup() {
      const categoryStore = useCategoryWasteStore()
      const wastesStore = useWastesStore()

      // Mock category data (id must be string based on error log, though store might define it differently)
      // If store defines it as string, we use string.
      categoryStore.category = [
        { id: '1', name: 'พลาสติก' },
        { id: '2', name: 'กระดาษ' },
        { id: '3', name: 'โลหะ' },
        { id: '4', name: 'แก้ว' },
      ] as any // Cast to any to bypass potential strict type checks if store definition varies

      // Mock wasteToEdit
      wastesStore.wasteToEdit = {
        id: 123 as any,
        name: 'ขวดพลาสติกใส',
        price: 10,
        category: 'พลาสติก',
        url: 'https://images.unsplash.com/photo-1560272564-c83b66b1ad12?w=400&h=300',
        last_update: new Date().toISOString(),
      }

      // Mock updateWaste action

      wastesStore.updateWaste = async (id: any, data: any) => {
        console.log('Mock updateWaste called with:', id, data)
        await new Promise((resolve) => setTimeout(resolve, 1000))
        alert(`Updated waste ID ${id}: ${data.name}`)
        return Promise.resolve(Number(id)) // Return number
      }

      wastesStore.clearWasteToEdit = () => {
        console.log('Mock clearWasteToEdit called')
      }

      const openModal = () => {
        const modal = document.getElementById('modal-edit-waste') as HTMLDialogElement
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
          Open Edit Waste Modal
        </button>
        <EditWasteModal />
      </div>
    `,
  }),
} satisfies Meta<typeof EditWasteModal>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

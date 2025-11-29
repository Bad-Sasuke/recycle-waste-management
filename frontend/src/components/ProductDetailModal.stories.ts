import type { Meta, StoryObj } from '@storybook/vue3-vite'
import { ref, onMounted } from 'vue'
import ProductDetailModal from './ProductDetailModal.vue'

const meta = {
  title: 'Components/ProductDetailModal',
  component: ProductDetailModal,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component: 'Modal displaying detailed product information in the marketplace.',
      },
    },
  },
  render: (args) => ({
    components: { ProductDetailModal },
    setup() {
      const modalRef = ref<any>(null)

      const mockProduct = {
        name: 'ขวดพลาสติกใส (PET)',
        price: 8.5,
        category: 'พลาสติก',
        url: 'https://images.unsplash.com/photo-1560272564-c83b66b1ad12?w=400&h=300',
        last_update: new Date().toISOString(),
        shops: [
          {
            shop_name: 'ร้านรีไซเคิลสุขุมวิท',
            shop_image_url:
              'https://images.unsplash.com/photo-1532996122724-e3c354a0b15b?w=100&h=100',
          },
          { shop_name: 'วงษ์พาณิชย์ สาขา 1', shop_image_url: '' },
          { shop_name: 'รับซื้อของเก่า ลุงแดง', shop_image_url: '' },
        ],
      }

      const openModal = () => {
        modalRef.value?.openModal(mockProduct)
      }

      onMounted(() => {
        // Automatically open modal for demo purposes
        setTimeout(() => {
          openModal()
        }, 500)
      })

      return { modalRef, openModal, args }
    },
    template: `
      <div style="padding: 20px;">
        <button @click="openModal" class="btn btn-primary">
          Open Product Detail Modal
        </button>
        <ProductDetailModal ref="modalRef" />
      </div>
    `,
  }),
} satisfies Meta<typeof ProductDetailModal>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const NoImage: Story = {
  render: (args) => ({
    components: { ProductDetailModal },
    setup() {
      const modalRef = ref<any>(null)

      const mockProduct = {
        name: 'กระดาษลัง',
        price: 4.0,
        category: 'กระดาษ',
        url: '',
        last_update: new Date().toISOString(),
        shops: [{ shop_name: 'ร้าน A', shop_image_url: '' }],
      }

      const openModal = () => {
        modalRef.value?.openModal(mockProduct)
      }

      return { modalRef, openModal, args }
    },
    template: `
      <div style="padding: 20px;">
        <button @click="openModal" class="btn btn-primary">
          Open No Image Product Modal
        </button>
        <ProductDetailModal ref="modalRef" />
      </div>
    `,
  }),
}

export const NoShops: Story = {
  render: (args) => ({
    components: { ProductDetailModal },
    setup() {
      const modalRef = ref<any>(null)

      const mockProduct = {
        name: 'เศษเหล็ก',
        price: 12.0,
        category: 'โลหะ',
        url: 'https://images.unsplash.com/photo-1610701596007-11502861dcfa?w=400&h=300',
        last_update: new Date().toISOString(),
        shops: [],
      }

      const openModal = () => {
        modalRef.value?.openModal(mockProduct)
      }

      return { modalRef, openModal, args }
    },
    template: `
      <div style="padding: 20px;">
        <button @click="openModal" class="btn btn-primary">
          Open No Shops Product Modal
        </button>
        <ProductDetailModal ref="modalRef" />
      </div>
    `,
  }),
}

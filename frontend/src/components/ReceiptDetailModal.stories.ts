import type { Meta, StoryObj } from '@storybook/vue3-vite'
import { ref, onMounted } from 'vue'
import ReceiptDetailModal from './ReceiptDetailModal.vue'
import type { ReceiptDetailData } from './ReceiptDetailModal.vue'

const meta = {
  title: 'Components/ReceiptDetailModal',
  component: ReceiptDetailModal,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component: 'Modal displaying detailed receipt information including items and summary.',
      },
    },
  },
  render: (args) => ({
    components: { ReceiptDetailModal },
    setup() {
      const modalRef = ref<any>(null)

      const mockReceiptData: ReceiptDetailData = {
        receipt: {
          id: '550e8400-e29b-41d4-a716-446655440000',
          created_at: new Date().toISOString(),
          payment_method: 'เงินสด',
          total_amount: 1250.0,
          vat_rate: 0.07,
          vat: 87.5,
          net_total: 1337.5,
        },
        items: [
          {
            id: '1',
            name: 'ขวดพลาสติกใส (PET)',
            category: 'พลาสติก',
            weight: 50,
            unit_price: 8.5,
            price: 425.0,
          },
          {
            id: '2',
            name: 'กระดาษลัง',
            category: 'กระดาษ',
            weight: 100,
            unit_price: 4.0,
            price: 400.0,
          },
          {
            id: '3',
            name: 'เศษอลูมิเนียม',
            category: 'โลหะ',
            weight: 25,
            unit_price: 17.0,
            price: 425.0,
          },
        ],
        shop: {
          name: 'ร้านรีไซเคิลสุขุมวิท',
        },
      }

      const openModal = () => {
        modalRef.value?.openModal()
      }

      onMounted(() => {
        // Automatically open modal for demo purposes
        setTimeout(() => {
          openModal()
        }, 500)
      })

      return { modalRef, openModal, args, mockReceiptData }
    },
    template: `
      <div style="padding: 20px;">
        <button @click="openModal" class="btn btn-primary">
          เปิดรายละเอียดใบเสร็จ
        </button>
        <ReceiptDetailModal
          ref="modalRef"
          :receipt-data="mockReceiptData"
          :is-loading="false"
        />
      </div>
    `,
  }),
} satisfies Meta<typeof ReceiptDetailModal>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Loading: Story = {
  render: (args) => ({
    components: { ReceiptDetailModal },
    setup() {
      const modalRef = ref<any>(null)

      const openModal = () => {
        modalRef.value?.openModal()
      }

      onMounted(() => {
        setTimeout(() => {
          openModal()
        }, 500)
      })

      return { modalRef, openModal, args }
    },
    template: `
      <div style="padding: 20px;">
        <button @click="openModal" class="btn btn-primary">
          เปิดรายละเอียดใบเสร็จ (Loading)
        </button>
        <ReceiptDetailModal
          ref="modalRef"
          :receipt-data="null"
          :is-loading="true"
        />
      </div>
    `,
  }),
}

export const NoData: Story = {
  render: (args) => ({
    components: { ReceiptDetailModal },
    setup() {
      const modalRef = ref<any>(null)

      const openModal = () => {
        modalRef.value?.openModal()
      }

      onMounted(() => {
        setTimeout(() => {
          openModal()
        }, 500)
      })

      return { modalRef, openModal, args }
    },
    template: `
      <div style="padding: 20px;">
        <button @click="openModal" class="btn btn-primary">
          เปิดรายละเอียดใบเสร็จ (ไม่มีข้อมูล)
        </button>
        <ReceiptDetailModal
          ref="modalRef"
          :receipt-data="null"
          :is-loading="false"
        />
      </div>
    `,
  }),
}

export const EmptyItems: Story = {
  render: (args) => ({
    components: { ReceiptDetailModal },
    setup() {
      const modalRef = ref<any>(null)

      const mockReceiptData: ReceiptDetailData = {
        receipt: {
          id: '550e8400-e29b-41d4-a716-446655440000',
          created_at: new Date().toISOString(),
          payment_method: 'เงินสด',
          total_amount: 0,
          vat_rate: 0.07,
          vat: 0,
          net_total: 0,
        },
        items: [],
        shop: {
          name: 'ร้านรีไซเคิลธนบุรี',
        },
      }

      const openModal = () => {
        modalRef.value?.openModal()
      }

      onMounted(() => {
        setTimeout(() => {
          openModal()
        }, 500)
      })

      return { modalRef, openModal, args, mockReceiptData }
    },
    template: `
      <div style="padding: 20px;">
        <button @click="openModal" class="btn btn-primary">
          เปิดรายละเอียดใบเสร็จ (ไม่มีรายการ)
        </button>
        <ReceiptDetailModal
          ref="modalRef"
          :receipt-data="mockReceiptData"
          :is-loading="false"
        />
      </div>
    `,
  }),
}

export const SingleItem: Story = {
  render: (args) => ({
    components: { ReceiptDetailModal },
    setup() {
      const modalRef = ref<any>(null)

      const mockReceiptData: ReceiptDetailData = {
        receipt: {
          id: '650e8400-e29b-41d4-a716-446655440001',
          created_at: new Date().toISOString(),
          payment_method: 'โอนเงิน',
          total_amount: 850.0,
          vat_rate: 0.07,
          vat: 59.5,
          net_total: 909.5,
        },
        items: [
          {
            id: '1',
            name: 'ขวดแก้ว',
            category: 'แก้ว',
            weight: 100,
            unit_price: 8.5,
            price: 850.0,
          },
        ],
        shop: {
          name: 'ร้านวงษ์พาณิชย์',
        },
      }

      const openModal = () => {
        modalRef.value?.openModal()
      }

      onMounted(() => {
        setTimeout(() => {
          openModal()
        }, 500)
      })

      return { modalRef, openModal, args, mockReceiptData }
    },
    template: `
      <div style="padding: 20px;">
        <button @click="openModal" class="btn btn-primary">
          เปิดรายละเอียดใบเสร็จ (รายการเดียว)
        </button>
        <ReceiptDetailModal
          ref="modalRef"
          :receipt-data="mockReceiptData"
          :is-loading="false"
        />
      </div>
    `,
  }),
}

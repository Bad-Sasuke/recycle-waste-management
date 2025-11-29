import type { Meta, StoryObj } from '@storybook/vue3-vite'
import { fn } from 'storybook/test'
import CardComponent from './CardComponent.vue'

const meta = {
  title: 'Components/CardComponent',
  component: CardComponent,
  tags: ['autodocs'],
  argTypes: {
    id: { control: 'text' },
    name: { control: 'text' },
    price: { control: 'number' },
    category: { control: 'text' },
    last_update: { control: 'text' },
    url: { control: 'text' },
  },
  args: {
    onClick: fn(),
  },
} satisfies Meta<typeof CardComponent>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    id: '1',
    name: 'ขวดพลาสติก',
    price: 8,
    category: 'พลาสติก',
    last_update: '2024-11-29T10:00:00Z',
    url: 'https://images.unsplash.com/photo-1560272564-c83b66b1ad12?w=400&h=300',
  },
}

export const ExpensiveItem: Story = {
  args: {
    id: '2',
    name: 'เศษอลูมิเนียม',
    price: 50,
    category: 'โลหะ',
    last_update: '2024-11-28T15:30:00Z',
    url: 'https://images.unsplash.com/photo-1610701596007-11502861dcfa?w=400&h=300',
  },
}

export const CheapItem: Story = {
  args: {
    id: '3',
    name: 'กระดาษหนังสือพิมพ์',
    price: 2,
    category: 'กระดาษ',
    last_update: '2024-11-27T08:00:00Z',
    url: 'https://images.unsplash.com/photo-1586075010923-2dd4570fb338?w=400&h=300',
  },
}

export const LongName: Story = {
  args: {
    id: '4',
    name: 'ขวดแก้วสีใสขนาดใหญ่พิเศษ',
    price: 15,
    category: 'แก้ว',
    last_update: '2024-11-29T12:00:00Z',
    url: 'https://images.unsplash.com/photo-1532996122724-e3c354a0b15b?w=400&h=300',
  },
}

export const MultipleCards: Story = {
  render: () => ({
    components: { CardComponent },
    template: `
      <div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(256px, 1fr)); gap: 20px; padding: 20px;">
        <CardComponent
          id="1"
          name="ขวดพลาสติก"
          :price="8"
          category="พลาสติก"
          last_update="2024-11-29T10:00:00Z"
          url="https://images.unsplash.com/photo-1560272564-c83b66b1ad12?w=400&h=300"
        />
        <CardComponent
          id="2"
          name="เศษอลูมิเนียม"
          :price="50"
          category="โลหะ"
          last_update="2024-11-28T15:30:00Z"
          url="https://images.unsplash.com/photo-1610701596007-11502861dcfa?w=400&h=300"
        />
        <CardComponent
          id="3"
          name="กระดาษหนังสือพิมพ์"
          :price="2"
          category="กระดาษ"
          last_update="2024-11-27T08:00:00Z"
          url="https://images.unsplash.com/photo-1586075010923-2dd4570fb338?w=400&h=300"
        />
      </div>
    `,
  }),
}

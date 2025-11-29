import type { Meta, StoryObj } from '@storybook/vue3-vite'
import SwitchLangMobile from './SwitchLangMobile.vue'
import { useI18nStore } from '@/stores/i18n'

const meta = {
  title: 'Components/SwitchLang/SwitchLangMobile',
  component: SwitchLangMobile,
  tags: ['autodocs'],
  parameters: {
    docs: {
      description: {
        component:
          'Mobile language switcher dropdown component. Visible only on mobile devices (hidden on md+ screens).',
      },
    },
    viewport: {
      defaultViewport: 'mobile1', // Set default viewport to mobile
    },
  },
  render: (args) => ({
    components: { SwitchLangMobile },
    setup() {
      const i18nStore = useI18nStore()

      // Mock setLocale action
      i18nStore.setLocale = (lang: string) => {
        console.log('Switch language to:', lang)
        // In a real app, this would update the i18n locale
      }

      return { args }
    },
    template: `
      <div style="padding: 20px; min-height: 200px; width: 300px; border: 1px dashed #ccc;">
        <p class="mb-4 text-sm text-gray-500">Mobile Viewport Simulation (Width: 300px)</p>
        <SwitchLangMobile v-bind="args" />
      </div>
    `,
  }),
} satisfies Meta<typeof SwitchLangMobile>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const ForceVisible: Story = {
  parameters: {
    viewport: { defaultViewport: 'responsive' }, // Reset viewport
  },
  render: (args) => ({
    components: { SwitchLangMobile },
    setup() {
      const i18nStore = useI18nStore()
      i18nStore.setLocale = (lang: string) => console.log('Switch language to:', lang)
      return { args }
    },
    template: `
      <div style="padding: 20px; min-height: 200px;">
        <p class="mb-4 text-sm text-gray-500">Forced visibility by overriding md:hidden class (for demo purposes)</p>
        <!-- Remove md:hidden class for demo -->
        <div class="dropdown mx-4">
          <div tabindex="0" role="button" class="btn btn-outline p-1 text-neutral w-full">
            <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-language-hiragana" width="24" height="24" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M4 5h7" /><path d="M7 4c0 4.846 -1.415 9 -2 11.5" /><path d="M10 8.5c0 2.286 -2 4.5 -3.5 4.5" /><path d="M12 5c0 8.5 -4 12 -6 12" /><path d="M16 4a5 5 0 0 0 -5 5h8" /><path d="M19.5 16l-3.5 -6l-3.5 6" /><path d="M14 13l7 0" /></svg>
            <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-chevron-down" width="24" height="24" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M6 9l6 6l6 -6" /></svg>
          </div>
          <ul tabindex="0" class="dropdown-content menu bg-base-100 rounded-box z-[1] w-32 p-2 shadow w-full">
            <li><a>English</a></li>
            <li><a>ไทย</a></li>
          </ul>
        </div>
      </div>
    `,
  }),
}

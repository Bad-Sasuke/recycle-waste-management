<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useUsersStore } from '@/stores/users'
import PopupLogin from './PopupLogin.vue'

const { t } = useI18n()
const router = useRouter()
const usersStore = useUsersStore()

// Testimonial data
const testimonials = [
  { id: 1, name: 'John Doe', role: 'Regular User', content: t('LandingHome.testimonials.item1.content'), avatar: 'https://api.dicebear.com/7.x/avataaars-neutral/svg?seed=1' },
  { id: 2, name: 'Jane Smith', role: 'Community Leader', content: t('LandingHome.testimonials.item2.content'), avatar: 'https://api.dicebear.com/7.x/avataaars-neutral/svg?seed=2' },
  { id: 3, name: 'Mike Johnson', role: 'Environmental Advocate', content: t('LandingHome.testimonials.item3.content'), avatar: 'https://api.dicebear.com/7.x/avataaars-neutral/svg?seed=3' },
]

const handleGetStarted = () => {
  if (usersStore.isLogin) {
    // If user is logged in, redirect to marketplace
    router.push('/marketplace')
  } else {
    // If user is not logged in, show login modal
    const modal = document.getElementById('popup-login') as HTMLDialogElement
    if (modal) {
      modal.showModal()
    }
  }
}

const scrollToHowItWorks = () => {
  // Scroll to the "How It Works" section
  const howItWorksSection = document.getElementById('how-it-works');
  if (howItWorksSection) {
    howItWorksSection.scrollIntoView({ behavior: 'smooth' });
  }
}
</script>

<template>
  <PopupLogin />

  <!-- Hero Section with DaisyUI -->
  <section class="bg-gradient-to-r from-green-400 to-emerald-500 text-white py-20">
    <div class="container mx-auto px-4 text-center">
      <h1 class="text-4xl md:text-6xl font-bold mb-6">{{ $t('LandingHome.hero.title') }}</h1>
      <p class="text-xl md:text-2xl mb-8 max-w-3xl mx-auto">{{ $t('LandingHome.hero.description') }}</p>
      <div class="flex flex-col sm:flex-row justify-center gap-4">
        <button class="btn btn-primary btn-wide bg-white text-green-600 hover:bg-gray-100 border-0"
          @click="handleGetStarted">
          {{ $t('LandingHome.hero.button') }}
        </button>
        <button class="btn btn-outline btn-wide text-white border-white hover:bg-white hover:text-green-600"
          @click="scrollToHowItWorks">
          {{ $t('LandingHome.hero.secondary_button') }}
        </button>
      </div>
    </div>
  </section>

  <!-- How It Works Section with DaisyUI Cards -->
  <section id="how-it-works" class="py-20 bg-base-100">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl md:text-4xl font-bold text-center mb-16">{{ $t('LandingHome.how_it_works.title') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        <div class="card bg-base-100 shadow-xl border border-base-200">
          <figure class="px-8 pt-8">
            <div class="bg-green-100 p-4 rounded-full">
              <div class="text-3xl text-green-600">1</div>
            </div>
          </figure>
          <div class="card-body items-center text-center">
            <h3 class="card-title text-green-600">{{ $t('LandingHome.how_it_works.steps.step1.title') }}</h3>
            <p class="text-gray-600 text-sm">{{ $t('LandingHome.how_it_works.steps.step1.description') }}</p>
            <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-500" fill="none" viewBox="0 0 24 24"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </div>
          </div>
        </div>

        <div class="card bg-base-100 shadow-xl border border-base-200">
          <figure class="px-8 pt-8">
            <div class="bg-green-100 p-4 rounded-full">
              <div class="text-3xl text-green-600">2</div>
            </div>
          </figure>
          <div class="card-body items-center text-center">
            <h3 class="card-title text-green-600">{{ $t('LandingHome.how_it_works.steps.step2.title') }}</h3>
            <p class="text-gray-600 text-sm">{{ $t('LandingHome.how_it_works.steps.step2.description') }}</p>
            <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-500" fill="none" viewBox="0 0 24 24"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>
          </div>
        </div>

        <div class="card bg-base-100 shadow-xl border border-base-200">
          <figure class="px-8 pt-8">
            <div class="bg-green-100 p-4 rounded-full">
              <div class="text-3xl text-green-600">3</div>
            </div>
          </figure>
          <div class="card-body items-center text-center">
            <h3 class="card-title text-green-600">{{ $t('LandingHome.how_it_works.steps.step3.title') }}</h3>
            <p class="text-gray-600 text-sm">{{ $t('LandingHome.how_it_works.steps.step3.description') }}</p>
            <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-500" fill="none" viewBox="0 0 24 24"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
              </svg>
            </div>
          </div>
        </div>

        <div class="card bg-base-100 shadow-xl border border-base-200">
          <figure class="px-8 pt-8">
            <div class="bg-green-100 p-4 rounded-full">
              <div class="text-3xl text-green-600">4</div>
            </div>
          </figure>
          <div class="card-body items-center text-center">
            <h3 class="card-title text-green-600">{{ $t('LandingHome.how_it_works.steps.step4.title') }}</h3>
            <p class="text-gray-600 text-sm">{{ $t('LandingHome.how_it_works.steps.step4.description') }}</p>
            <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-500" fill="none" viewBox="0 0 24 24"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Features Section with DaisyUI -->
  <section class="py-20 bg-base-200">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl md:text-4xl font-bold text-center mb-16">{{ $t('LandingHome.features.title') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h3 class="card-title text-green-600">
              <span class="bg-green-100 p-2 rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </span>
              {{ $t('LandingHome.features.items.feature1.title') }}
            </h3>
            <p>{{ $t('LandingHome.features.items.feature1.description') }}</p>
          </div>
        </div>

        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h3 class="card-title text-green-600">
              <span class="bg-green-100 p-2 rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </span>
              {{ $t('LandingHome.features.items.feature2.title') }}
            </h3>
            <p>{{ $t('LandingHome.features.items.feature2.description') }}</p>
          </div>
        </div>

        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h3 class="card-title text-green-600">
              <span class="bg-green-100 p-2 rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4 4 0 003 15z" />
                </svg>
              </span>
              {{ $t('LandingHome.features.items.feature3.title') }}
            </h3>
            <p>{{ $t('LandingHome.features.items.feature3.description') }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Testimonials Section -->
  <section class="py-20 bg-base-100">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl md:text-4xl font-bold text-center mb-16">{{ $t('LandingHome.testimonials.title') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div v-for="testimonial in testimonials" :key="testimonial.id" class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <div class="flex items-center mb-4">
              <div class="avatar mr-4">
                <div class="w-16 h-16 rounded-full">
                  <img :src="testimonial.avatar" alt="User avatar" />
                </div>
              </div>
              <div>
                <h4 class="font-bold">{{ testimonial.name }}</h4>
                <p class="text-sm text-gray-500">{{ testimonial.role }}</p>
              </div>
            </div>
            <p class="text-gray-600 italic">"{{ testimonial.content }}"</p>
            <div class="rating mt-4">
              <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400" checked />
              <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400" checked />
              <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400" checked />
              <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400" checked />
              <input type="radio" name="rating-1" class="mask mask-star-2 bg-orange-400" checked />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Call to Action Section -->
  <section class="py-20 bg-gradient-to-r from-green-400 to-emerald-500 text-white">
    <div class="container mx-auto px-4 text-center">
      <h2 class="text-3xl md:text-4xl font-bold mb-6">{{ $t('LandingHome.cta.title') }}</h2>
      <p class="text-xl mb-8 max-w-2xl mx-auto">{{ $t('LandingHome.cta.description') }}</p>
      <button class="btn btn-primary btn-wide bg-white text-green-600 hover:bg-gray-100 border-0"
        @click="handleGetStarted">
        {{ $t('LandingHome.cta.button') }}
      </button>
    </div>
  </section>
</template>

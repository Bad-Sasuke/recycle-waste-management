<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// Testimonial data
const testimonials = [
  { id: 1, name: 'John Doe', role: 'Regular User', content: t('LandingHome.testimonials.item1.content'), avatar: 'https://api.dicebear.com/7.x/avataaars-neutral/svg?seed=1' },
  { id: 2, name: 'Jane Smith', role: 'Community Leader', content: t('LandingHome.testimonials.item2.content'), avatar: 'https://api.dicebear.com/7.x/avataaars-neutral/svg?seed=2' },
  { id: 3, name: 'Mike Johnson', role: 'Environmental Advocate', content: t('LandingHome.testimonials.item3.content'), avatar: 'https://api.dicebear.com/7.x/avataaars-neutral/svg?seed=3' },
]
</script>

<template>
  <!-- Hero Section with DaisyUI -->
  <section class="bg-gradient-to-r from-green-400 to-emerald-500 text-white py-20">
    <div class="container mx-auto px-4 text-center">
      <h1 class="text-4xl md:text-6xl font-bold mb-6">{{ $t('LandingHome.hero.title') }}</h1>
      <p class="text-xl md:text-2xl mb-8 max-w-3xl mx-auto">{{ $t('LandingHome.hero.description') }}</p>
      <div class="flex flex-col sm:flex-row justify-center gap-4">
        <button class="btn btn-primary btn-wide bg-white text-green-600 hover:bg-gray-100 border-0">
          {{ $t('LandingHome.hero.button') }}
        </button>
        <button class="btn btn-outline btn-wide text-white border-white hover:bg-white hover:text-green-600">
          {{ $t('LandingHome.hero.secondary_button') }}
        </button>
      </div>
    </div>
  </section>

  <!-- How It Works Section with DaisyUI Cards -->
  <section class="py-20 bg-base-100">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl md:text-4xl font-bold text-center mb-16">{{ $t('LandingHome.how_it_works.title') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div class="card bg-base-100 shadow-xl border border-base-200">
          <figure class="px-8 pt-8">
            <div class="bg-green-100 p-4 rounded-full">
              <div class="text-3xl text-green-600">1</div>
            </div>
          </figure>
          <div class="card-body items-center text-center">
            <h3 class="card-title text-green-600">{{ $t('LandingHome.how_it_works.steps.step1.title') }}</h3>
            <p class="text-gray-600">{{ $t('LandingHome.how_it_works.steps.step1.description') }}</p>
            <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16" />
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
            <p class="text-gray-600">{{ $t('LandingHome.how_it_works.steps.step2.description') }}</p>
            <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16" />
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
            <p class="text-gray-600">{{ $t('LandingHome.how_it_works.steps.step3.description') }}</p>
            <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16" />
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
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
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
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
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
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4 4 0 003 15z" />
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
      <button class="btn btn-primary btn-wide bg-white text-green-600 hover:bg-gray-100 border-0">
        {{ $t('LandingHome.cta.button') }}
      </button>
    </div>
  </section>
</template>

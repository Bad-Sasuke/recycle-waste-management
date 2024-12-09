<script setup lang="ts"></script>

<template>
  <!-- Hero Section -->
  <section class="bg-green-200 py-20 text-center">
    <h1 class="text-4xl font-bold text-gray-800">{{ $t('LandingHome.hero.title') }}</h1>
    <p class="text-xl text-gray-700 mt-4">{{ $t('LandingHome.hero.description') }}</p>
    <button class="mt-6 px-6 py-2 bg-green-500 text-white rounded-full">
      {{ $t('LandingHome.hero.button') }}
    </button>
  </section>

  <!-- How It Works -->
  <section class="py-20">
    <div class="text-center">
      <h2 class="text-3xl font-bold text-gray-800">{{ $t('LandingHome.how_it_works.title') }}</h2>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-8 mt-8">
        <div class="p-4 border rounded-lg shadow-lg">
          <div class="text-2xl font-bold text-green-500">
            {{ $t('LandingHome.how_it_works.steps.step1.title') }}
          </div>
          <p class="text-gray-600 mt-2">
            {{ $t('LandingHome.how_it_works.steps.step1.description') }}
          </p>
        </div>
        <div class="p-4 border rounded-lg shadow-lg">
          <div class="text-2xl font-bold text-green-500">
            {{ $t('LandingHome.how_it_works.steps.step2.title') }}
          </div>
          <p class="text-gray-600 mt-2">
            {{ $t('LandingHome.how_it_works.steps.step2.description') }}
          </p>
        </div>
        <div class="p-4 border rounded-lg shadow-lg">
          <div class="text-2xl font-bold text-green-500">
            {{ $t('LandingHome.how_it_works.steps.step3.title') }}
          </div>
          <p class="text-gray-600 mt-2">
            {{ $t('LandingHome.how_it_works.steps.step3.description') }}
          </p>
        </div>
      </div>
    </div>
  </section>

  <!-- Features -->
  <section class="bg-gray-100 py-20">
    <div class="text-center">
      <h2 class="text-3xl font-bold text-gray-800">{{ $t('LandingHome.features.title') }}</h2>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-8 mt-8">
        <div class="p-6 bg-white shadow-md rounded-lg">
          <h3 class="text-2xl font-semibold text-green-500">
            {{ $t('LandingHome.features.items.feature1.title') }}
          </h3>
          <p class="text-gray-600 mt-2">
            {{ $t('LandingHome.features.items.feature1.description') }}
          </p>
        </div>
        <div class="p-6 bg-white shadow-md rounded-lg">
          <h3 class="text-2xl font-semibold text-green-500">
            {{ $t('LandingHome.features.items.feature2.title') }}
          </h3>
          <p class="text-gray-600 mt-2">
            {{ $t('LandingHome.features.items.feature2.description') }}
          </p>
        </div>
        <div class="p-6 bg-white shadow-md rounded-lg">
          <h3 class="text-2xl font-semibold text-green-500">
            {{ $t('LandingHome.features.items.feature3.title') }}
          </h3>
          <p class="text-gray-600 mt-2">
            {{ $t('LandingHome.features.items.feature3.description') }}
          </p>
        </div>
      </div>
    </div>
  </section>
</template>

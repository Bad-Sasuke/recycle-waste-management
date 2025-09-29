<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const activeTab = ref('about')

// Function to change active tab
const changeTab = (tab: string) => {
  activeTab.value = tab
  // Update the URL query parameter
  router.replace({ query: { ...route.query, tab } })
}

// Set active tab based on URL query parameter
const setActiveTabFromQuery = () => {
  const tab = route.query.tab as string
  if (['about', 'privacy', 'terms'].includes(tab)) {
    activeTab.value = tab
  } else {
    activeTab.value = 'about' // Default to about tab
  }
}

// Watch for query changes in the URL
watch(() => route.query, (newQuery) => {
  const tab = newQuery.tab as string
  if (['about', 'privacy', 'terms'].includes(tab)) {
    activeTab.value = tab
  } else {
    activeTab.value = 'about' // Default to about tab
  }
}, { immediate: true })

onMounted(() => {
  setActiveTabFromQuery()
})
</script>

<template>
  <div class="min-h-screen bg-base-100 py-12">
    <div class="container mx-auto px-4 max-w-6xl">
      <!-- Page Header -->
      <div class="text-center mb-12">
        <h1 class="text-4xl font-bold text-green-700 mb-4">{{ $t('Footer.about_us') }}</h1>
        <p class="text-lg text-gray-600 max-w-3xl mx-auto">
          {{ $t('About.description') }}
        </p>
      </div>

      <!-- Tab Navigation and Content with Context7 Component -->
      <div class="bg-base-100 rounded-lg p-8 shadow-lg">
        <div class="tabs tabs-border">
          <input type="radio" name="my_tabs_2"
            class="tab px-6 md:px-8 min-w-[160px] md:min-w-[200px] text-sm md:text-base whitespace-nowrap"
            :aria-label="$t('About.about_us_tab')" :checked="activeTab === 'about'" @click="changeTab('about')" />
          <div class="tab-content bg-base-100 p-10">
            <!-- About Us Section -->
            <div v-show="activeTab === 'about'" class="prose max-w-none">
              <h2 class="text-3xl font-bold text-green-700 mb-6">{{ $t('About.about_us_title') }}</h2>

              <div class="mb-8">
                <h3 class="text-2xl font-semibold text-gray-800 mb-4">{{ $t('About.our_mission_title') }}</h3>
                <p class="text-gray-700 mb-4">
                  {{ $t('About.our_mission_desc') }}
                </p>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
                <div class="card bg-base-200 p-6 rounded-lg">
                  <h3 class="text-xl font-semibold text-green-600 mb-3">{{ $t('About.vision_title') }}</h3>
                  <p class="text-gray-700">
                    {{ $t('About.vision_desc') }}
                  </p>
                </div>
                <div class="card bg-base-200 p-6 rounded-lg">
                  <h3 class="text-xl font-semibold text-green-600 mb-3">{{ $t('About.mission_title') }}</h3>
                  <p class="text-gray-700">
                    {{ $t('About.mission_desc') }}
                  </p>
                </div>
              </div>

              <div class="mb-8">
                <h3 class="text-2xl font-semibold text-gray-800 mb-4">{{ $t('About.our_story_title') }}</h3>
                <p class="text-gray-700 mb-4">
                  {{ $t('About.our_story_desc_p1') }}
                </p>
                <p class="text-gray-700">
                  {{ $t('About.our_story_desc_p2') }}
                </p>
              </div>

              <div>
                <h3 class="text-2xl font-semibold text-gray-800 mb-4">{{ $t('About.impact_title') }}</h3>
                <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  <div class="stat">
                    <div class="stat-figure text-primary">
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                        class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                        </path>
                      </svg>
                    </div>
                    <div class="stat-title">{{ $t('About.impact_users_title') }}</div>
                    <div class="stat-value text-primary">10K+</div>
                    <div class="stat-desc">{{ $t('About.impact_users_desc') }}</div>
                  </div>

                  <div class="stat">
                    <div class="stat-figure text-secondary">
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                        class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                      </svg>
                    </div>
                    <div class="stat-title">{{ $t('About.impact_materials_title') }}</div>
                    <div class="stat-value text-secondary">500T+</div>
                    <div class="stat-desc">{{ $t('About.impact_materials_desc') }}</div>
                  </div>

                  <div class="stat">
                    <div class="stat-figure text-accent">
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                        class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4">
                        </path>
                      </svg>
                    </div>
                    <div class="stat-title">{{ $t('About.impact_cities_title') }}</div>
                    <div class="stat-value text-accent">25+</div>
                    <div class="stat-desc">{{ $t('About.impact_cities_desc') }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <input type="radio" name="my_tabs_2"
            class="tab px-6 md:px-8 min-w-[160px] md:min-w-[200px] text-sm md:text-base whitespace-nowrap"
            :aria-label="$t('About.privacy_policy_tab')" :checked="activeTab === 'privacy'"
            @click="changeTab('privacy')" />
          <div class="tab-content bg-base-100 p-10">
            <!-- Privacy Policy Section -->
            <div v-show="activeTab === 'privacy'" class="prose max-w-none">
              <h2 class="text-3xl font-bold text-green-700 mb-6">{{ $t('About.privacy_policy_title') }}</h2>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.privacy_intro_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.privacy_intro_desc') }}
                </p>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.privacy_information_title') }}</h3>
                <p class="text-gray-700 mb-3">{{ $t('About.privacy_information_desc') }}</p>
                <ul class="list-disc pl-6 text-gray-700 space-y-2">
                  <li>{{ $t('About.privacy_information_list1') }}</li>
                  <li>{{ $t('About.privacy_information_list2') }}</li>
                  <li>{{ $t('About.privacy_information_list3') }}</li>
                  <li>{{ $t('About.privacy_information_list4') }}</li>
                </ul>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.privacy_use_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.privacy_use_desc') }}
                </p>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.privacy_disclosure_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.privacy_disclosure_desc') }}
                </p>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.privacy_security_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.privacy_security_desc') }}
                </p>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.privacy_changes_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.privacy_changes_desc') }}
                </p>
              </div>

              <div>
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.privacy_contact_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.privacy_contact_desc') }}
                </p>
              </div>
            </div>
          </div>

          <input type="radio" name="my_tabs_2"
            class="tab px-6 md:px-8 min-w-[160px] md:min-w-[200px] text-sm md:text-base whitespace-nowrap"
            :aria-label="$t('About.terms_tab')" :checked="activeTab === 'terms'" @click="changeTab('terms')" />
          <div class="tab-content bg-base-100 p-10">
            <!-- Terms of Service Section -->
            <div v-show="activeTab === 'terms'" class="prose max-w-none">
              <h2 class="text-3xl font-bold text-green-700 mb-6">{{ $t('About.terms_title') }}</h2>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.terms_acceptance_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.terms_acceptance_desc') }}
                </p>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.terms_description_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.terms_description_desc') }}
                </p>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.terms_user_registration_title') }}
                </h3>
                <p class="text-gray-700">
                  {{ $t('About.terms_user_registration_desc') }}
                </p>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.terms_user_responsibilities_title') }}
                </h3>
                <p class="text-gray-700 mb-3">{{ $t('About.terms_user_responsibilities_desc') }}</p>
                <ul class="list-disc pl-6 text-gray-700 space-y-2">
                  <li>{{ $t('About.terms_user_responsibilities_list1') }}</li>
                  <li>{{ $t('About.terms_user_responsibilities_list2') }}</li>
                  <li>{{ $t('About.terms_user_responsibilities_list3') }}</li>
                  <li>{{ $t('About.terms_user_responsibilities_list4') }}</li>
                </ul>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.terms_platform_use_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.terms_platform_use_desc') }}
                </p>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.terms_content_title') }}</h3>
                <p class="text-gray-700 mb-3">{{ $t('About.terms_content_desc') }}</p>
                <ul class="list-disc pl-6 text-gray-700 space-y-2">
                  <li>{{ $t('About.terms_content_list1') }}</li>
                  <li>{{ $t('About.terms_content_list2') }}</li>
                  <li>{{ $t('About.terms_content_list3') }}</li>
                </ul>
              </div>

              <div class="mb-6">
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.terms_termination_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.terms_termination_desc') }}
                </p>
              </div>

              <div>
                <h3 class="text-xl font-semibold text-gray-800 mb-3">{{ $t('About.terms_governing_title') }}</h3>
                <p class="text-gray-700">
                  {{ $t('About.terms_governing_desc') }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">

import { useRouter } from 'vue-router'
import { useUsersStore } from '@/stores/users'
import PopupLogin from './PopupLogin.vue'
import { IconRecycle, IconChartBar, IconCheck, IconStar, IconBuildingWarehouse, IconMapPin, IconBuildingStore } from '@tabler/icons-vue'


const router = useRouter()
const usersStore = useUsersStore()

const handleGetStarted = () => {
  if (usersStore.isLogin) {
    router.push('/marketplace')
  } else {
    const modal = document.getElementById('popup-login') as HTMLDialogElement
    if (modal) {
      modal.showModal()
    }
  }
}

const scrollToHowItWorks = () => {
  const howItWorksSection = document.getElementById('how-it-works')
  if (howItWorksSection) {
    howItWorksSection.scrollIntoView({ behavior: 'smooth' })
  }
}

// Fixed English testimonials for SaaS demo look, can be i18n'd later
const testimonials = [
  { id: 1, name: 'Alex Morgan', role: 'EcoTech Solutions', content: 'This platform revolutionized how we handle industrial waste. The analytics are game-changing.', avatar: 'https://ui-avatars.com/api/?name=Alex+Morgan&background=0D8ABC&color=fff' },
  { id: 2, name: 'Sarah Chen', role: 'City Planner', content: 'Incredible visibility into our recycling metrics. A must-have for any smart city initiative.', avatar: 'https://ui-avatars.com/api/?name=Sarah+Chen&background=10B981&color=fff' },
  { id: 3, name: 'James Wilson', role: 'Sustainability Director', content: 'Efficiency increased by 40% within the first month. The marketplace feature is brilliant.', avatar: 'https://ui-avatars.com/api/?name=James+Wilson&background=F59E0B&color=fff' },
]


</script>

<template>
  <PopupLogin />

  <!-- Premium SaaS Hero Section -->
  <section class="relative bg-slate-900 overflow-hidden font-sans">
    <div class="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20"></div>
    <!-- Background Gradient Glob -->
    <div
      class="absolute top-0 right-0 -mr-20 -mt-20 w-96 h-96 bg-primary rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob">
    </div>
    <div
      class="absolute bottom-0 left-0 -ml-20 -mb-20 w-96 h-96 bg-blue-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000">
    </div>

    <div class="container mx-auto px-6 py-24 md:py-32 relative z-10">
      <div class="flex flex-col lg:flex-row items-center gap-12">
        <div class="lg:w-1/2 text-center lg:text-left">
          <div class="hero-content-anim" style="animation-delay: 0.1s">
            <div
              class="inline-block px-4 py-1.5 mb-6 text-sm font-semibold tracking-wider text-primary uppercase bg-primary/30 rounded-full border border-primary/30">
              Start Your Eco Journey Today
            </div>
          </div>
          <h1 class="text-4xl md:text-6xl font-extrabold text-white leading-tight mb-6 hero-content-anim"
            style="animation-delay: 0.2s">
            Recycle Waste <br />
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-primary to-cyan-400">Management</span>
          </h1>
          <p class="text-lg md:text-xl text-slate-300 mb-8 max-w-2xl mx-auto lg:mx-0 hero-content-anim"
            style="animation-delay: 0.3s">
            Connect with local recycling shops, trade waste materials, and track your environmental impact.
            The complete platform for modern waste management.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start hero-content-anim"
            style="animation-delay: 0.4s">
            <button @click="handleGetStarted"
              class="btn btn-primary bg-primary hover:bg-accent text-white border-0 px-8 text-lg shadow-lg shadow-primary/30 transition-all hover:-translate-y-1">
              Start for Free
            </button>
            <button @click="scrollToHowItWorks"
              class="btn btn-ghost text-white hover:bg-white/10 px-8 text-lg border border-white/20">
              How it works
            </button>
          </div>

          <div
            class="mt-10 flex items-center justify-center lg:justify-start gap-4 text-sm text-slate-400 hero-content-anim"
            style="animation-delay: 0.5s">
            <div class="flex -space-x-2">
              <div class="w-8 h-8 rounded-full bg-slate-700 border-2 border-slate-900"></div>
              <div class="w-8 h-8 rounded-full bg-slate-600 border-2 border-slate-900"></div>
              <div class="w-8 h-8 rounded-full bg-slate-500 border-2 border-slate-900"></div>
            </div>
            <p>Trusted by 500+ Eco-Leaders</p>
          </div>
        </div>

        <!-- Glassmorphism Dashboard Mockup -->
        <div class="lg:w-1/2 relative floating-dashboard">
          <div class="relative rounded-2xl bg-slate-800/50 border border-slate-700/50 backdrop-blur-xl p-6 shadow-2xl">
            <!-- Header -->
            <div class="flex items-center justify-between mb-6">
              <div class="flex items-center gap-2">
                <div class="w-3 h-3 rounded-full bg-red-400"></div>
                <div class="w-3 h-3 rounded-full bg-yellow-400"></div>
                <div class="w-3 h-3 rounded-full bg-green-400"></div>
              </div>
              <div class="h-2 w-32 bg-slate-700 rounded-full"></div>
            </div>

            <!-- Content -->
            <div class="grid grid-cols-2 gap-4 mb-6">
              <div class="bg-slate-700/50 rounded-lg p-4">
                <div class="h-8 w-8 bg-primary/20 rounded-lg flex items-center justify-center mb-2">
                  <IconRecycle class="text-primary" size="20" />
                </div>
                <div class="text-2xl font-bold text-white mb-1">2,543</div>
                <div class="text-xs text-slate-400">Tons Recycled</div>
              </div>
              <div class="bg-slate-700/50 rounded-lg p-4">
                <div class="h-8 w-8 bg-blue-500/20 rounded-lg flex items-center justify-center mb-2">
                  <IconChartBar class="text-blue-400" size="20" />
                </div>
                <div class="text-2xl font-bold text-white mb-1">$45.2k</div>
                <div class="text-xs text-slate-400">Value Generated</div>
              </div>
            </div>

            <!-- Graph placeholder -->
            <div class="bg-slate-700/30 rounded-lg p-4 h-32 flex items-end justify-between gap-2">
              <div class="w-full bg-primary/20 rounded-t-sm h-[40%]"></div>
              <div class="w-full bg-primary/40 rounded-t-sm h-[60%]"></div>
              <div class="w-full bg-primary/30 rounded-t-sm h-[50%]"></div>
              <div class="w-full bg-primary/60 rounded-t-sm h-[80%]"></div>
              <div class="w-full bg-primary/50 rounded-t-sm h-[70%]"></div>
              <div class="w-full bg-primary/80 rounded-t-sm h-[90%]"></div>
            </div>
          </div>

          <!-- Floating Elements -->
          <div
            class="absolute -right-6 top-10 bg-slate-800 p-4 rounded-xl border border-slate-700 shadow-xl animate-bounce"
            style="animation-duration: 3s">
            <div class="flex items-center gap-3">
              <div class="bg-primary/20 p-2 rounded-full">
                <IconCheck class="text-primary" size="20" />
              </div>
              <div>
                <p class="text-xs text-slate-400">Pickup Status</p>
                <p class="text-sm font-bold text-white">Completed</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>



  <!-- Features Grid -->
  <section class="py-24 bg-white">
    <div class="container mx-auto px-6">
      <div class="text-center max-w-3xl mx-auto mb-20">
        <h2 class="text-3xl md:text-5xl font-bold text-slate-900 mb-6">Powerful Features for Modern Waste Management
        </h2>
        <p class="text-lg text-slate-600">Everything you need to manage recycling, track impact, and connect with waste
          collectors.</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-10">
        <div
          class="group p-8 rounded-2xl bg-white border border-slate-100 shadow-lg hover:shadow-2xl transition-all duration-300 hover:-translate-y-2">
          <div
            class="w-14 h-14 bg-primary/10 rounded-xl flex items-center justify-center mb-6 group-hover:bg-primary transition-colors duration-300">
            <IconBuildingWarehouse class="text-primary group-hover:text-white" size="32" />
          </div>
          <h3 class="text-xl font-bold text-slate-900 mb-3">Marketplace</h3>
          <p class="text-slate-600 leading-relaxed">Browse and trade recyclable materials. Connect directly with buyers
            and sellers in your community.</p>
        </div>

        <div
          class="group p-8 rounded-2xl bg-white border border-slate-100 shadow-lg hover:shadow-2xl transition-all duration-300 hover:-translate-y-2">
          <div
            class="w-14 h-14 bg-blue-100 rounded-xl flex items-center justify-center mb-6 group-hover:bg-blue-500 transition-colors duration-300">
            <IconMapPin class="text-blue-600 group-hover:text-white" size="32" />
          </div>
          <h3 class="text-xl font-bold text-slate-900 mb-3">Shop Locator</h3>
          <p class="text-slate-600 leading-relaxed">Find nearby recycling shops easily using our interactive map. View
            shop details, prices, and accepted materials.</p>
        </div>

        <div
          class="group p-8 rounded-2xl bg-white border border-slate-100 shadow-lg hover:shadow-2xl transition-all duration-300 hover:-translate-y-2">
          <div
            class="w-14 h-14 bg-purple-100 rounded-xl flex items-center justify-center mb-6 group-hover:bg-purple-500 transition-colors duration-300">
            <IconBuildingStore class="text-purple-600 group-hover:text-white" size="32" />
          </div>
          <h3 class="text-xl font-bold text-slate-900 mb-3">Shop Management</h3>
          <p class="text-slate-600 leading-relaxed">Shop owners can manage inventory, update pricing, and connect with
            nearby customers for pickups.</p>
        </div>
      </div>
    </div>
  </section>

  <!-- How It Works (Dark) -->
  <section id="how-it-works" class="py-24 bg-slate-900 text-white relative overflow-hidden">
    <!-- Decorative lines -->
    <div
      class="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:40px_40px]">
    </div>

    <div class="container mx-auto px-6 relative z-10">
      <h2 class="text-3xl md:text-5xl font-bold text-center mb-16">How It Works</h2>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8 relative">
        <!-- Connecting Line (Desktop) -->
        <div
          class="hidden md:block absolute top-24 left-[10%] right-[10%] h-0.5 bg-gradient-to-r from-primary via-blue-500 to-purple-500 opacity-30">
        </div>

        <div class="relative text-center">
          <div
            class="w-16 h-16 mx-auto bg-slate-800 border-2 border-primary rounded-full flex items-center justify-center text-2xl font-bold mb-6 hover:bg-primary hover:text-white transition-colors duration-300 z-10 relative">
            1
          </div>
          <h3 class="text-xl font-bold mb-3">Create Profile</h3>
          <p class="text-slate-400 text-sm px-4">Sign up easily and set up your personal or business profile to get
            started.</p>
        </div>

        <div class="relative text-center">
          <div
            class="w-16 h-16 mx-auto bg-slate-800 border-2 border-primary rounded-full flex items-center justify-center text-2xl font-bold mb-6 hover:bg-primary hover:text-white transition-colors duration-300 z-10 relative">
            2
          </div>
          <h3 class="text-xl font-bold mb-3">Locate Shops</h3>
          <p class="text-slate-400 text-sm px-4">Use the Shop Locator map to find participating recycling centers near
            you.</p>
        </div>

        <div class="relative text-center">
          <div
            class="w-16 h-16 mx-auto bg-slate-800 border-2 border-primary rounded-full flex items-center justify-center text-2xl font-bold mb-6 hover:bg-primary hover:text-white transition-colors duration-300 z-10 relative">
            3
          </div>
          <h3 class="text-xl font-bold mb-3">Manage Waste</h3>
          <p class="text-slate-400 text-sm px-4">Shops can manage their inventory, while users can track their recycling
            contributions.</p>
        </div>

        <div class="relative text-center">
          <div
            class="w-16 h-16 mx-auto bg-slate-800 border-2 border-primary rounded-full flex items-center justify-center text-2xl font-bold mb-6 hover:bg-primary hover:text-white transition-colors duration-300 z-10 relative">
            4
          </div>
          <h3 class="text-xl font-bold mb-3">Connect & Trade</h3>
          <p class="text-slate-400 text-sm px-4">Use the Marketplace to buy or sell materials and connect via our
            platform.</p>
        </div>
      </div>
    </div>
  </section>



  <!-- Testimonials Carousel (Simplified Grid for now) -->
  <section class="py-24 bg-white">
    <div class="container mx-auto px-6">
      <h2 class="text-3xl md:text-5xl font-bold text-center text-slate-900 mb-16">Loved by Green Teams</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div v-for="t in testimonials" :key="t.id"
          class="p-8 bg-slate-50 rounded-2xl hover:bg-white hover:shadow-xl transition-all border border-transparent hover:border-primary/20">
          <div class="flex gap-2 mb-4 text-primary">
            <IconStar v-for="i in 5" :key="i" size="18" fill="currentColor" />
          </div>
          <p class="text-slate-700 mb-6 leading-relaxed">"{{ t.content }}"</p>
          <div class="flex items-center gap-4">
            <img :src="t.avatar" class="w-12 h-12 rounded-full ring-2 ring-primary/20" :alt="t.name">
            <div>
              <h4 class="font-bold text-slate-900">{{ t.name }}</h4>
              <p class="text-sm text-slate-500">{{ t.role }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- CTA Section -->
  <section class="py-24 bg-gradient-to-r from-accent to-teal-600 relative overflow-hidden">
    <div class="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20"></div>
    <div class="container mx-auto px-6 text-center relative z-10 text-white pb-16">
      <h2 class="text-3xl md:text-5xl font-bold mb-6">Ready to Optimize Your Waste Management?</h2>
      <p class="text-xl text-primary-content mb-10 max-w-2xl mx-auto">Join the platform that's powering the circular
        economy. Start for free today.</p>
      <button
        class="btn bg-white text-primary hover:bg-primary/10 border-0 px-10 py-4 h-auto text-lg rounded-full font-bold shadow-2xl"
        @click="handleGetStarted">
        Get Started Now
      </button>
    </div>

    <!-- Wave Divider -->
    <div class="absolute bottom-0 w-full overflow-hidden leading-[0] transform translate-y-[1px]">
      <svg class="relative block w-[calc(140%+1.3px)] h-[50px] md:h-[100px] fill-slate-900" data-name="Layer 1"
        xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 120" preserveAspectRatio="none">
        <path
          d="M985.66,92.83C906.67,72,823.78,31,743.84,14.19c-82.26-17.34-168.06-16.33-250.45.39-57.84,11.73-114,31.07-172,41.86A600.21,600.21,0,0,1,0,27.35V120H1200V95.8C1132.19,118.92,1055.71,111.31,985.66,92.83Z">
        </path>
      </svg>
    </div>
  </section>
</template>

<style scoped>
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.hero-content-anim {
  opacity: 0;
  animation: fadeInUp 0.8s ease-out forwards;
}

@keyframes float {
  0% {
    transform: translateY(0px);
  }

  50% {
    transform: translateY(-10px);
  }

  100% {
    transform: translateY(0px);
  }
}

.floating-dashboard {
  animation: float 6s ease-in-out infinite;
}

/* Custom Animation for Gradient Blobs */
@keyframes blob {
  0% {
    transform: translate(0px, 0px) scale(1);
  }

  33% {
    transform: translate(30px, -50px) scale(1.1);
  }

  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }

  100% {
    transform: translate(0px, 0px) scale(1);
  }
}

.animate-blob {
  animation: blob 7s infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}
</style>

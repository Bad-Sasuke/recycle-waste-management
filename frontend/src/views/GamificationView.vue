<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { IconLeaf, IconTrophy, IconMedal, IconShare, IconTrendingUp } from '@tabler/icons-vue'

// --- Mock Data ---
const userStats = ref({
    totalPoints: 1250,
    currentRank: 'Sapling Guardian',
    nextRank: 'Forest Ranger',
    pointsToNextRank: 250,
    treesPlanted: 5,
    co2Saved: 120.5, // kg
})

const recentActivities = ref([
    { id: 1, action: 'Recycled Plastic Bottles', points: +50, date: '2 hours ago', icon: 'bottle' },
    { id: 2, action: 'Recycled Paper', points: +30, date: '1 day ago', icon: 'paper' },
    { id: 3, action: 'Daily Login Bonus', points: +10, date: '1 day ago', icon: 'star' },
])

const leaderboard = ref([
    { rank: 1, name: 'EcoWarrior99', points: 5400, avatar: 'https://ui-avatars.com/api/?name=Eco+Warrior&background=10B981&color=fff' },
    { rank: 2, name: 'GreenLife', points: 4850, avatar: 'https://ui-avatars.com/api/?name=Green+Life&background=3B82F6&color=fff' },
    { rank: 3, name: 'RecycleKing', points: 4200, avatar: 'https://ui-avatars.com/api/?name=Recycle+King&background=F59E0B&color=fff' },
    { rank: 4, name: 'Sarah J.', points: 3900, avatar: 'https://ui-avatars.com/api/?name=Sarah+J&background=EC4899&color=fff' },
    { rank: 5, name: 'Mike Chen', points: 3500, avatar: 'https://ui-avatars.com/api/?name=Mike+Chen&background=6366F1&color=fff' },
])

// --- Computed ---
const progressPercentage = computed(() => {
    return 85 // Mock 85% for demo
})

// --- Animation State ---
const treeScale = ref(0.8)
const animateTree = () => {
    treeScale.value = 1
}

onMounted(() => {
    setTimeout(animateTree, 500)
})
</script>

<template>
    <div class="min-h-screen bg-slate-50 pb-20 pt-10">
        <!-- Header Hero -->
        <div class="bg-emerald-600 text-white py-12 mb-8 relative overflow-hidden">
            <div class="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20"></div>
            <div class="container mx-auto px-6 relative z-10">
                <div class="flex flex-col md:flex-row items-center justify-between gap-6">
                    <div>
                        <h1 class="text-3xl md:text-5xl font-bold mb-2">My Eco Forest 🌳</h1>
                        <p class="text-emerald-100 text-lg">Your recycling efforts are growing cloud forests!</p>
                    </div>
                    <div
                        class="flex items-center gap-4 bg-white/10 backdrop-blur-md p-4 rounded-xl border border-white/20">
                        <div class="bg-yellow-400 p-3 rounded-full text-yellow-900">
                            <IconTrophy size="32" stroke="2" />
                        </div>
                        <div>
                            <div class="text-xs uppercase tracking-wider font-semibold opacity-80">Current Rank</div>
                            <div class="text-2xl font-bold">{{ userStats.currentRank }}</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="container mx-auto px-6">
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">

                <!-- Left Column: My Forest & Stats -->
                <div class="lg:col-span-2 space-y-8">

                    <!-- Tree Visualization -->
                    <div class="bg-white rounded-3xl shadow-xl overflow-hidden relative border border-slate-100">
                        <div class="p-6 border-b border-slate-100 flex justify-between items-center">
                            <h2 class="text-xl font-bold text-slate-800 flex items-center gap-2">
                                <IconLeaf class="text-emerald-500" /> My Virtual Tree
                            </h2>
                            <span class="bg-emerald-100 text-emerald-700 text-xs font-bold px-3 py-1 rounded-full">Level
                                12</span>
                        </div>

                        <div
                            class="h-96 relative bg-gradient-to-b from-sky-100 to-emerald-50 flex items-end justify-center overflow-hidden group cursor-pointer">
                            <!-- Animated Sun -->
                            <div
                                class="absolute top-10 right-10 w-16 h-16 bg-yellow-300 rounded-full blur-xl opacity-60 animate-pulse">
                            </div>

                            <!-- The Tree SVG -->
                            <svg viewBox="0 0 400 400"
                                class="w-80 h-80 transition-transform duration-1000 ease-out drop-shadow-2xl z-10"
                                :style="{ transform: `scale(${treeScale})` }">
                                <defs>
                                    <linearGradient id="trunkGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                                        <stop offset="0%" style="stop-color:#8B4513;stop-opacity:1" />
                                        <stop offset="50%" style="stop-color:#A0522D;stop-opacity:1" />
                                        <stop offset="100%" style="stop-color:#8B4513;stop-opacity:1" />
                                    </linearGradient>
                                    <linearGradient id="leafGrad1" x1="0%" y1="0%" x2="0%" y2="100%">
                                        <stop offset="0%" style="stop-color:#4ADE80;stop-opacity:1" />
                                        <stop offset="100%" style="stop-color:#22C55E;stop-opacity:1" />
                                    </linearGradient>
                                    <linearGradient id="leafGrad2" x1="0%" y1="0%" x2="0%" y2="100%">
                                        <stop offset="0%" style="stop-color:#34D399;stop-opacity:1" />
                                        <stop offset="100%" style="stop-color:#059669;stop-opacity:1" />
                                    </linearGradient>
                                    <filter id="glow" x="-20%" y="-20%" width="140%" height="140%">
                                        <feGaussianBlur stdDeviation="5" result="blur" />
                                        <feComposite in="SourceGraphic" in2="blur" operator="over" />
                                    </filter>
                                </defs>

                                <!-- Ground Shadow -->
                                <ellipse cx="200" cy="360" rx="80" ry="15" fill="#000000" opacity="0.1" />

                                <!-- Trunk -->
                                <path
                                    d="M185 360 Q180 300 190 250 Q160 220 150 180 M215 360 Q220 300 210 250 Q240 220 250 180"
                                    stroke="url(#trunkGrad)" stroke-width="24" stroke-linecap="round" fill="none" />
                                <path d="M200 260 L200 160" stroke="url(#trunkGrad)" stroke-width="20"
                                    stroke-linecap="round" />

                                <!-- Foliage Group (Bottom Layer) -->
                                <g class="animate-sway" style="animation-duration: 5s; transform-origin: 200px 200px;">
                                    <circle cx="150" cy="190" r="50" fill="url(#leafGrad2)" />
                                    <circle cx="250" cy="190" r="50" fill="url(#leafGrad2)" />
                                    <circle cx="200" cy="170" r="55" fill="url(#leafGrad2)" />
                                </g>

                                <!-- Foliage Group (Top Layer - Lighter) -->
                                <g class="animate-sway"
                                    style="animation-duration: 4s; animation-delay: 0.5s; transform-origin: 200px 150px;">
                                    <circle cx="175" cy="140" r="55" fill="url(#leafGrad1)" />
                                    <circle cx="225" cy="140" r="55" fill="url(#leafGrad1)" />
                                    <circle cx="200" cy="100" r="60" fill="url(#leafGrad1)" />
                                </g>

                                <!-- Shine/Highlight -->
                                <circle cx="180" cy="120" r="15" fill="white" opacity="0.2" />
                            </svg>

                            <!-- Floating Points -->
                            <div
                                class="absolute top-1/2 left-1/2 -translate-x-1/2 -mt-20 opacity-0 group-hover:opacity-100 group-hover:-translate-y-20 transition-all duration-700 pointer-events-none">
                                <div
                                    class="bg-white/90 backdrop-blur text-emerald-600 px-4 py-2 rounded-xl shadow-lg font-bold border border-emerald-100 flex items-center gap-2">
                                    <IconLeaf size="16" class="fill-current" />
                                    +12 Oxygen
                                </div>
                            </div>
                        </div>

                        <div class="p-6 bg-white relative z-10">
                            <div class="flex justify-between text-sm mb-2 font-medium text-slate-600">
                                <span>Progress to {{ userStats.nextRank }}</span>
                                <span>{{ userStats.pointsToNextRank }} pts needed</span>
                            </div>
                            <div class="w-full bg-slate-100 rounded-full h-4 overflow-hidden">
                                <div class="bg-emerald-500 h-4 rounded-full transition-all duration-1000 relative overflow-hidden"
                                    :style="{ width: `${progressPercentage}%` }">
                                    <div class="absolute inset-0 bg-white/20 animate-[shimmer_2s_infinite]"></div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Statistics Cards -->
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                        <div
                            class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 text-center hover:shadow-md transition-shadow">
                            <div
                                class="bg-blue-100 w-12 h-12 rounded-xl flex items-center justify-center mx-auto mb-4 text-blue-600">
                                <IconTrendingUp size="24" />
                            </div>
                            <div class="text-3xl font-bold text-slate-800">{{ userStats.totalPoints }}</div>
                            <div class="text-sm text-slate-500">Total Points</div>
                        </div>
                        <div
                            class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 text-center hover:shadow-md transition-shadow">
                            <div
                                class="bg-green-100 w-12 h-12 rounded-xl flex items-center justify-center mx-auto mb-4 text-green-600">
                                <IconLeaf size="24" />
                            </div>
                            <div class="text-3xl font-bold text-slate-800">{{ userStats.treesPlanted }}</div>
                            <div class="text-sm text-slate-500">Trees Planted</div>
                        </div>
                        <div
                            class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 text-center hover:shadow-md transition-shadow">
                            <div
                                class="bg-teal-100 w-12 h-12 rounded-xl flex items-center justify-center mx-auto mb-4 text-teal-600">
                                <IconShare size="24" />
                            </div>
                            <div class="text-3xl font-bold text-slate-800">{{ userStats.co2Saved }}</div>
                            <div class="text-sm text-slate-500">kg CO₂ Saved</div>
                        </div>
                    </div>

                </div>

                <!-- Right Column: Leaderboard -->
                <div class="space-y-8">
                    <!-- Leaderboard Card -->
                    <div class="bg-white p-6 rounded-3xl shadow-lg border border-slate-100">
                        <h2 class="text-xl font-bold mb-6 flex items-center gap-2">
                            <IconMedal class="text-yellow-500" />
                            Top Green Heroes
                        </h2>

                        <div class="space-y-4">
                            <div v-for="(user, index) in leaderboard" :key="index"
                                class="flex items-center gap-4 p-3 rounded-xl transition-colors hover:bg-slate-50"
                                :class="index === 0 ? 'bg-yellow-50 border border-yellow-200' : ''">
                                <div class="font-bold w-6 text-center"
                                    :class="index === 0 ? 'text-yellow-600 text-xl' : 'text-slate-400'">
                                    {{ user.rank }}
                                </div>
                                <img :src="user.avatar" class="w-10 h-10 rounded-full ring-2 ring-white shadow-sm">
                                <div class="flex-1">
                                    <div class="font-bold text-slate-800">{{ user.name }}</div>
                                    <div class="text-xs text-slate-500">{{ user.points }} pts</div>
                                </div>
                                <div v-if="index < 3">
                                    <IconTrophy v-if="index === 0" size="20" class="text-yellow-500" />
                                    <IconMedal v-else size="20" class="text-slate-300" />
                                </div>
                            </div>
                        </div>

                        <button
                            class="w-full mt-6 btn btn-outline border-slate-200 hover:bg-slate-50 hover:border-slate-300 text-slate-600">
                            View Full Leaderboard
                        </button>
                    </div>

                    <!-- Recent Activity Card -->
                    <div class="bg-white p-6 rounded-3xl shadow-sm border border-slate-100">
                        <h2 class="text-lg font-bold mb-4">Recent Activity</h2>
                        <div class="relative pl-6 border-l-2 border-slate-100 space-y-6">
                            <div v-for="activity in recentActivities" :key="activity.id" class="relative">
                                <div
                                    class="absolute -left-[31px] bg-emerald-100 p-1.5 rounded-full border-4 border-white">
                                    <div class="w-2 h-2 bg-emerald-500 rounded-full"></div>
                                </div>
                                <div class="text-sm font-semibold text-slate-800">{{ activity.action }}</div>
                                <div class="text-xs text-slate-500 flex justify-between mt-1">
                                    <span>{{ activity.date }}</span>
                                    <span class="text-emerald-600 font-bold">+{{ activity.points }} pts</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
@keyframes shimmer {
    0% {
        transform: translateX(-100%);
    }

    100% {
        transform: translateX(100%);
    }
}

@keyframes sway {

    0%,
    100% {
        transform: rotate(-3deg);
    }

    50% {
        transform: rotate(3deg);
    }
}

.animate-sway {
    animation: sway ease-in-out infinite;
    transform-box: fill-box;
}
</style>

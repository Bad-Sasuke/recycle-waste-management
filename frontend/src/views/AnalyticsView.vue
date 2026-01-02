<template>
    <div class="container mx-auto px-4 py-8 max-w-6xl">
        <div class="mb-8">
            <h1 class="text-3xl font-bold bg-gradient-to-r from-green-600 to-teal-500 bg-clip-text text-transparent">
                My Recycling Analytics 📉
            </h1>
            <p class="text-gray-500 dark:text-gray-400 mt-2">
                Detailed insights into your recycling habits and earnings.
            </p>
        </div>

        <div v-if="loading" class="flex justify-center items-center h-64">
            <span class="loading loading-spinner loading-lg text-success"></span>
        </div>

        <div v-else-if="error" class="alert alert-error">
            <span>{{ error }}</span>
        </div>

        <div v-else>
            <!-- Summary Cards -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                <!-- Earnings -->
                <div
                    class="stats shadow bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 border border-green-100 dark:border-green-800">
                    <div class="stat">
                        <div class="stat-figure text-green-600 dark:text-green-400">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                class="inline-block w-8 h-8 stroke-current">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z">
                                </path>
                            </svg>
                        </div>
                        <div class="stat-title text-gray-600 dark:text-gray-400">Total Earnings (This Month)</div>
                        <div class="stat-value text-green-600 dark:text-green-400">฿{{
                            formatNumber(analytics?.total_earnings_this_month || 0) }}</div>
                        <div class="stat-desc">From recycling sales</div>
                    </div>
                </div>

                <!-- Weight -->
                <div class="stats shadow bg-base-100 dark:bg-base-200 border border-base-200 dark:border-base-700">
                    <div class="stat">
                        <div class="stat-figure text-secondary">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                class="inline-block w-8 h-8 stroke-current">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16H9m3 0h3">
                                </path>
                            </svg>
                        </div>
                        <div class="stat-title">Recycled Weight (This Month)</div>
                        <div class="stat-value text-secondary">{{ formatNumber(analytics?.total_weight_this_month || 0)
                            }} kg</div>
                        <div class="stat-desc" :class="getDiffColorClass(analytics?.weight_diff_percent)">
                            <span v-if="(analytics?.weight_diff_percent || 0) > 0">↗︎ </span>
                            <span v-else-if="(analytics?.weight_diff_percent || 0) < 0">↘︎ </span>
                            {{ Math.abs(analytics?.weight_diff_percent || 0).toFixed(1) }}% vs Last Month
                        </div>
                    </div>
                </div>

                <!-- Last Month Comparison -->
                <div class="stats shadow bg-base-100 dark:bg-base-200 border border-base-200 dark:border-base-700">
                    <div class="stat">
                        <div class="stat-figure text-primary">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                class="inline-block w-8 h-8 stroke-current">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z">
                                </path>
                            </svg>
                        </div>
                        <div class="stat-title">Last Month Total</div>
                        <div class="stat-value text-primary">{{ formatNumber(analytics?.total_weight_last_month || 0) }}
                            kg</div>
                        <div class="stat-desc">Previous month performance</div>
                    </div>
                </div>
            </div>

            <!-- Charts Section -->
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <!-- Category Breakdown -->
                <div class="card bg-base-100 shadow-xl border border-base-200">
                    <div class="card-body">
                        <h2 class="card-title text-lg mb-4">♻️ Waste Category Breakdown</h2>
                        <div class="h-64 relative">
                            <Doughnut v-if="chartDataCategory" Conversations :data="chartDataCategory"
                                :options="chartOptions" />
                            <div v-else class="flex items-center justify-center h-full text-gray-400">No data available
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Daily Trend -->
                <div class="card bg-base-100 shadow-xl border border-base-200">
                    <div class="card-body">
                        <h2 class="card-title text-lg mb-4">📅 Daily Activity (This Month)</h2>
                        <div class="h-64 relative">
                            <Bar v-if="chartDataDaily" :data="chartDataDaily" :options="chartOptions" />
                            <div v-else class="flex items-center justify-center h-full text-gray-400">No data available
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { receiptService } from '@/services/receipt';
import type { UserAnalyticsResponse } from '@/types/receipt';
import {
    Chart as ChartJS,
    ArcElement,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale,
    Title
} from 'chart.js';
import { Doughnut, Bar } from 'vue-chartjs';

ChartJS.register(ArcElement, Tooltip, Legend, BarElement, CategoryScale, LinearScale, Title);

const analytics = ref<UserAnalyticsResponse | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

const fetchAnalytics = async () => {
    try {
        loading.value = true;
        analytics.value = await receiptService.getMyAnalytics();
    } catch (e: any) {
        error.value = e.message || 'Failed to load analytics';
    } finally {
        loading.value = false;
    }
};

const formatNumber = (num: number) => {
    return new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(num);
};

const getDiffColorClass = (percent: number | undefined) => {
    if (!percent) return 'text-gray-500';
    return percent > 0 ? 'text-success' : percent < 0 ? 'text-error' : 'text-gray-500';
};

// Chart Data
const chartDataCategory = computed(() => {
    if (!analytics.value || !analytics.value.monthly_category_stats.length) return null;

    const labels = analytics.value.monthly_category_stats.map(s => s.category);
    const data = analytics.value.monthly_category_stats.map(s => s.weight);

    // Vibrant Colors
    const backgroundColors = [
        '#34d399', // Emerald
        '#60a5fa', // Blue
        '#f472b6', // Pink
        '#fbbf24', // Amber
        '#a78bfa', // Purple
        '#2dd4bf', // Teal
    ];

    return {
        labels,
        datasets: [{
            backgroundColor: backgroundColors,
            data
        }]
    };
});

const chartDataDaily = computed(() => {
    if (!analytics.value || !analytics.value.daily_stats.length) return null;

    // Sort by date
    const sortedStats = [...analytics.value.daily_stats].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());

    const labels = sortedStats.map(s => {
        const date = new Date(s.date);
        return date.toLocaleDateString('en-US', { day: 'numeric', month: 'short' });
    });
    const data = sortedStats.map(s => s.weight);

    return {
        labels,
        datasets: [{
            label: 'Weight (kg)',
            backgroundColor: '#10b981',
            data
        }]
    };
});

const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: {
            position: 'bottom' as const,
        }
    }
};

onMounted(() => {
    fetchAnalytics();
});
</script>

<style scoped>
/* Gradient Text Utility */
.bg-clip-text {
    -webkit-background-clip: text;
}
</style>

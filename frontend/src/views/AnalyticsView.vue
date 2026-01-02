<template>
  <div class="container mx-auto px-4 py-8 max-w-6xl">
    <div class="mb-8">
      <h1
        class="text-3xl font-bold bg-gradient-to-r from-green-600 to-teal-500 bg-clip-text text-transparent flex items-center">
        My Recycling Analytics
        <IconChartLine size="32" class="ml-2 text-teal-600" />
      </h1>
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mt-2">
        <p class="text-gray-500 dark:text-gray-400">
          Detailed insights into your recycling habits and earnings.
        </p>
        <button @click="showImpactStory = true"
          class="btn btn-sm btn-primary text-white gap-2 shadow-sm hover:shadow-md transition-all">
          <IconShare size="16" />
          Share Impact
        </button>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <span class="loading loading-spinner loading-lg text-success"></span>
    </div>

    <div v-else-if="error" class="alert alert-error">
      <span>{{ error }}</span>
    </div>

    <div v-else>
      <!-- Empty State -->
      <div v-if="!analytics?.total_weight_this_month && !analytics?.total_weight_last_month"
        class="text-center py-16 bg-white rounded-xl shadow-sm border border-gray-100">
        <div class="inline-flex items-center justify-center w-24 h-24 rounded-full bg-green-50 mb-6">
          <IconSeeding size="48" class="text-green-500" stroke="1.5" />
        </div>
        <h2 class="text-2xl font-bold text-gray-800 mb-2">Start Your Recycling Journey!</h2>
        <p class="text-gray-500 max-w-md mx-auto mb-8">
          It looks like you haven't completed any recycling requests yet.
          Start selling your waste to see your impact stats here!
        </p>
        <router-link to="/waste-purchase" class="btn btn-primary text-white">
          Create Sell Request
        </router-link>
      </div>

      <!-- Summary Cards -->
      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <!-- Earnings -->
        <div
          class="stats shadow bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 border border-green-100 dark:border-green-800">
          <div class="stat">
            <div class="stat-figure text-green-600 dark:text-green-400">
              <IconCurrencyBaht size="32" stroke="1.5" />
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
              <IconScale size="32" stroke="1.5" />
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
              <IconClock size="32" stroke="1.5" />
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
            <h2 class="card-title text-lg mb-4 flex items-center">
              <IconRecycle size="20" class="mr-2 text-green-600" />
              Waste Category Breakdown
            </h2>
            <div class="h-64 relative">
              <Doughnut v-if="chartDataCategory" :data="chartDataCategory" :options="chartOptions" />
              <div v-else class="flex items-center justify-center h-full text-gray-400">No data available
              </div>
            </div>
          </div>
        </div>

        <!-- Daily Trend -->
        <div class="card bg-base-100 shadow-xl border border-base-200">
          <div class="card-body">
            <h2 class="card-title text-lg mb-4 flex items-center">
              <IconCalendar size="20" class="mr-2 text-blue-500" />
              Daily Activity (This Month)
            </h2>
            <div class="h-64 relative">
              <Bar v-if="chartDataDaily" :data="chartDataDaily" :options="chartOptions" />
              <div v-else class="flex items-center justify-center h-full text-gray-400">No data available
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Impact Story Modal -->
    <ImpactStory v-if="showImpactStory && analytics" :data="analytics" @close="showImpactStory = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { receiptService } from '@/services/receipt';
import type { UserAnalyticsResponse } from '@/types/receipt';
import ImpactStory from '@/components/ImpactStory.vue';
import {
  IconChartLine,
  IconShare,
  IconSeeding,
  IconRecycle,
  IconCalendar,
  IconCurrencyBaht,
  IconScale,
  IconClock
} from '@tabler/icons-vue';
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
const showImpactStory = ref(false);

const fetchAnalytics = async () => {
  try {
    loading.value = true;
    analytics.value = await receiptService.getMyAnalytics();
  } catch (e: unknown) {
    const message = e instanceof Error ? e.message : 'Failed to load analytics';
    error.value = message;
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
  const stats = analytics.value?.monthly_category_stats || [];
  if (!stats.length) return null;

  const labels = stats.map(s => s.category);
  const data = stats.map(s => s.weight);

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
  const stats = analytics.value?.daily_stats || [];
  if (!stats.length) return null;

  // Sort by date
  const sortedStats = [...stats].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());

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
  background-clip: text;
  -webkit-background-clip: text;
}
</style>

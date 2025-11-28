<template>
    <div class="min-h-screen bg-gradient-to-br from-base-200 to-base-300">
        <!-- Loading State -->
        <div v-if="isLoading" class="flex justify-center items-center min-h-screen">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>

        <!-- Content -->
        <div v-else-if="shop" class="container mx-auto px-4 py-8 max-w-7xl">
            <!-- Header Section -->
            <div class="card bg-base-100 shadow-2xl mb-8">
                <div class="card-body">
                    <div class="flex flex-col md:flex-row gap-6">
                        <!-- Shop Logo -->
                        <div class="flex-shrink-0">
                            <div class="avatar">
                                <div class="w-32 h-32 rounded-xl ring ring-primary ring-offset-base-100 ring-offset-2">
                                    <img v-if="shop.image_url" :src="shop.image_url" :alt="shop.name"
                                        class="w-full h-full object-cover" />
                                    <div v-else
                                        class="w-full h-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-4xl font-bold text-white">
                                        {{ shop.name.charAt(0) }}
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- Shop Info -->
                        <div class="flex-grow">
                            <h1 class="text-3xl font-bold text-gray-800 mb-2">{{ shop.name }}</h1>
                            <p v-if="shop.description" class="text-gray-600 mb-4">{{ shop.description }}</p>

                            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                                <div class="flex items-start gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg"
                                        class="w-5 h-5 text-primary flex-shrink-0 mt-0.5" fill="none"
                                        viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                                    </svg>
                                    <span class="text-sm text-gray-700">{{ shop.address }}</span>
                                </div>

                                <div v-if="shop.phone" class="flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-primary" fill="none"
                                        viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                                    </svg>
                                    <span class="text-sm text-gray-700">{{ shop.phone }}</span>
                                </div>

                                <div v-if="shop.email" class="flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-primary" fill="none"
                                        viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                                    </svg>
                                    <span class="text-sm text-gray-700">{{ shop.email }}</span>
                                </div>

                                <div v-if="shop.opening_time && shop.closing_time" class="flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-primary" fill="none"
                                        viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                    <span class="text-sm text-gray-700">{{ shop.opening_time }} - {{ shop.closing_time
                                        }}</span>
                                </div>
                            </div>
                        </div>

                        <!-- Stats Card -->
                        <div class="flex-shrink-0">
                            <div
                                class="stats stats-vertical shadow bg-gradient-to-br from-primary to-secondary text-primary-content">
                                <div class="stat place-items-center">
                                    <div class="stat-title text-white/80">Rating</div>
                                    <div class="stat-value text-white">{{ averageRating.toFixed(1) }}</div>
                                    <div class="stat-desc text-white/70">
                                        <div class="rating rating-sm gap-1">
                                            <input v-for="i in 5" :key="i" type="radio" name="shop-rating"
                                                class="mask mask-star-2 bg-white"
                                                :checked="i === Math.round(averageRating)" disabled />
                                        </div>
                                    </div>
                                </div>

                                <div class="stat place-items-center">
                                    <div class="stat-title text-white/80">Reviews</div>
                                    <div class="stat-value text-white">{{ totalReviews }}</div>
                                    <div class="stat-desc text-white/70">Total feedback</div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Reviews Section -->
            <div class="card bg-base-100 shadow-xl">
                <div class="card-body">
                    <h2 class="card-title text-2xl mb-6">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-primary" fill="none"
                            viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                        </svg>
                        Customer Reviews
                    </h2>

                    <!-- No Reviews -->
                    <div v-if="reviews.length === 0" class="text-center py-12">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none"
                            viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
                        </svg>
                        <p class="text-gray-500">No reviews yet</p>
                    </div>

                    <!-- Reviews List -->
                    <div v-else class="space-y-4">
                        <div v-for="review in reviews" :key="review.review_id" class="card bg-base-200 shadow">
                            <div class="card-body">
                                <div class="flex items-start justify-between">
                                    <div class="flex items-center gap-3">
                                        <div class="avatar placeholder">
                                            <div class="bg-neutral text-neutral-content rounded-full w-10">
                                                <span class="text-sm">{{ review.user_id.substring(0, 2).toUpperCase()
                                                    }}</span>
                                            </div>
                                        </div>
                                        <div>
                                            <p class="font-semibold text-gray-800">Customer Review</p>
                                            <p class="text-xs text-gray-500">{{ formatDate(review.created_at) }}</p>
                                        </div>
                                    </div>
                                    <div class="rating rating-sm">
                                        <input v-for="i in 5" :key="i" type="radio" :name="'rating-' + review.review_id"
                                            class="mask mask-star-2 bg-warning" :checked="i === review.rating"
                                            disabled />
                                    </div>
                                </div>
                                <p v-if="review.comment" class="text-gray-700 mt-3">{{ review.comment }}</p>
                            </div>
                        </div>
                    </div>

                    <!-- Pagination -->
                    <div v-if="totalPages > 1" class="flex justify-center mt-6">
                        <div class="join">
                            <button class="join-item btn btn-sm" :disabled="currentPage === 1"
                                @click="changePage(currentPage - 1)">
                                «
                            </button>

                            <button v-for="page in displayPages" :key="page" class="join-item btn btn-sm"
                                :class="{ 'btn-active': currentPage === page }" @click="changePage(page)">
                                {{ page }}
                            </button>

                            <button class="join-item btn btn-sm" :disabled="currentPage === totalPages"
                                @click="changePage(currentPage + 1)">
                                »
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Error State -->
        <div v-else class="flex flex-col justify-center items-center min-h-screen">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 text-error mb-4" fill="none" viewBox="0 0 24 24"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <p class="text-xl text-gray-700">Shop not found</p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import config from '../config';
import { fetchShopById } from '@/services/shop';

const route = useRoute();
const webAPI = config.webAPI;

const shop = ref<any>(null);
const reviews = ref<any[]>([]);
const isLoading = ref(true);
const currentPage = ref(1);
const pageSize = 10;
const totalReviews = ref(0);
const totalPages = ref(0);

const averageRating = computed(() => {
    if (shop.value && shop.value.average_rating !== undefined) {
        return shop.value.average_rating;
    }
    return 0;
});

const displayPages = computed(() => {
    const pages = [];
    const maxDisplay = 5;
    let start = Math.max(1, currentPage.value - Math.floor(maxDisplay / 2));
    let end = Math.min(totalPages.value, start + maxDisplay - 1);

    if (end - start + 1 < maxDisplay) {
        start = Math.max(1, end - maxDisplay + 1);
    }

    for (let i = start; i <= end; i++) {
        pages.push(i);
    }
    return pages;
});

const fetchShop = async () => {
    try {
        const shopId = route.params.shop_id as string;
        const response = await fetchShopById(shopId);

        if (response.data) {
            shop.value = response.data;
            // Use total reviews from shop data if available
            if (shop.value.total_reviews !== undefined) {
                totalReviews.value = shop.value.total_reviews;
            }
        }
    } catch (error) {
        console.error('Error fetching shop:', error);
    }
};

const fetchReviews = async () => {
    try {
        const shopId = route.params.shop_id as string;
        const response = await fetch(`${webAPI}/api/reviews/shop/${shopId}?page=${currentPage.value}&page_size=${pageSize}`);
        const data = await response.json();

        if (data.reviews) {
            reviews.value = data.reviews;
            totalReviews.value = data.total_count;
            totalPages.value = data.total_pages;
        }
    } catch (error) {
        console.error('Error fetching reviews:', error);
    }
};

const changePage = async (page: number) => {
    currentPage.value = page;
    await fetchReviews();
    window.scrollTo({ top: 0, behavior: 'smooth' });
};

const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('th-TH', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
};

onMounted(async () => {
    isLoading.value = true;
    await Promise.all([fetchShop(), fetchReviews()]);
    isLoading.value = false;
});
</script>

<script setup lang="ts">
defineProps({
    href: {
        type: String,
        default: null
    },
    target: {
        type: String,
        default: '_self'
    },
    variant: {
        type: String,
        default: 'primary' // primary, warning, danger, info
    }
})

const getGradientClass = (variant: string) => {
    switch (variant) {
        case 'warning': return 'bg-gradient-warning shadow-orange'
        case 'danger': return 'bg-gradient-danger shadow-red'
        case 'info': return 'bg-gradient-info shadow-blue'
        default: return 'bg-gradient-primary shadow-green' // Assuming primary is green-ish or default
    }
}
</script>

<template>
    <component :is="href ? 'a' : 'button'" :href="href" :target="target" class="navigate-btn"
        :class="getGradientClass(variant)">
        <slot></slot>
    </component>
</template>

<style scoped>
.navigate-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    width: 100%;
    padding: 10px 18px;
    color: white;
    font-weight: 600;
    font-size: 0.875rem;
    border-radius: 10px;
    text-decoration: none;
    margin-top: 6px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    border: none;
    cursor: pointer;
    position: relative;
    overflow: hidden;
}

.navigate-btn::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
    transition: left 0.5s;
}

.navigate-btn:hover {
    transform: translateY(-2px);
}

.navigate-btn:hover::before {
    left: 100%;
}

.navigate-btn:active {
    transform: translateY(0);
}
</style>

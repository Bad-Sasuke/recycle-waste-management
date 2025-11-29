<script setup lang="ts">
defineProps({
    colorClass: {
        type: String,
        default: 'bg-primary'
    }
})
</script>

<template>
    <div class="liquid-marker">
        <div class="liquid-pin" :class="colorClass">
            <div class="liquid-glow"></div>
            <slot>
                <!-- Default Icon -->
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" width="20" height="20"
                    style="position: relative; z-index: 2;">
                    <path
                        d="M21 13.242V20H22V22H2V20H3V13.242L1.515 9.583C1.195 8.772 1.633 7.862 2.444 7.541C2.626 7.466 2.821 7.428 3.018 7.428H20.982C21.889 7.428 22.623 8.163 22.623 9.07C22.623 9.267 22.585 9.462 22.51 9.644L21 13.242ZM19 13.972L20.234 10.428H3.766L5 13.972V20H19V13.972ZM14 10H10V12C10 13.105 10.895 14 12 14C13.105 14 14 13.105 14 12V10Z" />
                </svg>
            </slot>
        </div>
        <div class="liquid-pulse" :class="colorClass"></div>
    </div>
</template>

<style scoped>
.liquid-marker {
    position: relative;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.liquid-pin {
    width: 40px;
    height: 40px;
    border-radius: 50% 50% 50% 0;
    transform: rotate(-45deg);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow:
        0 4px 12px rgba(0, 0, 0, 0.2),
        inset 0 2px 4px rgba(255, 255, 255, 0.3),
        inset 0 -2px 4px rgba(0, 0, 0, 0.1);
    position: relative;
    z-index: 2;
    animation: float 3s ease-in-out infinite;
}

.liquid-pin :deep(svg) {
    transform: rotate(45deg);
}

.liquid-glow {
    position: absolute;
    top: 5px;
    left: 5px;
    width: 15px;
    height: 15px;
    background: radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.8), transparent);
    border-radius: 50%;
    z-index: 3;
}

.liquid-pulse {
    position: absolute;
    bottom: -10px;
    left: 50%;
    transform: translateX(-50%) rotateX(60deg);
    width: 30px;
    height: 30px;
    border-radius: 50%;
    opacity: 0.5;
    filter: blur(4px);
    animation: pulse-shadow 3s ease-in-out infinite;
    z-index: 1;
}

@keyframes float {

    0%,
    100% {
        transform: translateY(0) rotate(-45deg);
    }

    50% {
        transform: translateY(-6px) rotate(-45deg);
    }
}

@keyframes pulse-shadow {

    0%,
    100% {
        transform: translateX(-50%) rotateX(60deg) scale(1);
        opacity: 0.5;
    }

    50% {
        transform: translateX(-50%) rotateX(60deg) scale(0.8);
        opacity: 0.2;
    }
}
</style>

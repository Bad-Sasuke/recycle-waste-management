<template>
    <div class="glass-card">
        <div class="glass-card-content">
            <slot></slot>
        </div>
    </div>
</template>

<style scoped>
.glass-card {
    background: rgba(255, 255, 255, 0.65);
    backdrop-filter: blur(20px) saturate(180%);
    -webkit-backdrop-filter: blur(20px) saturate(180%);
    border-radius: 20px;
    border: 1px solid rgba(255, 255, 255, 0.3);
    box-shadow:
        0 8px 32px 0 rgba(31, 38, 135, 0.15),
        0 0 0 1px rgba(255, 255, 255, 0.1) inset,
        0 1px 2px 0 rgba(0, 0, 0, 0.05);
    overflow: hidden;
    position: relative;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.glass-card:hover {
    transform: translateY(-2px);
    box-shadow:
        0 12px 40px 0 rgba(31, 38, 135, 0.2),
        0 0 0 1px rgba(255, 255, 255, 0.2) inset;
}

.glass-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: -50%;
    width: 200%;
    height: 100%;
    background: linear-gradient(to right,
            transparent,
            rgba(255, 255, 255, 0.3),
            transparent);
    transform: skewX(-25deg);
    animation: shine 6s infinite;
    pointer-events: none;
}

@keyframes shine {
    0% {
        transform: translateX(-100%) skewX(-25deg);
    }

    20% {
        transform: translateX(100%) skewX(-25deg);
    }

    100% {
        transform: translateX(100%) skewX(-25deg);
    }
}

.glass-card-content {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 12px;
}
</style>

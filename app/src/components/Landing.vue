<template>
    <div class="landing-page">
        <!-- Section 1: Header -->
        <header class="landing-header">
            <router-link to="/" class="landing-logo-link">
                <span class="landing-logo"></span>
            </router-link>
            <div class="landing-actions">
                <router-link to="/signup" class="landing-btn landing-btn-signup">
                    [SIGN UP]
                </router-link>
                <router-link to="/login" class="landing-btn landing-btn-login">
                    [LOG IN]
                </router-link>
            </div>
        </header>

        <!-- Section 2: Hero Content -->
        <section class="hero-section">
            <div class="hero-content">
                <!-- Background patterns - hidden on mobile, shown on tablet/desktop -->
                <div class="hero-background light-only" aria-hidden="true"></div>
                <div class="hero-background dark-only" aria-hidden="true"></div>
                
                <div class="hero-text">
                    <h1 class="hero-title">
                        RESIST_EMAIL_<br>SURVEILLANCE
                    </h1>
                    <div class="hero-features">
                        <div class="feature-item">
                            <span class="feature-bullet">&gt;</span>
                            <p class="feature-text">Forward emails through isolated aliases with Mailx.</p>
                        </div>
                        <div class="feature-item">
                            <span class="feature-bullet">&gt;</span>
                            <p class="feature-text">Send and reply without exposing your primary email address.</p>
                        </div>
                        <div class="feature-item">
                            <span class="feature-bullet">&gt;</span>
                            <p class="feature-text">Supports PGP encryption, domain choice, multiple</p>
                        </div>
                        <div class="feature-item">
                            <span class="feature-bullet">&gt;</span>
                            <p class="feature-text">recipients, and catch-all aliases.</p>
                        </div>
                    </div>
                </div>
                <div class="hero-cta">
                    <router-link to="/signup" class="cta-btn cta-primary">
                        ./GET_ACCESS
                    </router-link>
                    <a href="https://github.com/ivpn/email" target="_blank" rel="noopener noreferrer" class="cta-btn cta-secondary">
                        ./VIEW_SOURCE
                    </a>
                </div>
            </div>
            <div class="hero-screenshot">
                <div class="screenshot-container">
                    <div class="screenshot-border">
                        <!-- Light mode screenshot -->
                        <img 
                            src="../assets/landing-screenshot-light.png" 
                            srcset="../assets/landing-screenshot-light.png 1x, ../assets/landing-screenshot-light@2x.png 2x"
                            alt="Mailx Application Screenshot"
                            class="screenshot-image light-only"
                        />
                        <!-- Dark mode screenshot -->
                        <img 
                            src="../assets/landing-screenshot-dark.png" 
                            srcset="../assets/landing-screenshot-dark.png 1x, ../assets/landing-screenshot-dark@2x.png 2x"
                            alt="Mailx Application Screenshot"
                            class="screenshot-image dark-only"
                        />
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped>
.landing-page {
    @apply w-full;
    @apply bg-white dark:bg-[#0a0a0a];
    @apply min-h-screen;
}

.landing-header {
    @apply flex items-center justify-between;
    @apply border-b border-solid;
    @apply border-light-gray-6 dark:border-dark-neutral-3;
    @apply py-3;
    /* Responsive padding: 110px desktop, 32px tablet, 16px mobile */
    @apply px-4 md:px-8 lg:px-[110px];
    @apply w-full;
}

.landing-logo-link {
    @apply p-0 flex items-center;
}

.landing-logo {
    @apply block bg-no-repeat bg-contain bg-center;
    @apply bg-[url("../assets/mailX-5.svg")] dark:bg-[url("../assets/mailX-3.svg")];
    /* Desktop/Tablet: 118px x 25px, Mobile: 98px x 20px */
    @apply w-[98px] h-5 md:w-[118px] md:h-[25px];
}

.landing-actions {
    @apply flex gap-2 items-center;
}

.landing-btn {
    @apply flex items-center justify-center;
    @apply font-medium text-sm leading-none;
    @apply py-[11px] whitespace-nowrap;
    @apply transition-colors;
}

.landing-btn-signup {
    @apply px-2;
    @apply text-[#1a75da] dark:text-[#0179ff];
}

.landing-btn-signup:hover {
    @apply opacity-80;
}

.landing-btn-login {
    @apply px-3;
    @apply bg-[#1a75da] dark:bg-[#0179ff];
    @apply text-[#fffeff] dark:text-[#0a0a0a];
}

.landing-btn-login:hover {
    @apply opacity-90;
}

/* Hero Section */
.hero-section {
    @apply w-full flex flex-col;
    @apply relative overflow-hidden;
}

/* Background patterns - hidden on mobile, shown on tablet/desktop */
.hero-background {
    @apply absolute pointer-events-none;
    @apply hidden md:block;
    /* Tether to top, right, bottom - resize with container while maintaining 1:1 ratio */
    top: 0;
    right: 0;
    bottom: 0;
    /* Use aspect ratio to maintain 1:1 while filling height */
    aspect-ratio: 1 / 1;
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center right;
}

/* Light mode backgrounds */
.hero-background.light-only {
    background-image: url('../assets/hero-bg-light-tablet.svg');
}

@media (min-width: 1024px) {
    .hero-background.light-only {
        background-image: url('../assets/hero-bg-light-desktop.svg');
    }
}

/* Dark mode backgrounds */
.hero-background.dark-only {
    background-image: url('../assets/hero-bg-dark-tablet.svg');
}

@media (min-width: 1024px) {
    .hero-background.dark-only {
        background-image: url('../assets/hero-bg-dark-desktop.svg');
    }
}

.hero-content {
    @apply w-full flex flex-col;
    @apply py-8 md:py-8 lg:py-8;
    @apply px-4 md:px-8 lg:px-[220px];
    @apply relative;
}

.hero-text {
    @apply flex flex-col gap-5 md:gap-6;
}

.hero-title {
    @apply font-bold uppercase;
    @apply text-[#1a75da] dark:text-[#449cf8];
    /* Responsive font sizes: 60px desktop, 52px tablet, 42px mobile */
    @apply text-[42px] leading-[46px] md:text-[52px] md:leading-[52px] lg:text-[60px] lg:leading-[60px];
    @apply m-0;
}

.hero-features {
    @apply flex flex-col gap-3;
}

.feature-item {
    @apply flex gap-2 items-center;
    @apply relative z-10;
}

.feature-bullet {
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-base;
    @apply shrink-0;
    @apply flex items-center justify-center;
    @apply w-6 h-6;
}

.feature-text {
    @apply text-[#115190] dark:text-[#82b8ed];
    @apply text-base leading-[21px] md:leading-4;
    @apply m-0;
    @apply flex-1;
}

.hero-cta {
    @apply flex gap-2 items-center mt-8;
    /* Full width buttons on mobile, auto width on desktop */
    @apply flex-col md:flex-row w-full md:w-auto;
}

.cta-btn {
    @apply flex items-center justify-center;
    @apply px-4 py-[13px];
    @apply font-medium text-sm leading-4;
    @apply transition-colors;
    @apply whitespace-nowrap;
    /* Full width on mobile */
    @apply w-full md:w-auto;
}

.cta-primary {
    @apply bg-[#1a75da] dark:bg-[#0179ff];
    @apply text-[#fffeff] dark:text-[#0a0a0a];
}

.cta-primary:hover {
    @apply opacity-90;
}

.cta-secondary {
    @apply border border-solid;
    @apply border-[#1a75da] dark:border-[#0179ff];
    @apply text-[#1a75da] dark:text-[#0179ff];
}

.cta-secondary:hover {
    @apply opacity-80;
}

/* Screenshot Section */
.hero-screenshot {
    @apply w-full;
    @apply relative z-10;
    @apply px-2 md:px-0 lg:px-[220px];
}

.screenshot-container {
    @apply w-full;
    @apply bg-[#f3f4f5] dark:bg-[#141212];
    @apply p-2;
}

.screenshot-border {
    @apply w-full;
    @apply bg-white dark:bg-[#191919];
    @apply border border-solid;
    @apply border-[#dbdfe5] dark:border-[#262626];
    @apply p-1;
    @apply overflow-hidden;
}

.screenshot-image {
    @apply w-full;
    @apply block;
    @apply h-auto;
}

/* Show/hide images based on theme */
.light-only {
    @apply dark:hidden;
}

.dark-only {
    @apply hidden dark:block;
}
</style>
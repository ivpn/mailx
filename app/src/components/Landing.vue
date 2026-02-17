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
                <!-- Background patterns - hidden on mobile, shown on tablet/desktop -->
                <div class="hero-background light-only" aria-hidden="true"></div>
                <div class="hero-background dark-only" aria-hidden="true"></div>
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

        <!-- Section 3: How It Works -->
        <section class="how-it-works-section">
            <div class="how-it-works-container">
                <!-- Title and Subtitle -->
                <div class="how-it-works-header">
                    <h2 class="how-it-works-title">HOW_IT_WORKS</h2>
                    <p class="how-it-works-subtitle">Mailx generates unique aliases for each service and sender.</p>
                </div>

                <!-- Flow Diagrams -->
                <div class="flow-diagrams">
                    <!-- Receive Flow -->
                    <div class="flow-card">
                        <div class="flow-header">
                            <span class="flow-title"># RECEIVE_FLOW</span>
                        </div>
                        <div class="flow-diagram">
                            <div class="flow-steps">
                                <!-- Step 1: Sender -->
                                <div class="flow-step">
                                    <div class="flow-icon-row">
                                        <img src="../assets/icons/flow-icon-user.svg" alt="" class="flow-icon" />
                                        <div class="flow-arrow">
                                            <svg width="100%" height="6" viewBox="0 0 100 6" preserveAspectRatio="none">
                                                <line x1="0" y1="3" x2="100%" y2="3" stroke="currentColor" stroke-width="1" stroke-dasharray="4 2" />
                                            </svg>
                                            <span class="flow-arrow-char">&gt;</span>
                                        </div>
                                    </div>
                                    <span class="flow-label">Sender</span>
                                </div>

                                <!-- Step 2: Alias -->
                                <div class="flow-step">
                                    <div class="flow-icon-row">
                                        <img src="../assets/icons/flow-icon-mail.svg" alt="" class="flow-icon" />
                                        <div class="flow-arrow">
                                            <svg width="100%" height="6" viewBox="0 0 100 6" preserveAspectRatio="none">
                                                <line x1="0" y1="3" x2="100%" y2="3" stroke="currentColor" stroke-width="1" stroke-dasharray="4 2" />
                                            </svg>
                                            <span class="flow-arrow-char">&gt;</span>
                                        </div>
                                    </div>
                                    <span class="flow-label">Alias</span>
                                </div>

                                <!-- Step 3: Mailx -->
                                <div class="flow-step">
                                    <div class="flow-icon-row">
                                        <img src="../assets/icons/flow-icon-server.svg" alt="" class="flow-icon" />
                                        <div class="flow-arrow">
                                            <svg width="100%" height="6" viewBox="0 0 100 6" preserveAspectRatio="none">
                                                <line x1="0" y1="3" x2="100%" y2="3" stroke="currentColor" stroke-width="1" stroke-dasharray="4 2" />
                                            </svg>
                                            <span class="flow-arrow-char">&gt;</span>
                                        </div>
                                    </div>
                                    <span class="flow-label">Mailx</span>
                                </div>

                                <!-- Step 4: Inbox -->
                                <div class="flow-step flow-step-last">
                                    <div class="flow-icon-row">
                                        <img src="../assets/icons/flow-icon-mail.svg" alt="" class="flow-icon" />
                                    </div>
                                    <span class="flow-label">Inbox</span>
                                </div>
                            </div>
                        </div>
                        <div class="flow-instructions">
                            <div class="instruction-item">
                                <span class="instruction-number">#1</span>
                                <span class="instruction-text">Create an email alias</span>
                            </div>
                            <div class="instruction-item">
                                <span class="instruction-number">#2</span>
                                <span class="instruction-text">Share the alias instead of your real email</span>
                            </div>
                            <div class="instruction-item">
                                <span class="instruction-number">#3</span>
                                <span class="instruction-text">Receive forwarded emails in your inbox</span>
                            </div>
                        </div>
                    </div>

                    <!-- Send/Reply Flow -->
                    <div class="flow-card">
                        <div class="flow-header">
                            <span class="flow-title"># SEND/REPLY_FLOW</span>
                        </div>
                        <div class="flow-diagram">
                            <div class="flow-steps">
                                <!-- Step 1: Outbox -->
                                <div class="flow-step">
                                    <div class="flow-icon-row">
                                        <img src="../assets/icons/flow-icon-mail.svg" alt="" class="flow-icon" />
                                        <div class="flow-arrow">
                                            <svg width="100%" height="6" viewBox="0 0 100 6" preserveAspectRatio="none">
                                                <line x1="0" y1="3" x2="100%" y2="3" stroke="currentColor" stroke-width="1" stroke-dasharray="4 2" />
                                            </svg>
                                            <span class="flow-arrow-char">&gt;</span>
                                        </div>
                                    </div>
                                    <span class="flow-label">Outbox</span>
                                </div>

                                <!-- Step 2: Mailx -->
                                <div class="flow-step">
                                    <div class="flow-icon-row">
                                        <img src="../assets/icons/flow-icon-server.svg" alt="" class="flow-icon" />
                                        <div class="flow-arrow">
                                            <svg width="100%" height="6" viewBox="0 0 100 6" preserveAspectRatio="none">
                                                <line x1="0" y1="3" x2="100%" y2="3" stroke="currentColor" stroke-width="1" stroke-dasharray="4 2" />
                                            </svg>
                                            <span class="flow-arrow-char">&gt;</span>
                                        </div>
                                    </div>
                                    <span class="flow-label">Mailx</span>
                                </div>

                                <!-- Step 3: Alias -->
                                <div class="flow-step">
                                    <div class="flow-icon-row">
                                        <img src="../assets/icons/flow-icon-mail.svg" alt="" class="flow-icon" />
                                        <div class="flow-arrow">
                                            <svg width="100%" height="6" viewBox="0 0 100 6" preserveAspectRatio="none">
                                                <line x1="0" y1="3" x2="100%" y2="3" stroke="currentColor" stroke-width="1" stroke-dasharray="4 2" />
                                            </svg>
                                            <span class="flow-arrow-char">&gt;</span>
                                        </div>
                                    </div>
                                    <span class="flow-label">Alias</span>
                                </div>

                                <!-- Step 4: Recipient -->
                                <div class="flow-step flow-step-last">
                                    <div class="flow-icon-row">
                                        <img src="../assets/icons/flow-icon-user.svg" alt="" class="flow-icon" />
                                    </div>
                                    <span class="flow-label">Recipient</span>
                                </div>
                            </div>
                        </div>
                        <div class="flow-instructions">
                            <div class="instruction-item">
                                <span class="instruction-number">#1</span>
                                <span class="instruction-text">Create a send address for your recipient</span>
                            </div>
                            <div class="instruction-item">
                                <span class="instruction-number">#2</span>
                                <span class="instruction-text">Email that address from your inbox</span>
                            </div>
                            <div class="instruction-item">
                                <span class="instruction-number">#3</span>
                                <span class="instruction-text">Recipient gets your email from the alias</span>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Feature Cards -->
                <div class="feature-cards">
                    <div class="feature-card">
                        <div class="feature-card-header">
                            <img src="../assets/icons/feature-icon-lock.svg" alt="" class="feature-icon" />
                            <h3 class="feature-card-title">BLOCK_SPAM_AND_PHISHING</h3>
                        </div>
                        <p class="feature-card-text">Disable or delete any alias instantly. Even if one gets compromised, your real email remains unexposed.</p>
                    </div>

                    <div class="feature-card">
                        <div class="feature-card-header">
                            <img src="../assets/icons/feature-icon-data-leak.svg" alt="" class="feature-icon" />
                            <h3 class="feature-card-title">DETECT_DATA_LEAKS</h3>
                        </div>
                        <p class="feature-card-text">One alias per service means you'll know exactly who shared your email without your consent.</p>
                    </div>

                    <div class="feature-card">
                        <div class="feature-card-header">
                            <img src="../assets/icons/feature-icon-users.svg" alt="" class="feature-icon" />
                            <h3 class="feature-card-title">PREVENT_IDENTITY_CORRELATION</h3>
                        </div>
                        <p class="feature-card-text">Each service gets a unique alias. Advertisers and data brokers can't connect your accounts across platforms.</p>
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

    * {
        @apply font-secondary;
    }
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
    top: -30px;
    right: 0;
    bottom: -30px;
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
    @apply flex flex-col;
    @apply py-8 md:py-8 lg:py-8;
    @apply px-4 md:px-8;
    @apply max-w-[1000px] mx-auto;
    @apply relative;
    @apply w-full;
}

.hero-text {
    @apply flex flex-col gap-5 md:gap-6 z-10;
    @apply relative;
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
    @apply flex gap-2 items-center mt-8 z-10;
    /* Full width buttons on mobile, auto width on desktop */
    @apply flex-col md:flex-row w-full md:w-auto;
    @apply relative;
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
    @apply px-[15px] py-[12px];
    @apply border border-solid;
    @apply border-[#1a75da] dark:border-[#0179ff];
    @apply text-[#1a75da] dark:text-[#0179ff];
}

.cta-secondary:hover {
    @apply opacity-80;
}

/* Screenshot Section */
.hero-screenshot {
    @apply relative z-10;
    @apply px-4 md:px-8;
    @apply max-w-[1000px] mx-auto;
    @apply w-full;
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

/* Override theme visibility for backgrounds on mobile - force hide */
.hero-background.light-only {
    @apply dark:hidden md:block;
}

.hero-background.dark-only {
    @apply hidden dark:md:block;
}

/* Section 3: How It Works */
.how-it-works-section {
    @apply w-full;
    @apply bg-white dark:bg-[#0a0a0a];
    @apply py-[72px];
}

.how-it-works-container {
    @apply flex flex-col gap-8;
    @apply w-full;
    @apply max-w-[1000px] mx-auto;
    @apply px-4 md:px-8;
}

.how-it-works-header {
    @apply flex flex-col gap-4;
}

.how-it-works-title {
    @apply font-bold uppercase;
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-[36px] leading-[40px];
    @apply m-0;
}

.how-it-works-subtitle {
    @apply text-[#115190] dark:text-[#82b8ed];
    @apply text-base leading-4;
    @apply m-0;
    @apply max-w-[587px];
}

/* Flow Diagrams */
.flow-diagrams {
    @apply flex flex-col lg:flex-row gap-0;
    @apply border border-solid;
    @apply border-[#dbdfe5] dark:border-[#282727];
}

.flow-card {
    @apply flex-1;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply lg:border-r lg:last:border-r-0;
    @apply border-b lg:border-b-0 last:border-b-0;
}

.flow-header {
    @apply bg-[#f3f4f5] dark:bg-[#0a0a0a];
    @apply border-b border-dashed;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply px-5 py-3;
}

.flow-title {
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-xs leading-3;
    @apply font-normal;
}

.flow-diagram {
    @apply bg-[#fafbfc] dark:bg-[#0f0f11];
    @apply p-5;
}

.flow-steps {
    @apply flex items-center gap-1;
}

.flow-step {
    @apply flex-1;
    @apply flex flex-col gap-[10px];
}

.flow-step-last {
    @apply flex-none;
}

.flow-icon-row {
    @apply flex items-center gap-0;
}

.flow-icon {
    @apply w-5 h-5;
    @apply shrink-0;
}

.flow-arrow {
    @apply flex-1;
    @apply text-[#1a75da] dark:text-[#282727];
    @apply flex items-center;
    @apply relative;
    @apply mx-2;
}

.flow-arrow-char {
    @apply text-[#1a75da] dark:text-[#282727];
    @apply text-xs;
    @apply absolute right-[-4px] top-[-4px];
    @apply leading-none;
}

.flow-label {
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-xs leading-4;
}

.flow-instructions {
    @apply bg-[#fafbfc] dark:bg-[#0f0f11];
    @apply border-t border-solid;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply p-5;
    @apply flex flex-col gap-4;
}

.instruction-item {
    @apply flex gap-2 items-start;
}

.instruction-number {
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-sm leading-4;
    @apply w-5 h-5;
    @apply flex items-center justify-center;
    @apply shrink-0;
}

.instruction-text {
    @apply text-[#115190] dark:text-[#82b8ed];
    @apply text-sm leading-5;
    @apply flex-1;
}

/* Feature Cards */
.feature-cards {
    @apply border border-solid;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply flex flex-col lg:flex-row;
}

.feature-card {
    @apply flex-1;
    @apply p-6;
    @apply flex flex-col gap-4;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply border-b lg:border-b-0 lg:border-r;
    @apply last:border-b-0 lg:last:border-r-0;
}

.feature-card-header {
    @apply flex gap-3 items-center;
}

.feature-icon {
    @apply w-6 h-6;
    @apply shrink-0;
}

.feature-card-title {
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-base leading-4;
    @apply font-medium uppercase;
    @apply m-0;
}

.feature-card-text {
    @apply text-[#115190] dark:text-[#82b8ed];
    @apply text-sm leading-5;
    @apply m-0;
}
</style>
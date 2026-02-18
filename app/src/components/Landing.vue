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
                                        <img src="../assets/icons/flow-icon-user-light.svg" alt="" class="flow-icon light-only" />
                                        <img src="../assets/icons/flow-icon-user-dark.svg" alt="" class="flow-icon dark-only" />
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
                                        <img src="../assets/icons/flow-icon-mail-light.svg" alt="" class="flow-icon light-only" />
                                        <img src="../assets/icons/flow-icon-mail-dark.svg" alt="" class="flow-icon dark-only" />
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
                                        <img src="../assets/icons/flow-icon-server-light.svg" alt="" class="flow-icon light-only" />
                                        <img src="../assets/icons/flow-icon-server-dark.svg" alt="" class="flow-icon dark-only" />
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
                                        <img src="../assets/icons/flow-icon-mail-light.svg" alt="" class="flow-icon light-only" />
                                        <img src="../assets/icons/flow-icon-mail-dark.svg" alt="" class="flow-icon dark-only" />
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
                                        <img src="../assets/icons/flow-icon-mail-light.svg" alt="" class="flow-icon light-only" />
                                        <img src="../assets/icons/flow-icon-mail-dark.svg" alt="" class="flow-icon dark-only" />
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
                                        <img src="../assets/icons/flow-icon-server-light.svg" alt="" class="flow-icon light-only" />
                                        <img src="../assets/icons/flow-icon-server-dark.svg" alt="" class="flow-icon dark-only" />
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
                                        <img src="../assets/icons/flow-icon-mail-light.svg" alt="" class="flow-icon light-only" />
                                        <img src="../assets/icons/flow-icon-mail-dark.svg" alt="" class="flow-icon dark-only" />
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
                                        <img src="../assets/icons/flow-icon-user-light.svg" alt="" class="flow-icon light-only" />
                                        <img src="../assets/icons/flow-icon-user-dark.svg" alt="" class="flow-icon dark-only" />
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
                            <img src="../assets/icons/feature-icon-lock-light.svg" alt="" class="feature-icon light-only" />
                            <img src="../assets/icons/feature-icon-lock-dark.svg" alt="" class="feature-icon dark-only" />
                            <h3 class="feature-card-title">BLOCK_SPAM_AND_PHISHING</h3>
                        </div>
                        <p class="feature-card-text">Disable or delete any alias instantly. Even if one gets compromised, your real email remains unexposed.</p>
                    </div>

                    <div class="feature-card">
                        <div class="feature-card-header">
                            <img src="../assets/icons/feature-icon-data-leak-light.svg" alt="" class="feature-icon light-only" />
                            <img src="../assets/icons/feature-icon-data-leak-dark.svg" alt="" class="feature-icon dark-only" />
                            <h3 class="feature-card-title">DETECT_DATA_LEAKS</h3>
                        </div>
                        <p class="feature-card-text">One alias per service means you'll know exactly who shared your email without your consent.</p>
                    </div>

                    <div class="feature-card">
                        <div class="feature-card-header">
                            <img src="../assets/icons/feature-icon-users-light.svg" alt="" class="feature-icon light-only" />
                            <img src="../assets/icons/feature-icon-users-dark.svg" alt="" class="feature-icon dark-only" />
                            <h3 class="feature-card-title">PREVENT_IDENTITY_CORRELATION</h3>
                        </div>
                        <p class="feature-card-text">Each service gets a unique alias. Advertisers and data brokers can't connect your accounts across platforms.</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Section 4: Feature Set -->
        <section class="feature-set-section">
            <div class="feature-set-container">
                <div class="feature-set-header">
                    <div class="feature-set-command">$ mailx --list-features</div>
                    <h2 class="feature-set-title">FEATURE_SET</h2>
                </div>

                <div class="feature-set-grid">
                    <div class="feature-set-card">
                        <div class="feature-set-preview">
                            <div class="feature-set-preview-bg">
                                <div class="feature-set-preview-frame">
                                    <img
                                        src="../assets/feature-set-alias-light.png"
                                        alt="New alias feature preview"
                                        class="feature-set-image light-only"
                                    />
                                    <img
                                        src="../assets/feature-set-alias-dark.png"
                                        alt="New alias feature preview"
                                        class="feature-set-image dark-only"
                                    />
                                </div>
                            </div>
                        </div>

                        <div class="feature-set-list">
                            <div class="feature-set-item">
                                <span class="feature-set-bullet">&gt;</span>
                                <p class="feature-set-text">Create aliases via Mailx browser extension</p>
                            </div>
                            <div class="feature-set-item">
                                <span class="feature-set-bullet">&gt;</span>
                                <p class="feature-set-text">Use preset domains or your own for aliases</p>
                            </div>
                            <div class="feature-set-item">
                                <span class="feature-set-bullet">&gt;</span>
                                <p class="feature-set-text">Send and reply from any alias</p>
                            </div>
                            <div class="feature-set-item">
                                <span class="feature-set-bullet">&gt;</span>
                                <p class="feature-set-text">Generate aliases on-the-fly with wildcards</p>
                            </div>
                        </div>
                    </div>

                    <div class="feature-set-card">
                        <div class="feature-set-preview">
                            <div class="feature-set-preview-bg">
                                <div class="feature-set-preview-frame">
                                    <img
                                        src="../assets/feature-set-stats-light.png"
                                        alt="Forwarding stats feature preview"
                                        class="feature-set-image light-only"
                                    />
                                    <img
                                        src="../assets/feature-set-stats-dark.png"
                                        alt="Forwarding stats feature preview"
                                        class="feature-set-image dark-only"
                                    />
                                </div>
                            </div>
                        </div>

                        <div class="feature-set-list">
                            <div class="feature-set-item">
                                <span class="feature-set-bullet">&gt;</span>
                                <p class="feature-set-text">Route forwards to one or more recipients</p>
                            </div>
                            <div class="feature-set-item">
                                <span class="feature-set-bullet">&gt;</span>
                                <p class="feature-set-text">Encrypt forwarded messages with your PGP key</p>
                            </div>
                            <div class="feature-set-item">
                                <span class="feature-set-bullet">&gt;</span>
                                <p class="feature-set-text">Enable 2FA to prevent unauthorized access</p>
                            </div>
                            <div class="feature-set-item">
                                <span class="feature-set-bullet">&gt;</span>
                                <p class="feature-set-text">Review forwarding stats and delivery logs</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Section 5: Verifiable Privacy -->
        <section class="verifiable-privacy-section">
            <div class="verifiable-privacy-container">
                <div class="verifiable-privacy-header">
                    <div class="verifiable-privacy-command">/etc/mailx/trust.cfg</div>
                    <h2 class="verifiable-privacy-title">
                        VERIFIABLE_<br class="verifiable-privacy-break">PRIVACY
                    </h2>
                </div>

                <div class="verifiable-privacy-grid">
                    <div class="verifiable-privacy-card trust-card-accountable">
                        <div class="verifiable-privacy-card-content">
                            <div class="verifiable-privacy-card-title-row">
                                <img src="../assets/icons/trust-accountable-light.svg" alt="" class="verifiable-privacy-icon light-only" />
                                <img src="../assets/icons/trust-accountable-dark.svg" alt="" class="verifiable-privacy-icon dark-only" />
                                <h3 class="verifiable-privacy-card-title">ACCOUNTABLE_OPERATORS</h3>
                            </div>
                            <p class="verifiable-privacy-card-text">Built by the public team behind IVPN, with a 15-year history in operating privacy services.</p>
                        </div>
                        <a href="https://www.ivpn.net/" target="_blank" rel="noopener noreferrer" class="verifiable-privacy-card-link">./MEET_TEAM</a>
                    </div>

                    <div class="verifiable-privacy-card trust-card-open-source">
                        <div class="verifiable-privacy-card-content">
                            <div class="verifiable-privacy-card-title-row">
                                <img src="../assets/icons/trust-open-source-light.svg" alt="" class="verifiable-privacy-icon light-only" />
                                <img src="../assets/icons/trust-open-source-dark.svg" alt="" class="verifiable-privacy-icon dark-only" />
                                <h3 class="verifiable-privacy-card-title">OPEN_SOURCE</h3>
                            </div>
                            <p class="verifiable-privacy-card-text">The entire Mailx project is open-source. Our implementation is public and available for review.</p>
                        </div>
                        <a href="https://github.com/ivpn/email" target="_blank" rel="noopener noreferrer" class="verifiable-privacy-card-link">./VIEW_SOURCE</a>
                    </div>

                    <div class="verifiable-privacy-card trust-card-security">
                        <div class="verifiable-privacy-card-content">
                            <div class="verifiable-privacy-card-title-row">
                                <img src="../assets/icons/trust-audit-light.svg" alt="" class="verifiable-privacy-icon light-only" />
                                <img src="../assets/icons/trust-audit-dark.svg" alt="" class="verifiable-privacy-icon dark-only" />
                                <h3 class="verifiable-privacy-card-title">SECURITY_AUDIT</h3>
                            </div>
                            <p class="verifiable-privacy-card-text">Mailx has undergone a third-party security audit to validate our claims and architecture.</p>
                        </div>
                        <a href="https://github.com/ivpn/email" target="_blank" rel="noopener noreferrer" class="verifiable-privacy-card-link">./READ_AUDIT</a>
                    </div>

                    <div class="verifiable-privacy-card trust-card-no-tracking">
                        <div class="verifiable-privacy-card-content">
                            <div class="verifiable-privacy-card-title-row">
                                <img src="../assets/icons/trust-no-tracking-light.svg" alt="" class="verifiable-privacy-icon light-only" />
                                <img src="../assets/icons/trust-no-tracking-dark.svg" alt="" class="verifiable-privacy-icon dark-only" />
                                <h3 class="verifiable-privacy-card-title">NO_TRACKING</h3>
                            </div>
                            <p class="verifiable-privacy-card-text">Your IP address is never logged. Forwarded emails are automatically deleted after delivery.</p>
                        </div>
                        <router-link to="/privacy" class="verifiable-privacy-card-link">./PRIVACY_POLICY</router-link>
                    </div>
                </div>
            </div>
        </section>

        <!-- Section 6: Compartmentalized Services -->
        <section class="comp-services-section">
            <div class="comp-services-container">
                <div class="comp-services-header">
                    <div class="comp-services-command">
                        <span class="dark-only">$ cat /etc/ivpn/architecture.txt</span>
                        <span class="light-only">/etc/ivpn/unlinked_auth.txt</span>
                    </div>
                    <h2 class="comp-services-title">
                        COMPARTMENTALIZED_<br class="comp-services-break">SERVICES
                    </h2>
                    <p class="comp-services-subtitle">Unlike other privacy suites, the IVPN stack works without persistent account identifiers that correlate your activity across services.</p>
                </div>

                <div class="comp-services-panel">
                    <div class="comp-services-panel-head"># Service Architecture</div>

                    <div class="comp-services-architecture">
                        <div class="comp-arch-desktop">
                            <div class="comp-arch-left">
                                <img src="../assets/icons/comp-mailx-light.svg" alt="MailX" class="light-only comp-logo-mailx" />
                                <img src="../assets/icons/comp-mailx-dark.svg" alt="MailX" class="dark-only comp-logo-mailx" />
                                <img src="../assets/icons/comp-arrow-left-light.svg" alt="" class="light-only comp-arrow-side" />
                                <img src="../assets/icons/comp-arrow-left-dark.svg" alt="" class="dark-only comp-arrow-side" />
                            </div>

                            <div class="comp-arch-center">
                                <img src="../assets/icons/comp-ivpn-desktop.svg" alt="IVPN" class="comp-logo-ivpn" />
                                <img src="../assets/icons/comp-arrow-vertical-light.svg" alt="" class="light-only comp-arrow-vertical-down" />
                                <img src="../assets/icons/comp-arrow-vertical-dark.svg" alt="" class="dark-only comp-arrow-vertical-down" />

                                <div class="comp-arch-core">
                                    <img src="../assets/icons/comp-db-light.svg" alt="" class="light-only comp-logo-db" />
                                    <img src="../assets/icons/comp-db-dark.svg" alt="" class="dark-only comp-logo-db" />
                                    <p class="comp-arch-core-text">#ZERO-LINK<br>ACCESS ARCHITECTURE</p>
                                </div>

                                <img src="../assets/icons/comp-arrow-vertical-light.svg" alt="" class="light-only comp-arrow-vertical" />
                                <img src="../assets/icons/comp-arrow-vertical-dark.svg" alt="" class="dark-only comp-arrow-vertical" />
                                <img src="../assets/icons/comp-portmark-light.svg" alt="Portmaster" class="light-only comp-logo-port" />
                                <img src="../assets/icons/comp-portmark-dark.svg" alt="Portmaster" class="dark-only comp-logo-port" />
                            </div>

                            <div class="comp-arch-right">
                                <img src="../assets/icons/comp-arrow-right-light.svg" alt="" class="light-only comp-arrow-side" />
                                <img src="../assets/icons/comp-arrow-right-dark.svg" alt="" class="dark-only comp-arrow-side" />
                                <img src="../assets/icons/comp-moddns-light.svg" alt="modDNS" class="light-only comp-logo-moddns" />
                                <img src="../assets/icons/comp-moddns-dark.svg" alt="modDNS" class="dark-only comp-logo-moddns" />
                            </div>
                        </div>

                        <div class="comp-arch-mobile">
                            <div class="comp-arch-mobile-top">
                                <img src="../assets/icons/comp-ivpn-mobile.svg" alt="IVPN" class="comp-logo-ivpn-mobile" />
                                <img src="../assets/icons/comp-arrow-vertical-light.svg" alt="" class="light-only comp-arrow-vertical-mobile-down" />
                                <img src="../assets/icons/comp-arrow-vertical-dark.svg" alt="" class="dark-only comp-arrow-vertical-mobile-down" />
                                <img src="../assets/icons/comp-db-light.svg" alt="" class="light-only comp-logo-db" />
                                <img src="../assets/icons/comp-db-dark.svg" alt="" class="dark-only comp-logo-db" />
                                <p class="comp-arch-core-text comp-arch-core-text-mobile">#ZERO-LINK ACCESS ARCHITECTURE</p>
                            </div>

                            <div class="comp-arch-mobile-bottom">
                                <img src="../assets/icons/comp-mobile-left-light.svg" alt="MailX" class="light-only comp-mobile-side" />
                                <img src="../assets/icons/comp-mobile-left-dark.svg" alt="MailX" class="dark-only comp-mobile-side" />

                                <div class="comp-arch-mobile-center">
                                    <img src="../assets/icons/comp-mobile-vertical-light.svg" alt="" class="light-only comp-arrow-mobile-center-down" />
                                    <img src="../assets/icons/comp-mobile-vertical-dark.svg" alt="" class="dark-only comp-arrow-mobile-center-down" />
                                    <img src="../assets/icons/comp-portmark-light.svg" alt="Portmaster" class="light-only comp-logo-port-mobile" />
                                    <img src="../assets/icons/comp-portmark-dark.svg" alt="Portmaster" class="dark-only comp-logo-port-mobile" />
                                </div>

                                <img src="../assets/icons/comp-mobile-right-light.svg" alt="modDNS" class="light-only comp-mobile-side" />
                                <img src="../assets/icons/comp-mobile-right-dark.svg" alt="modDNS" class="dark-only comp-mobile-side" />
                            </div>
                        </div>
                    </div>

                    <div class="comp-services-cards">
                        <div class="comp-services-card">
                            <img src="../assets/icons/comp-signin-light.svg" alt="" class="comp-services-card-icon light-only" />
                            <img src="../assets/icons/comp-signin-dark.svg" alt="" class="comp-services-card-icon dark-only" />
                            <p class="comp-services-card-text">Sign up to Mailx with true account separation</p>
                        </div>

                        <div class="comp-services-card">
                            <img src="../assets/icons/comp-cross-light.svg" alt="" class="comp-services-card-icon light-only" />
                            <img src="../assets/icons/comp-cross-dark.svg" alt="" class="comp-services-card-icon dark-only" />
                            <p class="comp-services-card-text">No cross-service correlation, even by administrators</p>
                        </div>

                        <div class="comp-services-card">
                            <img src="../assets/icons/comp-repeat-light.svg" alt="" class="comp-services-card-icon light-only" />
                            <img src="../assets/icons/comp-repeat-dark.svg" alt="" class="comp-services-card-icon dark-only" />
                            <p class="comp-services-card-text">Automatic and anonymous subscription management</p>
                        </div>
                    </div>
                </div>

                <div class="comp-services-footnote">
                    <img src="../assets/icons/comp-bulb-light.svg" alt="" class="comp-services-footnote-icon light-only" />
                    <img src="../assets/icons/comp-bulb-dark.svg" alt="" class="comp-services-footnote-icon dark-only" />
                    <p class="comp-services-footnote-text">
                        Read more in our
                        <a href="https://www.ivpn.net/blog/" target="_blank" rel="noopener noreferrer" class="comp-services-link">./BLOG</a>
                        and review the
                        <a href="https://github.com/ivpn/email" target="_blank" rel="noopener noreferrer" class="comp-services-link">./SOURCE</a>
                        implementation.
                    </p>
                </div>
            </div>
        </section>

        <!-- Section 7: Constraints -->
        <section class="constraints-section">
            <div class="constraints-container">
                <div class="constraints-header">
                    <div class="constraints-command">/etc/mailx/limitations.txt</div>
                    <h2 class="constraints-title">CONSTRAINTS</h2>
                </div>

                <div class="constraints-list">
                    <div class="constraints-item">
                        <span class="constraints-bullet">&gt;</span>
                        <p class="constraints-text">Mailx is not an email provider, it can't replace your primary email</p>
                    </div>

                    <div class="constraints-item">
                        <span class="constraints-bullet">&gt;</span>
                        <p class="constraints-text">Messages are visible to Mailx servers during relay (use PGP)</p>
                    </div>

                    <div class="constraints-item">
                        <span class="constraints-bullet">&gt;</span>
                        <p class="constraints-text">No IMAP/POP3 access. Email forwarding only.</p>
                    </div>

                    <div class="constraints-item">
                        <span class="constraints-bullet">&gt;</span>
                        <p class="constraints-text">Not designed for protection against targeted surveillance</p>
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
    @apply max-w-[1060px] mx-auto;
    @apply px-4 md:px-8;
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
    @apply max-w-[1060px] mx-auto;
    @apply px-4 md:px-8;
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
img.light-only {
    @apply dark:hidden !important;
}

img.dark-only {
    @apply hidden dark:block !important;
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
    @apply max-w-[1060px] mx-auto;
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

/* Section 4: Feature Set */
.feature-set-section {
    @apply w-full;
    @apply bg-[#f3f4f5] dark:bg-[rgba(255,255,255,0.02)];
    @apply border-x border-solid;
    @apply py-6 md:py-[52px] lg:py-[72px];
    @apply border border-solid;
    @apply border-[#dbdfe5] dark:border-[#282727];
}

.feature-set-container {
    @apply w-full;
    @apply px-4 md:px-8 lg:px-[220px];
    @apply flex flex-col gap-8;
    @apply max-w-[1060px] mx-auto;
    @apply px-4 md:px-8;
}

.feature-set-header {
    @apply flex flex-col gap-4;
}

.feature-set-command {
    @apply bg-[rgba(26,117,218,0.12)] dark:bg-[#12161b];
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-xs leading-3;
    @apply px-3 py-3;
    @apply w-fit;
}

.feature-set-title {
    @apply m-0;
    @apply uppercase font-bold;
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-[32px] leading-[32px] md:text-[36px] md:leading-[40px];
}

.feature-set-grid {
    @apply w-full;
    @apply flex flex-col md:flex-row;
    @apply pb-px md:pb-0;
}

.feature-set-card {
    @apply flex-1;
    @apply border border-solid;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply -mb-px md:mb-0 md:-mr-px;
    @apply overflow-hidden;
    @apply flex flex-col;
}

.feature-set-preview {
    @apply bg-white dark:bg-[#0a0a0a];
    @apply border-b;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply px-4 pt-4 md:px-6 md:pt-6;
    @apply flex items-end justify-center;
    @apply relative;
}

.feature-set-preview::after {
    content: none;
}

.feature-set-preview-bg {
    @apply bg-[#f3f4f5] dark:bg-[#141212];
    @apply p-[6px];
    @apply w-full max-w-[450px];
}

.feature-set-preview-frame {
    @apply bg-white dark:bg-[#191919];
    @apply border border-solid;
    @apply border-[#dbdfe5] dark:border-[#262626];
    @apply p-[2px];
    @apply overflow-hidden;
}

.feature-set-image {
    @apply w-full h-auto;
    @apply block;
}

.feature-set-list {
    /* @apply border border-solid;
    @apply border-[#dbdfe5] dark:border-[#282727]; */
    @apply p-4 md:p-6;
    @apply flex flex-col;
    @apply gap-3 md:gap-4;
}

.feature-set-item {
    @apply flex gap-2 items-start;
}

.feature-set-bullet {
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-base;
    @apply w-5 h-5;
    @apply shrink-0;
    @apply flex items-center justify-center;
}

.feature-set-text {
    @apply m-0;
    @apply text-sm leading-5;
    @apply text-[#115190] dark:text-[#82b8ed];
}

/* Section 5: Verifiable Privacy */
.verifiable-privacy-section {
    @apply w-full;
    @apply bg-white dark:bg-[#050607];
    @apply py-6 md:py-[52px] lg:py-[72px];
}

.verifiable-privacy-container {
    @apply w-full;
    @apply max-w-[1060px] mx-auto;
    @apply px-4 md:px-8;
    @apply flex flex-col gap-8;
}

.verifiable-privacy-header {
    @apply flex flex-col gap-4;
}

.verifiable-privacy-command {
    @apply bg-[rgba(26,117,218,0.12)] dark:bg-[#12161b];
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-xs leading-3;
    @apply px-3 py-3;
    @apply w-fit;
}

.verifiable-privacy-title {
    @apply m-0;
    @apply uppercase font-bold;
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-[32px] leading-[32px] md:text-[36px] md:leading-[40px];
}

.verifiable-privacy-break {
    @apply md:hidden;
}

.verifiable-privacy-grid {
    @apply w-full;
    @apply flex flex-col md:flex-row md:flex-wrap;
    @apply pb-px md:pb-0 md:pr-px;
}

.verifiable-privacy-card {
    @apply w-full md:w-1/2;
    @apply border border-solid;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply -mb-px md:mb-0 md:-mr-px;
    @apply px-4 py-5 md:p-6;
    @apply flex flex-col gap-4;
    @apply md:mb-[-1px] lg:mb-[-1px];
}

.verifiable-privacy-card-content {
    @apply flex flex-col gap-3;
}

.verifiable-privacy-card-title-row {
    @apply flex items-center gap-3;
}

.verifiable-privacy-icon {
    @apply w-6 h-6;
    @apply shrink-0;
}

.verifiable-privacy-card-title {
    @apply m-0;
    @apply text-base leading-4 font-medium uppercase;
    @apply text-[#1a75da] dark:text-[#449cf8];
}

.verifiable-privacy-card-text {
    @apply m-0;
    @apply text-sm leading-5;
    @apply text-[#115190] dark:text-[#82b8ed];
}

.verifiable-privacy-card-link {
    @apply m-0;
    @apply inline-block;
    @apply pt-2;
    @apply text-sm leading-4 font-medium;
    @apply text-[#1a75da] dark:text-[#0179ff];
    @apply no-underline;
}

.verifiable-privacy-card-link:hover {
    @apply opacity-80;
}

.trust-card-accountable {
    @apply order-1;
}

.trust-card-open-source {
    @apply order-2 md:order-3 lg:order-2;
}

.trust-card-security {
    @apply order-3 md:order-2 lg:order-3;
}

.trust-card-no-tracking {
    @apply order-4;
}

/* Section 6: Compartmentalized Services */
.comp-services-section {
    @apply w-full;
    @apply bg-[#f3f4f5] dark:bg-[rgba(255,255,255,0.02)];
    @apply py-6 md:py-[52px] lg:py-[72px];
    @apply md:border md:border-solid;
    @apply md:border-[#dbdfe5] md:dark:border-[#282727];
}

.comp-services-container {
    @apply w-full;
    @apply px-4 md:px-8 lg:px-[220px];
    @apply flex flex-col gap-8 md:gap-[42px] lg:gap-8;
}

.comp-services-header {
    @apply flex flex-col gap-4;
}

.comp-services-command {
    @apply bg-[rgba(26,117,218,0.12)] dark:bg-[#12161b];
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-xs leading-3;
    @apply px-3 py-3;
    @apply w-fit;
}

.comp-services-title {
    @apply m-0;
    @apply uppercase font-bold;
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-[32px] leading-[32px] md:text-[36px] md:leading-[40px];
}

.comp-services-break {
    @apply md:hidden;
}

.comp-services-subtitle {
    @apply m-0;
    @apply text-base leading-[21px];
    @apply text-[#115190] dark:text-[#82b8ed];
    @apply max-w-[784px];
}

.comp-services-panel {
    @apply border border-solid;
    @apply border-[#dbdfe5] dark:border-[#282727];
}

.comp-services-panel-head {
    @apply px-5 md:px-6 py-3;
    @apply text-xs leading-3;
    @apply text-[#115190] dark:text-[#449cf8];
    @apply border-b border-dashed;
    @apply border-[#dbdfe5] dark:border-[#282727];
}

.comp-services-architecture {
    @apply border-b border-solid;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply bg-[#fffeff] dark:bg-[#0a0a0a];
    @apply px-5 py-5 md:px-6 md:py-[42px];
}

.comp-arch-desktop {
    @apply hidden md:flex;
    @apply items-center justify-center;
    @apply gap-3;
}

.comp-arch-left,
.comp-arch-right {
    @apply flex items-center gap-3;
    @apply pt-[18px];
}

.comp-arch-center {
    @apply flex flex-col items-center;
    @apply gap-2;
}

.comp-logo-mailx {
    @apply w-[88px] h-[19px];
}

.comp-logo-moddns {
    @apply w-[92px] h-[18px];
}

.comp-logo-ivpn {
    @apply w-[61px] h-[21px];
}

.comp-arrow-side {
    @apply w-8 h-[6px];
}

.comp-arrow-vertical {
    @apply w-8 h-[6px] rotate-90;
}

.comp-arrow-vertical-up {
    @apply w-8 h-[6px] rotate-90;
}

.comp-arrow-vertical-down {
    @apply w-8 h-[6px] rotate-[270deg];
}

.comp-arch-core {
    @apply flex flex-col items-center;
    @apply gap-3 p-3;
}

.comp-logo-db {
    @apply w-7 h-7;
}

.comp-arch-core-text {
    @apply m-0;
    @apply uppercase font-bold text-center;
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply text-sm leading-4;
}

.comp-logo-port {
    @apply w-[128px] h-[26px];
}

.comp-arch-mobile {
    @apply flex md:hidden;
    @apply flex-col items-center;
    @apply w-full;
    @apply gap-4;
}

.comp-arch-mobile-top {
    @apply flex flex-col items-center;
    @apply gap-2;
}

.comp-logo-ivpn-mobile {
    @apply w-[51px] h-[18px];
}

.comp-arrow-vertical-mobile {
    @apply w-8 h-[6px] rotate-90;
}

.comp-arrow-vertical-mobile-down {
    @apply w-8 h-[6px] rotate-[270deg];
}

.comp-arch-core-text-mobile {
    @apply text-center;
}

.comp-arch-mobile-bottom {
    @apply flex items-start justify-between;
    @apply w-full;
}

.comp-mobile-side {
    @apply w-[96px] h-[60px];
}

.comp-arch-mobile-center {
    @apply flex flex-col items-center;
    @apply gap-20;
}

.comp-arrow-mobile-center {
    @apply w-[72px] h-[6px] rotate-90;
}

.comp-arrow-mobile-center-down {
    @apply w-[72px] h-[6px] rotate-[270deg];
}

.comp-logo-port-mobile {
    @apply w-[128px] h-[26px];
}

.comp-services-cards {
    @apply flex flex-col md:flex-row;
}

.comp-services-card {
    @apply flex-1;
    @apply flex flex-col gap-4;
    @apply px-4 py-5 md:p-6;
    @apply border-[#dbdfe5] dark:border-[#282727];
    @apply border-b md:border-b-0 md:border-r md:last:border-r-0;
}

.comp-services-card-icon {
    @apply w-6 h-6;
}

.comp-services-card-text {
    @apply m-0;
    @apply text-sm leading-5;
    @apply text-[#115190] dark:text-[#82b8ed];
}

.comp-services-footnote {
    @apply flex items-start md:items-center gap-2;
}

.comp-services-footnote-icon {
    @apply w-4 h-4;
    @apply mt-[2px] md:mt-0;
}

.comp-services-footnote-text {
    @apply m-0;
    @apply text-sm leading-5;
    @apply text-[#115190] dark:text-[#82b8ed];
}

.comp-services-link {
    @apply text-sm leading-5 font-medium;
    @apply text-[#1a75da] dark:text-[#449cf8];
    @apply underline;
    @apply mx-1;
}

/* Section 7: Constraints */
.constraints-section {
    @apply w-full;
    @apply bg-[#fffeff] dark:bg-[#0a0a0a];
    @apply py-6 md:py-[52px] lg:py-[72px];
}

.constraints-container {
    @apply w-full;
    @apply px-4 md:px-8 lg:px-[220px];
    @apply flex flex-col gap-8;
}

.constraints-header {
    @apply flex flex-col gap-4;
}

.constraints-command {
    @apply bg-[rgba(195,31,31,0.05)] dark:bg-[rgba(214,64,64,0.05)];
    @apply lg:bg-[rgba(195,31,31,0.1)] lg:dark:bg-[rgba(214,64,64,0.05)];
    @apply text-[#c31f1f] dark:text-[#d64040];
    @apply text-xs leading-3;
    @apply px-3 py-3;
    @apply w-fit;
}

.constraints-title {
    @apply m-0;
    @apply uppercase font-bold;
    @apply text-[#c31f1f] dark:text-[#d64040];
    @apply text-[32px] leading-[32px] md:text-[36px] md:leading-[40px];
}

.constraints-list {
    @apply flex flex-col gap-3;
}

.constraints-item {
    @apply flex items-center gap-2;
}

.constraints-bullet {
    @apply text-[#c31f1f] dark:text-[#d64040];
    @apply text-base;
    @apply w-6 h-6;
    @apply shrink-0;
    @apply flex items-center justify-center;
}

.constraints-text {
    @apply m-0;
    @apply flex-1;
    @apply text-[#8d1719] dark:text-[#eb7676];
    @apply text-sm leading-5 md:text-base md:leading-4;
}
</style>
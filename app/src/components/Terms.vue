<template>
    <div class="container mx-auto max-w-screen-lg sm:p-10 p-5">
        <p class="py-5">
            <router-link to="/login" class="flex items-center gap-2">
                <i class="icon arrow-left-line icon-accent"></i>
                Back
            </router-link>
        </p>
        <h1>Terms of Service</h1>
        <p class="text-sm">Last updated: 07 July 2025</p>

        <h2>1. Introduction</h2>
        <p>These Terms of Service ("Terms") govern your use of Mailx, an email forwarding service operated by Privatus GmbH ("we," "us," or "our"). By subscribing to or using our service, you ("you" or "user") agree to abide by these Terms.</p>

        <h2>2. Acceptance of Terms</h2>
        <p>By accessing or using our service, you agree to be bound by these Terms and our <a href="/privacy">Privacy Policy</a>. If you do not agree with any part of these Terms, do not use our service.</p>

        <h2>3. User Responsibilities</h2>
        <p>You agree to use our service only for lawful purposes and in a manner that does not infringe or interfere with any third party’s use of the service. Prohibited activities include, but are not limited to:</p>
        <ul>
            <li>Engaging in illegal activities, including criminal or terrorist actions, fraud, identity theft, and money laundering.</li>
            <li>Attempting to hack, attack, or gain unauthorized access to any computers, networks, accounts, or systems.</li>
            <li>Transmitting or distributing viruses, worms, trojans, or other malicious software.</li>
            <li>Sending harassing, defamatory, abusive, threatening, or hateful messages.</li>
            <li>Sending unsolicited emails, bulk messages, or any form of spam, including commercial advertising.</li>
            <li>Engaging in phishing or activities that attempt to steal personal information.</li>
            <li>Distributing or transmitting copyrighted materials without proper authorization.</li>
            <li>Participating in activities related to child pornography or exploitation.</li>
            <li>Using stolen or unauthorized payment information, including credit card numbers and online payment accounts.</li>
        </ul>

        <h2>4. Account Usage and Security</h2>
        <p>Your account is for individual use only; sharing your account or login credentials is prohibited. You are responsible for all activities that occur under your account. Notify us immediately of any unauthorized use or security breaches.</p>

        <h2>5. Intellectual Property & Licence</h2>
        <p>All interfaces, graphics, logos, trademarks, and all other intellectual-property rights in or relating to Mailx (“Service IP”) are and shall remain the exclusive property of Privatus GmbH and its licensors. No rights are granted to you except as expressly set out in these Terms.<br>
        We grant you a personal, non-exclusive, non-transferable, revocable licence to access and use the Service for your own lawful, non-commercial purposes, strictly in accordance with these Terms.<br>
        You must <strong>not</strong>:</p>
        <ul>
            <li>Copy, modify, translate, reverse-engineer, decompile, disassemble or create derivative works from the Service IP, unless otherwise permitted by the source code license;</li>
            <li>Rent, lease, sell, sublicense, assign, distribute, publish, or otherwise commercialise the service or any access to it;</li>
            <li>Remove or alter any copyright, trademark or other proprietary notice contained in the Service;</li>
            <li>Use any automated means (bots, scrapers) to access or copy the service or its underlying databases.</li>
        </ul>
        <p>All rights not expressly granted to you in Section 2 are reserved by PrivatusGmbH. Any unauthorised use of the Service IP terminates the licence granted above and may violate applicable intellectual-property laws.</p>

        <h2>6. Service Availability and Modifications</h2>
        <p>We strive to provide reliable service but do not guarantee uninterrupted availability. We reserve the right to modify, suspend, or discontinue the service (or any part thereof) at any time, with or without notice.<br>
        Technical support is provided via email and is limited to issues related to the functionality of the service.</p>

        <h2>7. Termination and Suspension</h2>
        <p>We may suspend or terminate your account immediately, without prior notice or liability, upon breach of these Terms. No refunds will be provided if your account is terminated due to a breach of these Terms.</p>
        <p>Upon termination, your right to use the Mailx service ceases immediately. Provisions that should survive termination will remain in effect, including ownership provisions, warranty disclaimers, indemnity, and limitations of liability.</p>

        <h2>8. Changes to Terms</h2>
        <p>We reserve the right to modify or replace these Terms at any time. Significant changes will be communicated by posting the new Terms on our website and updating the "Effective Date." By continuing to use our service after revisions become effective, you agree to be bound by the revised Terms.</p>

        <h2>9. Limitation of Liability and Disclaimer of Warranties</h2>
        <p>The service is provided on an "as is" and "as available" basis. To the fullest extent permitted by Swiss law, we make no express or implied warranties or conditions of any kind, including warranties of merchantability, fitness for a particular purpose, title or non-infringement, nor do we warrant that the service will meet your requirements or be uninterrupted, timely, secure or error-free.</p>
        <p>In no event shall we be liable for any indirect, incidental, special, consequential, or punitive damages, including loss of profits, data, use, goodwill, or other intangible losses.</p>

        <h2>10. Indemnification</h2>
        <p>You agree to defend, indemnify, and hold harmless Privatus GmbH and its officers, directors, employees, and agents from any claims, damages, losses, liabilities, costs, or expenses arising from:</p>
        <ul>
            <li>Your use of and access to the service.</li>
            <li>Your violation of any term of these Terms.</li>
            <li>Your violation of any third-party right, including intellectual property or privacy rights.</li>
        </ul>

        <h2>11. Governing Law</h2>
        <p>These Terms are governed by and construed in accordance with the laws of Switzerland, without regard to its conflicts of law provisions.</p>

        <hr>

        <h2>Contact Us</h2>
        <p>If you have any questions about these Terms, please contact us at <a href="mailto:mailx@ivpn.net">mailx@ivpn.net</a>.</p>
        <p>This document constitutes the entire Terms and Conditions for your use of the Mailx email forwarding service operated by Privatus GmbH.</p>
        <p>Privatus GmbH<br>
        Baarerstrasse 82<br>
        6300 Zug<br>
        Switzerland</p>

        <br>

        <Footer />
    </div>
</template>

<script setup lang="ts">
import Footer from './Footer.vue'
</script>

<style scoped>
h2 {
    @apply mt-12;
}

p {
    @apply leading-7;
}

hr {
    @apply my-10;
}
</style>
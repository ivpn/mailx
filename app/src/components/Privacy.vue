<template>
    <div class="container mx-auto max-w-screen-lg sm:p-10 p-5">
        <p class="py-5">
            <router-link to="/login" class="flex items-center gap-2">
                <i class="icon arrow-left-line icon-accent"></i>
                Back
            </router-link>
        </p>
        <h1>Privacy Policy</h1>
        <p class="text-sm">Last updated: 15 June 2025</p>

        <p>MailX is developed and operated by the team behind IVPN. All services we offer to customers are built for privacy. Internally, we know what exactly that means. If a choice needs to be made between one practice that deepens a user’s privacy, and another that would diminish it but accelerate our growth, we’ll always take the slower, more private option.</p>
        <p>We are committed to being transparent about how we collect, use, and protect your data after you sign up to the MailX service. Below we offer a concise and human-readable summary of our policies.</p>

        <h2>What data don’t you log?</h2>

        <p>Our goal is minimal data collection. We don't log any customer related information unless it is absolutely required for the operation of the email forwarding service. As such, we don't log:</p>
        <ul>
            <li>information about non-authenticated website visits</li>
            <li>we don’t log customer IP addresses</li>
            <li>we don’t log bandwidth usage</li>
        </ul>

        <h2>Do you use cookies or third-party services?</h2>

        <p>We use cookies solely for the purpose of managing authenticated sessions.</p>
        <p>No third-party cookies, trackers, or analytics are deployed when visiting our website or using the service.</p>

        <h2>What information do you store about my account?</h2>

        <p>For users signing up with passkeys, we store the registered email address in our database.</p>
        <p>For users signing up with credentials, we store registered email address and a password hashed with Argon2 hashing algorithm.</p>
        <p>Further, we store additional email addresses optionally added as ‘Recipient’ for forwarded emails.</p>
        <p>We keep a record of the number of relayed or blocked emails for 90 days from the date and time the corresponding event occurred, used to provide relevant statistics to our customers.</p>

        <h2>How do you process and what information you store about forwarded emails?</h2>

        <p>Emails are processed through a self-hosted Postfix server.</p>
        <p>Delivered emails are immediately removed from the Postfix queues. Undelivered emails are stored in a deferred queue. After 5 days, Postfix stops attempting delivery and discards the emails from the deferred queue.</p>
        <p>The Postfix server is configured with the <a href="https://docker-mailserver.github.io/docker-mailserver/latest/config/environment/#log_level" target="_blank">"info" log level</a>, which means it records general information about email activity (such as when messages are sent or received). These logs are kept for 7 days to support system maintenance and troubleshooting.</p>

        <h2>I’ve signed up to the service through my IVPN subscription. Can you associate my MailX account with my IVPN Account ID?</h2>

        <p>When your sign up through the IVPN service, a temporary MailX signup link is generated in the IVPN database. Once the MailX signup is completed, the link and corresponding identifiers are removed from the IVPN database to prevent any association between MailX and IVPN accounts.</p>

        <h2>What information is retained when I stop using your service?</h2>

        <p>When a MailX account is terminated due to subscription ending, non-payment or for any other reason, all settings and data associated with that MailX account including the account itself is automatically deleted after 180 days. If you want to cancel your account and delete your data immediately, simply click on the ‘Delete account’ button in the Account settings.</p>

        <h2>Are there any ways to verify your claims regarding privacy?</h2>

        <p>The entire project is open source with repositories available on GitHub for customers to verify our privacy practices.</p>

        <h2>User Rights</h2>

        <p>Reasonable requests for release of a specific user’s data will be honoured within 28 days of an acceptable request from a user or person with a provable legal relationship with that user.</p>
        <p>We reserve the right to refuse or charge for requests that are manifestly unfounded or excessive. Any refused subject access requests will be responded to without undue delay including the refusal reason as well as recourse to refer to the supervisory authority.</p>

        <h2>Changes to policy</h2>

        <p>Privatus GmbH, operator of MailX reserves the right to change this privacy policy at any time. In such cases, we will take every reasonable step to ensure that these changes are brought to your attention by posting all changes prominently on the MailX website for a reasonable period of time, before the new policy becomes effective as well as emailing our existing customers.</p>
        <p>If you have any questions or concerns about our privacy practices, contact us at mailx@ivpn.net.</p>

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
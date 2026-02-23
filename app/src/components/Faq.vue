<template>
    <div class="container mx-auto max-w-screen-lg sm:p-10 p-5">
        <p class="py-5">
            <router-link to="/" class="flex items-center gap-2">
                <i class="icon arrow-left-line icon-accent"></i>
                Back
            </router-link>
        </p>
        <h1>Mailx FAQ</h1>
        <p class="text-sm">Last updated: 20 June 2025</p>

        <h2>Aliases: How do I add an email alias?</h2>
        <p>Click the 'Aliases' tab on the left side of the page.</p>
        <p>Click the 'New Alias' button on the top-right of the page, then choose an option.</p>
        <p>Add an optional 'Description' and one or more 'Recipients' from the list.</p>
        <p>Expand 'Custom settings' to select the 'Format' (Words, Random, or UUID), 'Domain', and optional 'From name'.</p>
        <p>Click the 'Create and copy to clipboard' button to add the alias to your account.</p>

        <h2>Aliases: What is a Wildcard Alias?</h2>
        <p>A 'Wildcard Alias' offers a way to use one address to receive mail for virtually infinite aliases.  If access to the Internet is not available and an alias is required, remember the Wildcard suffix and add an on-the-fly prefix to create a valid email alias.</p>
        <p>Note: The 'Wildcard Alias' limit per domain is two.</p>
        <p>To create a 'Wildcard Alias', navigate to the 'Wildcard' page, and click the 'New Wildcard' button.</p>
        <p>Choose an easy to remember suffix, an add an optional 'Description' and optional 'From name'.  Select one or more 'Recipient(s)', and choose the 'Domain'.</p>
        <p>Click the 'Create and copy to clipboard' button to add the alias to your account.</p>
        <p>To make use of the 'Wildcard Alias', add a prefix to your alias.  For example, use 'dentist' plus your 'Wildcard Alias' address for an instant alias: dentist+remember_this@domain</p>

        <h2>Aliases: How do I send a message using as email alias?</h2>
        <p>On the 'Aliases' page, choose an alias, and click the ellipsis menu button under the 'Actions' column on the right side.  Click the 'Send' option.</p>
        <p>In the 'To email:' field, enter the email address for your intended recipient.</p>
        <p>Click the 'Show address' button, and click the address to copy it.</p>
        <p>In your preferred mail client or webmail system associated with a verified 'Recipient' email address, send your message to this email address.</p>
        <p>Your recipient will see the alias, and your real email address remains private.</p>

        <h2>Aliases: How do I change the email alias properties after it has been created?</h2>
        <p>On the 'Aliases' page, choose an alias, and click the ellipsis menu button under the 'Actions' column on the right side.  Click the 'Edit' option.</p>
        <p>Update the optional 'Description', optional 'From name', and/or the 'Recipient(s)' list, then click the 'Save' button.</p>

        <h2>Aliases: How do I delete an alias?</h2>
        <p>On the 'Aliases' page, choose an alias, and click the ellipsis menu button under the 'Actions' column on the right side.  Click the 'Delete' option.</p>
        <p>A confirmation box will appear.  Click 'OK' to delete the alias, or click 'Cancel'.</p>

        <h2>Aliases: What happens to deleted aliases?</h2>
        <p>A deleted email alias cannot be recreated, resurrected, restored, or otherwise reused.  This is a privacy feature to prevent any access, intentional or otherwise, to incoming messages associated with any deleted alias.</p>

        <h2>Aliases: How do I stop an alias from receiving messages?</h2>
        <p>On the 'Aliases' page, choose an alias, and toggle the 'Status' switch to the left position.</p>
        <p>The switch's background becomes gray.  Incoming messages will be blocked.</p>

        <hr>

        <h2>Recipients: How do I add a new recipient email address?</h2>
        <p>Click the 'Recipients' tab on the left side of the page.</p>
        <p>Click the 'New Recipient' button on the top-right of the page, then choose an option.</p>
        <p>Mailx supports 10 recipient addresses, and each address needs to be verified.  Unverified recipients will not receive any forwarded emails, and will be deleted after seven days.</p>

        <h2>Recipients: How does the PGP option work?</h2>
        <p>Adding a PGP public key to a Recipient mail address enables PGP encryption between the Mailx forwarding server and your real email inbox.  This increases the level of privacy for forwarded messages against the access level of your mail provider.  A mail provider has full access to all email messages, and the value of PGP-encrypted content is significantly less than plain text content.</p>
        <p>PGP encryption is not available for outgoing/Send messages, or replies.</p>
        <p>Note: This option does not provide end-to-end encryption (E2EE) for email messages from their source (i.e. original sender) to the destination (i.e. your real inbox).</p>

        <h2>Recipients: How do I enable the PGP option for forwarded messages?</h2>
        <p>Add a PGP key to a recipient address via the 'Recipients' tab.</p>
        <p>Click the ellipsis menu button under the 'Actions' column on the right side.  Click the 'Add PGP Key' option.  Paste the public key, and click the 'Add PGP Public Key' button to confirm.</p>
        <p>Toggle the 'Encryption' switch ON or OFF to govern whether forwarded messages are encrypted or not between the Mailx server and your inbox.</p>

        <h2>Recipients: How do I remove a PGP key?</h2>
        <p>Click the 'Recipients' tab on the left side of the page.</p>
        <p>Click the ellipsis menu button under the 'Actions' column on the right side.  Click the 'Delete PGP Key' option.</p>
        <p>A confirmation box will appear.  Click 'OK' to delete the key, or click 'Cancel'.</p>

        <h2>Recipients: How do I delete a recipient email address?</h2>
        <p>Click the 'Recipients' tab on the left side of the page.</p>
        <p>Click the ellipsis menu button under the 'Actions' column on the right side.  Click the 'Delete' option.</p>
        <p>A confirmation box will appear.  Click 'OK' to delete the recipient, or click 'Cancel'.</p>

        <hr>

        <h2>Account: How do I delete my account?</h2>
        <p>Click the 'Account' tab on the left side of the page.</p>
        <p>Click the 'Delete Account' button at the bottom of the page.</p>
        <p>A confirmation windows will appear.  Entering an 8-symbol code is required for confirmation.  After entering the code, click the 'Delete Account' button.  A 'Cancel' button is also available on this confirmation window.</p>

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
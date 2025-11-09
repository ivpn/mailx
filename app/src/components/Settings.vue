<template>
    <div class="card-container">
        <header class="head">
            <h2>Settings</h2>
        </header>
        <div class="card-primary">
            <h4>Default Domain</h4>
            <p>
                Set the default alias domain for new aliases created.
            </p>
            <div class="max-w-xs mb-6">
                <label for="domain">
                    Select default domain:
                </label>
                <select id="domain">
                    <option
                        v-for="(domain, index) in domains"
                        v-bind:domain
                        :selected="domain == req.domain || index === 0"
                        :key="domain">
                        {{ domain }}
                    </option>
                </select>
            </div>
            <h4>Default Recipient</h4>
            <p>
                Set the default recipient for new aliases created.
            </p>
            <div class="max-w-xs mb-6">
                <label for="recipient">
                    Select default recipient:
                </label>
                <select id="recipient" :disabled="!recipients.length">
                    <option
                        v-for="recipient in recipients"
                        v-bind:value=recipient
                        :selected="recipient == req.recipient"
                        :key="recipient">
                        {{ recipient }}
                    </option>
                </select>
            </div>
            <h4>Default Alias Format</h4>
            <p>
                Set the default alias naming format for new aliases created. Options: 1. Words ('quiet.haze16') 2. Random ('uf1h0xi') 3. UUID ('550e8400-e29b-41d4-a716-446655440000')
            </p>
            <div class="max-w-xs mb-6">
                <label for="format">
                    Select default alias format:
                </label>
                <select id="format" :disabled="!aliasFormats.length">
                    <option
                        v-for="format in aliasFormats"
                        v-bind:value=format.toLowerCase()
                        :selected="format.toLowerCase() == req.alias_format"
                        :key="format">
                        {{ format }}
                    </option>
                </select>
            </div>
            <h4>From Name</h4>
            <p>
                Set the 'From name' used for replies and emails sent using an alias. Leave it blank to use the alias email address.
            </p>
            <div class="max-w-xs mb-5">
                <label for="from-name">
                    From name:
                </label>
                <input
                    v-model="req.from_name"
                    id="from-name"
                    type="text"
                >
            </div>
            <h4>Logs</h4>
            <h5>Failed Deliveries</h5>
            <p>
                Turn logging of failed email deliveries (bounces) from aliases on or off. When enabled, any failed delivery attempts are recorded and stored for 7 days.
            </p>
            <div class="mb-8">
                <label for="log-bounce">
                    Log Failed Deliveries:
                </label>
                <input
                    @change="saveSettings"
                    v-bind:checked="req.log_bounce"
                    v-model="req.log_bounce"
                    id="log-bounce"
                    type="checkbox"
                >
            </div>
            <h5>Discarded Emails</h5>
            <p>
                Turn logging of discarded emails on or off. When enabled, any discarded emails for your aliases are recorded and stored for 7 days.
            </p>
            <div class="mb-8">
                <label for="log-discard">
                    Log Discarded Emails:
                </label>
                <input
                    @change="saveSettings"
                    v-bind:checked="req.log_discard"
                    v-model="req.log_discard"
                    id="log-discard"
                    type="checkbox"
                >
            </div>
            <div class="mb-3">
                <button @click="saveSettings" class="cta">
                    Save Settings
                </button>
            </div>
            <p v-if="error" class="error">Error: {{ error }}</p>
            <p v-if="success" class="success">{{ success }}</p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { settingsApi } from '../api/settings.ts'
import { recipientApi } from '../api/recipient.ts'

const req = ref({
    id: '',
    domain: '',
    recipient: '',
    from_name: '',
    alias_format: '',
    log_bounce: false,
    log_discard: false,
})
const envDomains = import.meta.env.VITE_DOMAINS.split(',')
const domains = ref(envDomains)
const recipients = ref([])
const success = ref('')
const error = ref('')
const aliasFormats = ref(['Words', 'Random', 'UUID'])

const getSettings = async () => {
    try {
        const response = await settingsApi.get()
        req.value = response.data
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const saveSettings = async () => {
    req.value.domain = (document.getElementById('domain') as HTMLSelectElement).value
    req.value.recipient = (document.getElementById('recipient') as HTMLSelectElement).value
    req.value.alias_format = (document.getElementById('format') as HTMLSelectElement).value
    req.value.log_bounce = (document.getElementById('log-bounce') as HTMLInputElement).checked
    req.value.log_discard = (document.getElementById('log-discard') as HTMLInputElement).checked

    try {
        const res = await settingsApi.update(req.value)
        success.value = res.data.message
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            success.value = ''
            error.value = err.message
        }
    }
}

const getRecipients = async () => {
    try {
        const res = await recipientApi.getList()
        recipients.value = res.data
            .filter((item: { is_active: boolean }) => item.is_active)
            .map((recipient: { email: string }) => recipient.email)
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

onMounted(() => {
    getSettings()
    getRecipients()
})
</script>
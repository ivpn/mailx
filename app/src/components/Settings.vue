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
            <SettingsSelect
                v-if="loaded"
                id="domain"
                label="Select default domain:"
                button-label="Update Default Domain"
                :options="domainOptions"
                :initial-value="req.domain"
                :update="settingsApi.updateDefaultDomain"
            />
            <hr>
            <h4>Default Recipient</h4>
            <p>
                Set the default recipient for new aliases created.
            </p>
            <SettingsSelect
                v-if="loaded"
                id="recipient"
                label="Select default recipient:"
                button-label="Update Default Recipient"
                :options="recipientOptions"
                :initial-value="req.recipient"
                :update="settingsApi.updateDefaultRecipient"
            />
            <hr>
            <h4>Default Alias Format</h4>
            <p>
                Set the default alias naming format for new aliases created. Options: 1. Words ('quiet.haze16') 2. Random ('uf1h0xi') 3. UUID ('550e8400-e29b-41d4-a716-446655440000')
            </p>
            <SettingsSelect
                v-if="loaded"
                id="format"
                label="Select default alias format:"
                button-label="Update Alias Format"
                :options="aliasFormats"
                :initial-value="req.alias_format"
                :update="settingsApi.updateAliasFormat"
            />
            <hr>
            <h4>From Name</h4>
            <p>
                Set the 'From name' used for replies and emails sent using an alias. Leave it blank to use the alias email address.
            </p>
            <SettingsFromName
                v-if="loaded"
                :initial-from-name="req.from_name"
            />
            <hr>
            <h4>Mailx Header</h4>
            <p>
                Add Mailx header in forwarded messages - `Sent to &lt;alias&gt; from &lt;sender&gt;`.
            </p>
            <div v-if="loaded" class="mb-8">
                <label for="remove-header">
                    Add Mailx header:
                </label>
                <input
                    @change="saveToggle('header')"
                    v-bind:checked="!req.remove_header"
                    v-model="includeHeader"
                    id="remove-header"
                    type="checkbox"
                    class="checkbox-switch"
                >
                <p v-if="error && messageAt === 'header'" class="error">Error: {{ error }}</p>
                <p v-if="success && messageAt === 'header'" class="success">{{ success }}</p>
            </div>
            <hr>
            <SettingsBlockTrackers
                v-if="loaded"
                :initial-enabled="req.block_trackers"
            />
            <hr>
            <SettingsListUnsubscribe
                v-if="loaded"
                :initial-behaviour="req.list_unsubscribe_behaviour"
            />
            <hr>
            <h4>Diagnostic Logs</h4>
            <p>
                Track your failed email deliveries (bounces), blocked emails and forwarding issues in <router-link to="/account/diagnostics">Diagnostics</router-link>. When enabled, diagnostic logs are recorded and stored for 7 days.
            </p>
            <div v-if="loaded" class="mb-6">
                <label for="log-issues">
                    Enable Diagnostics:
                </label>
                <input
                    @change="saveToggle('diagnostics')"
                    v-bind:checked="req.log_issues"
                    v-model="req.log_issues"
                    id="log-issues"
                    type="checkbox"
                    class="checkbox-switch"
                >
                <p v-if="error && messageAt === 'diagnostics'" class="error">Error: {{ error }}</p>
                <p v-if="success && messageAt === 'diagnostics'" class="success">{{ success }}</p>
            </div>
            <div class="mb-8">
                <label for="log-issues">
                    Delete Diagnostic Logs:
                </label>
                <button @click="deleteAllLogs" class="cta sm delete">Delete All Logs</button>
                <p v-if="error && messageAt === 'logs'" class="error">Error: {{ error }}</p>
                <p v-if="success && messageAt === 'logs'" class="success">{{ success }}</p>
            </div>
            <hr>
            <p> 
                <strong>Note:</strong><br>
                If you use Mailx browser extension, in the extension go to "Settings" &gt; "Refresh Defaults" after saving to apply the new settings.
            </p>
            <p v-if="error && !messageAt" class="error">Error: {{ error }}</p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import axios from 'axios'
import { settingsApi } from '../api/settings.ts'
import { recipientApi } from '../api/recipient.ts'
import { logApi } from '../api/log.ts'
import SettingsListUnsubscribe from './SettingsListUnsubscribe.vue'
import SettingsBlockTrackers from './SettingsBlockTrackers.vue'
import SettingsSelect from './SettingsSelect.vue'
import SettingsFromName from './SettingsFromName.vue'

const req = ref({
    id: '',
    domain: '',
    recipient: '',
    from_name: '',
    alias_format: '',
    log_issues: false,
    remove_header: false,
    list_unsubscribe_behaviour: 0,
    block_trackers: true,
})
const envDomains = import.meta.env.VITE_DOMAINS.split(',')
const domains = ref<string[]>(envDomains)
const recipients = ref<string[]>([])
const success = ref('')
const error = ref('')
const aliasFormats = [
    { value: 'words', label: 'Words' },
    { value: 'random', label: 'Random' },
    { value: 'uuid', label: 'UUID' },
]
const domainOptions = computed(() => domains.value.map((domain) => ({ value: domain, label: domain })))
const recipientOptions = computed(() => recipients.value.map((email) => ({ value: email, label: email })))
const includeHeader = ref(true)
const loaded = ref(false)
// Where the latest success or error message is shown: right below the control that caused it,
// or at the bottom of the page for loading errors
type MessageAt = 'header' | 'diagnostics' | 'logs' | ''
const messageAt = ref<MessageAt>('')

const getSettings = async () => {
    try {
        const response = await settingsApi.getDefaults()
        req.value = response.data
        includeHeader.value = !req.value.remove_header
        const customDomains = response.data.custom_domains.map((item: { name: string }) => item.name)
        domains.value = [...new Set([...envDomains, ...customDomains])]
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    } finally {
        loaded.value = true
    }
}

const saveToggle = async (toggle: 'header' | 'diagnostics') => {
    messageAt.value = toggle
    success.value = ''
    error.value = ''

    try {
        const res = toggle === 'header'
            ? await settingsApi.updateRemoveHeader(!includeHeader.value)
            : await settingsApi.updateLogIssues(req.value.log_issues)
        req.value.remove_header = !includeHeader.value
        success.value = res.data.message
    } catch (err) {
        // Show the setting as it's still stored
        if (toggle === 'header') {
            includeHeader.value = !includeHeader.value
        } else {
            req.value.log_issues = !req.value.log_issues
        }
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        } else {
            error.value = 'Something went wrong'
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

const deleteAllLogs = async () => {
    if (!confirm('Are you sure you want to delete all diagnostic logs? This action cannot be undone.')) return

    messageAt.value = 'logs'
    try {
        const res = await logApi.deleteAll()
        error.value = ''
        success.value = res.data.message
    } catch (err) {
        if (axios.isAxiosError(err)) {
            success.value = ''
            error.value = err.message
        }
    }
}

onMounted(() => {
    getSettings()
    getRecipients()
})
</script>
<template>
    <div class="flex flex-col p-5 pb-4 my-8 bg-white dark:bg-neutral-800">
        <h1>Settings</h1>
        <h4>Default Domain</h4>
        <p>
            The default alias domain is the domain to be selected by default in the drop down options when generating a
            new alias.
        </p>
        <div class="max-w-xs mb-6">
            <label for="domain">
                Select default domain:
            </label>
            <select id="domain"
                class="form-select py-2.5 px-4 pe-9 block w-full border border-gray-500 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 focus:border-bluish-500 disabled:opacity-50 disabled:pointer-events-none outline-none focus:ring-transparent">
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
            The default recipient to be selected by default in the drop down options when creating a new recipient. You can add recipients <router-link to="/recipients">here</router-link>.
        </p>
        <div class="max-w-xs mb-6">
            <label for="recipient">
                Select default recipient:
            </label>
            <select id="recipient"
                :disabled="!recipients.length"
                class="form-select py-2.5 px-4 pe-9 block w-full border border-gray-500 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 focus:border-bluish-500 disabled:opacity-50 disabled:pointer-events-none outline-none focus:ring-transparent">
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
            The default alias format to be selected by default when creating a new alias. You can add aliases <router-link to="/aliases">here</router-link>.
        </p>
        <div class="max-w-xs mb-6">
            <label for="format">
                Select default alias format:
            </label>
            <select id="format"
                :disabled="!aliasFormats.length"
                class="form-select py-2.5 px-4 pe-9 block w-full border border-gray-500 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 focus:border-bluish-500 disabled:opacity-50 disabled:pointer-events-none outline-none focus:ring-transparent">
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
            The 'From Name' is shown when you send an email from an alias or reply anonymously to a forwarded email. If
            left blank, then the email alias will be used as the 'From Name'.
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
        <div class="mb-3">
            <button
                @click="saveSettings"
                class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-2 px-3 focus:outline-none focus:shadow-outline"
                type="submit">
                Save Settings
            </button>
        </div>
        <p v-if="error" class="error">Error: {{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
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
    alias_format: ''
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
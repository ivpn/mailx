<template>
    <div class="flex flex-col p-5 pb-4 my-8 bg-white dark:bg-neutral-800">
        <h1 class="text-3xl text-gray-800 dark:text-gray-100 font-semibold mb-5">Settings</h1>
        <h2 class="font-semibold text-gray-800 dark:text-gray-100 mb-3">
            Default Domain
        </h2>
        <p class="text-gray-500 dark:text-gray-400 mb-3">
            The default alias domain is the domain to be selected by default in the drop down options when generating a
            new alias.
        </p>
        <div class="max-w-xs mb-6">
            <label class="block text-gray-500 dark:text-gray-400 mb-3" for="domain">
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
        <h2 class="font-semibold text-gray-800 dark:text-gray-100 mb-3">
            Default Recipient
        </h2>
        <p class="text-gray-500 dark:text-gray-400 mb-3">
            The default recipient to be selected by default in the drop down options when creating a new recipient. You can add recipients <a class="text-bluish-500 hover:text-bluish-600" href="/recipients">here</a>.
        </p>
        <div class="max-w-xs mb-6">
            <label class="block text-gray-500 dark:text-gray-400 mb-3" for="recipient">
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
        <h2 class="font-semibold text-gray-800 dark:text-gray-100 mb-3">
            From Name
        </h2>
        <p class="text-gray-500 dark:text-gray-400 mb-3">
            The 'From Name' is shown when you send an email from an alias or reply anonymously to a forwarded email. If
            left blank, then the email alias will be used as the 'From Name'.
        </p>
        <div class="max-w-xs mb-5">
            <label class="block text-gray-500 dark:text-gray-400 mb-3" for="from-name">
                From name:
            </label>
            <input
                v-model="req.from_name"
                class="appearance-none outline-none border border-gray-500 w-full py-3 px-4 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 leading-tight focus:border-bluish-500 mb-2"
                id="from-name" type="text">
        </div>
        <div class="mb-3">
            <button
                @click="saveSettings"
                class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-2 px-3 focus:outline-none focus:shadow-outline"
                type="submit">
                Save Settings
            </button>
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-3">Error: {{ error }}</p>
        <p v-if="success" class="text-emerald-600 dark:text-emerald-500 text-sm mb-3">{{ success }}</p>
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
    from_name: ''
})
const envDomains = import.meta.env.VITE_DOMAINS.split(',')
const domains = ref(envDomains)
const recipients = ref([])
const success = ref('')
const error = ref('')

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
    const domainInput = document.getElementById('domain') as HTMLInputElement
    req.value.domain = domainInput.value

    const recipientInput = document.getElementById('recipient') as HTMLInputElement
    req.value.recipient = recipientInput.value

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
        const list = res.data.filter((item: { is_active: boolean }) => item.is_active)
        recipients.value = list.map((recipient: { email: string }) => recipient.email)
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
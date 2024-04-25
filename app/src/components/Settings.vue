<template>
    <div class="flex flex-col bg-white shadow-sm rounded-xl p-5 pb-4 my-8">
        <h1 class="text-xl font-bold text-gray-800 mb-5">Settings</h1>
        <h2 class="font-semibold text-gray-800 mb-3">
            Default Domain
        </h2>
        <p class="text-sm text-gray-500 mb-3">
            The default alias domain is the domain to be selected by default in the drop down options when generating a
            new alias.
        </p>
        <div class="max-w-xs mb-6">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="domain">
                Select default domain:
            </label>
            <select id="domain"
                class="form-select py-3 px-4 pe-9 block w-full border-2 border-gray-200 rounded-lg text-gray-700 focus:border-violet-600 disabled:opacity-50 disabled:pointer-events-none outline-none">
                <option
                    v-for="(domain, index) in domains"
                    v-bind:domain
                    :selected="domain == res.domain || index === 0"
                    :key="domain">
                    {{ domain }}
                </option>
            </select>
        </div>
        <h2 class="font-semibold text-gray-800 mb-3">
            Default Recipient
        </h2>
        <p class="text-sm text-gray-500 mb-3">
            The default recipient to be selected by default in the drop down options when creating a new recipient. You can add recipients <a class="text-violet-600 hover:text-violet-700 font-semibold" href="/recipients">here</a>.
        </p>
        <div class="max-w-xs mb-6">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="recipient">
                Select default recipient:
            </label>
            <select id="recipient"
                :disabled="!recipients.length"
                class="form-select py-3 px-4 pe-9 block w-full border-2 border-gray-200 rounded-lg text-gray-700 focus:border-violet-600 disabled:opacity-50 disabled:pointer-events-none outline-none">
                <option
                    v-for="recipient in recipients"
                    v-bind:value=recipient
                    :selected="recipient == res.recipient"
                    :key="recipient">
                    {{ recipient }}
                </option>
            </select>
        </div>
        <h2 class="font-semibold text-gray-800 mb-3">
            From Name
        </h2>
        <p class="text-sm text-gray-500 mb-3">
            The 'From Name' is shown when you send an email from an alias or reply anonymously to a forwarded email. If
            left blank, then the email alias will be used as the 'From Name'.
        </p>
        <div class="max-w-xs mb-5">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="from-name">
                From name:
            </label>
            <input
                v-model="res.from_name"
                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-violet-600 mb-2"
                id="from-name" type="text">
        </div>
        <div class="mb-3">
            <button
                @click="saveSettings"
                class="bg-violet-600 hover:bg-violet-700 text-white font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                type="submit">
                Save Settings
            </button>
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-3">{{ error }}</p>
        <p v-if="success" class="text-green-500 text-sm mb-3">{{ success }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { settingsApi } from '../api/settings.ts'
import { recipientApi } from '../api/recipient.ts'
import env from "../env.json"

const res = ref({
    id: '',
    domain: '',
    recipient: '',
    from_name: ''
})
const domains = ref(env.DOMAINS)
const recipients = ref([])
const success = ref('')
const error = ref('')

const getSettings = async () => {
    try {
        const response = await settingsApi.get()
        res.value = response.data
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const saveSettings = async () => {
    const domainInput = document.getElementById('domain') as HTMLInputElement
    res.value.domain = domainInput.value

    const recipientInput = document.getElementById('recipient') as HTMLInputElement
    res.value.recipient = recipientInput.value

    try {
        const response = await settingsApi.update(res.value)
        success.value = response.data.message
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
        const response = await recipientApi.getList()
        recipients.value = response.data.map((recipient: { email: string }) => recipient.email)
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
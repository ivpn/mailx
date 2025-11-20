<template>
    <div class="card-container">
        <header class="head">
            <h2>Logs</h2>
        </header>
        <div v-if="!logs.length && loaded" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon alert icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">You have no Logs</h4>
            <p class="text-tertiary mb-6">
                <span v-if="!settings.log_issues">To get started, first enable "Log Issues" in <router-link to="/settings">Settings</router-link>.<br></br></span>
                <span v-if="settings.log_issues">Failed email deliveries and forwarding issues will be logged here.</span>
            </p>
        </div>
        <div v-bind:class="{ 'hidden': !logs.length || !loaded }" class="card-primary">
            <div class="table-container">
                <table class="text-auto">
                    <thead>
                        <tr>
                            <th class="start">Log</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="log in logs" :key="rowKey">
                            <td class="start text-wrap leading-6">
                                <span class="text-tertiary">Type</span>: 
                                <span class="badge small" v-bind="{class: log.log_type}">{{ formatLogType(log) }}</span><br>
                                <span class="text-tertiary">ID</span>: {{ log.id }}<br>
                                <span class="text-tertiary">From</span>: {{ log.from }}<br>
                                <span class="text-tertiary">To</span>: {{ log.destination }}<br>
                                <span class="text-tertiary">Reason</span>: {{ log.message }}<br>
                                <span v-if="log.status"><span class="text-tertiary">Status</span>: {{ log.status }}<br></span>
                                <span v-if="log.remote_mta"><span class="text-tertiary">Remote MTA</span>: {{ log.remote_mta }}<br></span>
                                <span class="text-tertiary">Attempted At</span>: {{ attemptedAt(log) }}<br>
                                <button v-if="log.log_type === 'bounce'" v-bind:data-hs-overlay="'#modal-delivery-log' + log.id" class="cta mt-3">Full log</button>
                                <hr v-if="log.id !== logs[logs.length - 1]?.id" class="mt-8 mb-0">
                            </td>
                            <FailedDeliveryLog v-if="log.log_type === 'bounce'" :log="log" />
                        </tr>
                    </tbody>
                </table>
            </div>
            <p v-if="error" class="error">Error: {{ error }}</p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { settingsApi } from '../api/settings.ts'
import { logApi } from '../api/log.ts'
import FailedDeliveryLog from './FailedDeliveryLog.vue'

const log = {
    id: '',
    created_at: '',
    attempted_at: '',
    log_type: '',
    user_id: '',
    alias_id: '',
    from: '',
    destination: '',
    status: '',
    message: '',
    remote_mta: ''
}

const settings = ref({
    id: '',
    log_issues: false,
})

const logs = ref([] as typeof log[])
const error = ref('')
const loaded = ref(false)
const rowKey = ref(0)

const getSettings = async () => {
    try {
        const res = await settingsApi.get()
        settings.value = res.data
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const getLogs = async () => {
    try {
        const res = await logApi.getList()
        logs.value = res.data
        loaded.value = true
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

// Treat "0001-01-01T00:00:00Z" (and variant with .000Z) as null/unset attempted_at
const NULL_DATES = new Set([
    '0001-01-01T00:00:00Z',
    '0001-01-01T00:00:00.000Z'
])

const attemptedAt = (item: typeof log) => {
    const a = item.attempted_at
    const dateToUse = (a && !NULL_DATES.has(a)) ? a : item.created_at
    return formatDate(dateToUse)
}

const formatDate = (date: string) => {
    const d = new Date(date)
    return d.toLocaleString()
}

const formatLogType = (item: typeof log) => {
    if (!item.log_type) return ''
    return item.log_type.replace('_', ' ').charAt(0).toUpperCase() + item.log_type.replace('_', ' ').slice(1)
}

onMounted(async () => {
    await getSettings()
    getLogs()
})
</script>
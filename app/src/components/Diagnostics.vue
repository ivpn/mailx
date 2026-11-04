<template>
    <div class="card-container">
        <header class="head">
            <h2>Diagnostics</h2>
        </header>
        <div v-if="!logs.length && loaded" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon alert icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">No Diagnostic logs</h4>
            <p class="text-tertiary mb-6">
                <span v-if="!settings.log_issues">Enable email delivery diagnostics in <router-link to="/account/settings">Settings</router-link>.<br></br></span>
                <span v-if="settings.log_issues">Email delivery diagnostics will be logged here.</span>
            </p>
        </div>
        <div v-bind:class="{ 'hidden': !logs.length || !loaded }" class="card-primary">
            <div class="table-container">
                <table class="text-auto">
                    <thead>
                        <tr>
                            <th class="start">
                                <div class="hs-dropdown [--placement:bottom-left]">
                                    <button id="hs-dropdown-log-type" class="sort">
                                        <span class="flex mr-1">Show:</span>
                                        <span class="text-accent flex items-center">
                                            {{ logTypeLabel }}
                                        </span>
                                    </button>
                                    <div
                                        class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                                        aria-labelledby="hs-dropdown-log-type"
                                    >
                                        <button @click="setLogType('')">All</button>
                                        <button v-for="type in LOG_TYPES" :key="type" @click="setLogType(type)">{{ formatType(type) }}</button>
                                    </div>
                                </div>
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="!filteredLogs.length">
                            <td class="text-center text-tertiary py-10">No logs found</td>
                        </tr>
                        <tr v-for="log in filteredLogs" :key="log.id">
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
                                <button v-if="log.log_type === 'bounce'" v-bind:data-hs-overlay="'#modal-delivery-log' + log.id" class="cta sm mt-3">Full log</button>
                                <hr v-if="log.id !== filteredLogs[filteredLogs.length - 1]?.id" class="mt-8 mb-0">
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
import { computed, nextTick, onMounted, ref } from 'vue'
import axios from 'axios'
import { settingsApi } from '../api/settings.ts'
import { logApi } from '../api/log.ts'
import FailedDeliveryLog from './FailedDeliveryLog.vue'
import { closeDropdowns, initDropdowns } from '../lib/preline.ts'

// Mirrors model.LogType in the API.
const LOG_TYPES = [
    'bounce',
    'deferred_delivery',
    'failed_delivery',
    'disabled_alias',
    'disabled_domain',
    'paused_alias',
    'no_recipients',
    'unauthorised_send',
    'inactive_subscription',
    'blocked_sender',
    'blocked_by_rule',
]

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
const logType = ref('')
const logTypeLabel = computed(() => logType.value ? formatType(logType.value) : 'All')

// The API returns every log at once, so the type filter is applied client-side.
const filteredLogs = computed(() => logType.value
    ? logs.value.filter(l => l.log_type === logType.value)
    : logs.value)

const setLogType = (value: string) => {
    if (value === logType.value) return
    closeDropdowns(false)
    logType.value = value
}

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
        await nextTick()
        initDropdowns()
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

const formatType = (type: string) => {
    const label = type.replace(/_/g, ' ')
    return label.charAt(0).toUpperCase() + label.slice(1)
}

const formatLogType = (item: typeof log) => {
    if (!item.log_type) return ''
    return formatType(item.log_type)
}

onMounted(async () => {
    await getSettings()
    getLogs()
})
</script>
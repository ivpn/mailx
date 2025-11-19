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
                To get started, first enable "Log Issues" in <router-link to="/settings">Settings</router-link>.
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
                            <span class="text-tertiary">ID</span>: {{ log.id }}<br>
                            <span class="text-tertiary">From</span>: {{ log.from }}<br>
                            <span class="text-tertiary">To</span>: {{ log.destination }}<br>
                            <span class="text-tertiary">Type</span>: {{ log.type }}<br>
                            <span class="text-tertiary">Reason</span>: {{ log.message }}<br></br>
                            <span class="text-tertiary">Status</span>: {{ log.status }}<br>
                            <span class="text-tertiary">Remote MTA</span>: {{ log.remote_mta }}<br>
                            <span class="text-tertiary">Attempted At</span>: {{ formatDate(log.attempted_at || log.created_at) }}<br>
                            <button v-bind:data-hs-overlay="'#modal-delivery-log' + log.id" class="cta mt-3">Full log</button>
                            <hr v-if="log.id !== logs[logs.length - 1]?.id" class="mt-8 mb-0">
                        </td>
                        <FailedDeliveryLog :log="log" />
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
import { logApi } from '../api/log.ts'
import FailedDeliveryLog from './FailedDeliveryLog.vue'

const log = {
    id: '',
    created_at: '',
    attempted_at: '',
    type: '',
    user_id: '',
    alias_id: '',
    from: '',
    destination: '',
    status: '',
    message: '',
    remote_mta: ''
}

const logs = ref([] as typeof log[])
const error = ref('')
const loaded = ref(false)
const rowKey = ref(0)

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

const formatDate = (date: string) => {
    const d = new Date(date)
    return d.toLocaleString()
}

onMounted(async () => {
    getLogs()
})
</script>
<template>
    <div class="card-container">
        <header class="head">
            <h2>Logs</h2>
        </header>
        <div v-if="!list.length && loaded" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon alert icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">You have no Failed Deliveries yet</h4>
            <p class="text-tertiary mb-6">
                To get started, first enable Failed Deliveries in <router-link to="/settings">Settings</router-link>.
            </p>
        </div>
        <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="card-primary">
            <div class="table-container">
                <table class="text-auto">
                    <thead>
                        <tr>
                            <th class="start">Log</th>
                        </tr>
                    </thead>
                    <tbody>
                    <tr v-for="bounce in list" :key="rowKey">
                        <td class="start text-wrap leading-6">
                            <span class="text-tertiary">ID</span>: {{ bounce.id }}<br>
                            <span class="text-tertiary">From</span>: {{ bounce.from }}<br>
                            <span class="text-tertiary">To</span>: {{ bounce.destination }}<br>
                            <span class="text-tertiary">Status</span>: {{ bounce.status }}<br>
                            <span class="text-tertiary">Remote MTA</span>: {{ bounce.remote_mta }}<br>
                            <span class="text-tertiary">Diagnostic Code</span>: {{ bounce.diagnostic_code }}<br>
                            <span class="text-tertiary">Attempted At</span>: {{ formatDate(bounce.attempted_at) }}<br>
                            <button v-bind:data-hs-overlay="'#modal-delivery-log' + bounce.id" class="cta mt-3">Full log</button>
                        </td>
                        <FailedDeliveryLog :log="bounce" />
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
import { bounceApi } from '../api/bounce.ts'
import FailedDeliveryLog from './FailedDeliveryLog.vue'

const bounce = {
    id: '',
    created_at: '',
    attempted_at: '',
    user_id: '',
    alias_id: '',
    from: '',
    destination: '',
    status: '',
    diagnostic_code: '',
    remote_mta: ''
}

const list = ref([] as typeof bounce[])
const error = ref('')
const loaded = ref(false)
const rowKey = ref(0)

const getList = async () => {
    try {
        const res = await bounceApi.getList()
        list.value = res.data
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
    getList()
})
</script>
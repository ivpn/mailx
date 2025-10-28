<template>
    <div class="card-container">
        <header class="head">
            <h2>Failed Deliveries</h2>
        </header>
        <div v-if="!list.length && loaded" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon at icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">You have no Failed Deliveries yet</h4>
            <p class="text-tertiary mb-6">
                To get started, first enable Failed Deliveries in <router-link to="/settings">Settings</router-link>.
            </p>
        </div>
        <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="card-primary">
            <div  class="table-container">
                <table>
                    <thead class="desktop">
                        <tr>
                            <th>Log</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                    <tr v-for="bounce in list" :key="rowKey" class="desktop">
                        <td>
                            <span class="text-tertiary">ID</span>: {{ bounce.id }}<br>
                            <span class="text-tertiary">From</span>: {{ bounce.from }}<br>
                            <span class="text-tertiary">To</span>: {{ bounce.destination }}<br>
                            <span class="text-tertiary">Status</span>: {{ bounce.status }}<br>
                            <span class="text-tertiary">Remote MTA</span>: {{ bounce.remote_mta }}<br>
                            <span class="text-tertiary">Diagnostic Code</span>: {{ bounce.diagnostic_code }}<br>
                            <span class="text-tertiary">Created At</span>: {{ bounce.created_at }}<br>
                            <span class="text-tertiary">Attempted At</span>: {{ bounce.attempted_at }}
                        </td>
                        <td>
                            GET
                        </td>
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

onMounted(async () => {
    getList()
})
</script>
<template>
    <div class="card-container">

        <header class="head">
            <h2>Logs</h2>
        </header>

        <nav class="flex gap-x-2" aria-label="Tabs" role="tablist">
            <button type="button" class="plain active" id="unstyled-tabs-item-1" aria-selected="true" data-hs-tab="#unstyled-tabs-1" aria-controls="unstyled-tabs-1" role="tab">
                Failed Deliveries
            </button>
            <button type="button" class="plain" id="unstyled-tabs-item-2" aria-selected="false" data-hs-tab="#unstyled-tabs-2" aria-controls="unstyled-tabs-2" role="tab">
                Discarded Emails
            </button>
        </nav>

        <div class="mt-6 card-fill">
            <div id="unstyled-tabs-1" role="tabpanel" aria-labelledby="unstyled-tabs-item-1" class="card-fill">
                <div v-if="!bounces.length && loaded" class="card-empty">
                    <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                        <i class="icon alert icon-accent text-2xl"></i>
                    </span>
                    <h4 class="mb-6">You have no Failed Deliveries yet</h4>
                    <p class="text-tertiary mb-6">
                        To get started, first enable Failed Deliveries in <router-link to="/settings">Settings</router-link>.
                    </p>
                </div>
                <div v-bind:class="{ 'hidden': !bounces.length || !loaded }" class="card-primary">
                    <div class="table-container">
                        <table class="text-auto">
                            <thead>
                                <tr>
                                    <th class="start">Log</th>
                                </tr>
                            </thead>
                            <tbody>
                            <tr v-for="bounce in bounces" :key="rowKey">
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
                    <p v-if="errorBounces" class="error">Error: {{ errorBounces }}</p>
                </div>
            </div>
            <div id="unstyled-tabs-2" class="card-fill hidden" role="tabpanel" aria-labelledby="unstyled-tabs-item-2">
                <div v-if="!discards.length && loaded" class="card-empty">
                    <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                        <i class="icon alert icon-accent text-2xl"></i>
                    </span>
                    <h4 class="mb-6">You have no Discarded Emails yet</h4>
                    <p class="text-tertiary mb-6">
                        To get started, first enable Discarded Emails in <router-link to="/settings">Settings</router-link>.
                    </p>
                </div>
                <div v-bind:class="{ 'hidden': !discards.length || !loaded }" class="card-primary">
                    <div class="table-container">
                        <table class="text-auto">
                            <thead>
                                <tr>
                                    <th class="start">Log</th>
                                </tr>
                            </thead>
                            <tbody>
                            <tr v-for="discard in discards" :key="rowKey">
                                <td class="start text-wrap leading-6">
                                    <span class="text-tertiary">ID</span>: {{ discard.id }}<br>
                                    <span class="text-tertiary">From</span>: {{ discard.from }}<br>
                                    <span class="text-tertiary">To</span>: {{ discard.destination }}<br>
                                    <span class="text-tertiary">Reason</span>: {{ discard.message }}<br>
                                    <span class="text-tertiary">Attempted At</span>: {{ formatDate(discard.created_at) }}<br>
                                    <button v-bind:data-hs-overlay="'#modal-delivery-log' + discard.id" class="cta mt-3">Full log</button>
                                </td>
                            </tr>
                        </tbody>
                        </table>
                    </div>
                    <p v-if="errorDiscards" class="error">Error: {{ errorDiscards }}</p>
                </div>
            </div>
        </div>

    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { bounceApi } from '../api/bounce.ts'
import { discardApi } from '../api/discard.ts'
import FailedDeliveryLog from './FailedDeliveryLog.vue'
import tabs from '@preline/tabs'

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

const discard = {
    id: '',
    created_at: '',
    user_id: '',
    alias_id: '',
    from: '',
    destination: '',
    message: ''
}

const bounces = ref([] as typeof bounce[])
const discards = ref([] as typeof discard[])
const errorBounces = ref('')
const errorDiscards = ref('')
const loaded = ref(false)
const rowKey = ref(0)

const getBounces = async () => {
    try {
        const res = await bounceApi.getList()
        bounces.value = res.data
        loaded.value = true
        errorBounces.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            errorBounces.value = err.message
        }
    }
}

const getDiscards = async () => {
    try {
        const res = await discardApi.getList()
        discards.value = res.data
        loaded.value = true
        errorDiscards.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            errorDiscards.value = err.message
        }
    }
}

const formatDate = (date: string) => {
    const d = new Date(date)
    return d.toLocaleString()
}

onMounted(async () => {
    getBounces()
    getDiscards()
    tabs.autoInit()
})
</script>
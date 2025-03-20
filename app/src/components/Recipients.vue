<template>
    <div class="card-container">
        <header class="head">
            <h2>Recipients</h2>
            <div class="flex items-center justify-between">
                <RecipientCreate />
            </div>
        </header>
        <div v-if="!list.length && loaded" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon mailbox icon-accent text-2xl"></i>
            </span>
            <h4>You did not create any recipients yet</h4>
            <p class="text-tertiary pb-6">
                This area will light up with items once you start creating recipients.
            </p>
            <RecipientCreate />
        </div>
        <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="card-primary">
            <div class="table-container">
                <table>
                    <thead>
                        <tr>
                            <th>Created</th>
                            <th>Email</th>
                            <th>Verified</th>
                            <th>Encryption</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <RecipientRow v-for="recipient in list" :recipient="recipient" :key="rowKey" />
                    </tbody>
                </table>
                <p v-if="error" class="error">Error: {{ error }}</p>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { recipientApi } from '../api/recipient.ts'
import RecipientRow from './RecipientRow.vue'
import RecipientCreate from './RecipientCreate.vue'
import events from '../events.ts'

const recipient = {
    id: '',
    created_at: '',
    email: '',
    is_active: false,
}

const props = defineProps(['dashboard'])
const isDashboard = props.dashboard
const list = ref([] as typeof recipient[])
const error = ref('')
const loaded = ref(false)
const rowKey = ref(0)

const getList = async () => {
    try {
        const response = await recipientApi.getList()
        list.value = response.data
        if (isDashboard) list.value = list.value.slice(0, 5)
        loaded.value = true
        error.value = ''
        renderRow()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const deleteRecipient = async (id: string) => {
    try {
        await recipientApi.delete(id)
        error.value = ''
        reload()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const onDeleteRecipient = (payload: { id: string }) => {
    deleteRecipient(payload.id)
}

const reload = () => {
    if (!isDashboard) getList()
    if (isDashboard) location.reload()
}

const renderRow = () => {
    rowKey.value++
}

onMounted(() => {
    getList()
    events.on('recipient.create', getList)
    events.on('recipient.update', reload)
    events.on('recipient.verify', reload)
    events.on('recipient.delete', onDeleteRecipient)
})
</script>
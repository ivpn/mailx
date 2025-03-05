<template>
    <div v-if="!list.length && loaded" class="flex flex-col my-14">
        <div class="flex flex-col items-center text-center">
            <h3>Add Recipients</h3>
            <p>To get started, add a recipient.</p>
            <div class="flex gap-4">
                <RecipientCreate />
            </div>
        </div>
    </div>
    <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="card">
        <h1>Recipients</h1>
        <div>
            <div class="flex items-center justify-between mb-6">
                <RecipientCreate />
            </div>
            <div class="flex flex-col">
                <div class="-m-1.5 overflow-x-auto">
                    <div class="p-1.5 min-w-full inline-block align-middle">
                        <div class="overflow-x-auto">
                            <table class="table-auto w-full divide-y divide-gray-200 dark:divide-neutral-600">
                                <thead>
                                    <tr>
                                        <th scope="col"
                                            class="pr-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                            CREATED</th>
                                        <th scope="col"
                                            class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                            EMAIL</th>
                                        <th scope="col"
                                            class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                            VERIFIED</th>
                                            <th scope="col"
                                            class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                            ENCRYPTION</th>
                                        <th scope="col"
                                            class="pl-5 py-3 text-end text-xs font-medium text-gray-500 dark:text-gray-400">
                                            ACTIONS</th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-200 dark:divide-neutral-600">
                                    <RecipientRow v-for="recipient in list" :recipient="recipient" :key="rowKey" />
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <p v-if="isDashboard" class="my-4">
            <router-link to="/recipients">All Recipients</router-link>
        </p>
        <p v-if="error" class="error">Error: {{ error }}</p>
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
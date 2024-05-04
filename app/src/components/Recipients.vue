<template>
    <div v-if="!list.length && loaded" class="flex flex-col p-5 pb-4 my-8">
        <div class="flex flex-col items-center p-4 text-center py-5">
            <h3 class="text-lg font-bold text-gray-800">
                No recipients yet
            </h3>
            <p class="my-2 text-gray-500">
                To get started, add a recipient.
            </p>
            <div class="flex gap-4">
                <RecipientCreate @onCreateRecipient="getList" />
            </div>
        </div>
    </div>
    <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="flex flex-col bg-white shadow-sm rounded-xl p-5 pb-4 my-8">
        <h1 v-if="!isDashboard" class="text-2xl font-bold text-gray-800 mb-5">Recipients</h1>
        <h1 v-if="isDashboard" class="text-2xl font-bold text-gray-800 mb-5">Latest Recipients</h1>
        <div>
            <div class="flex items-center justify-between mb-6">
                <RecipientCreate @onCreateRecipient="getList" />
            </div>
            <div class="flex flex-col">
                <div class="-m-1.5 overflow-x-auto">
                    <div class="p-1.5 min-w-full inline-block align-middle">
                        <div class="overflow-hidden">
                            <table class="min-w-full divide-y divide-gray-200">
                                <thead>
                                    <tr>
                                        <th scope="col"
                                            class="pr-5 py-3 text-start text-xs font-medium text-gray-500">
                                            CREATED</th>
                                        <th scope="col"
                                            class="px-5 py-3 text-start text-xs font-medium text-gray-500">
                                            EMAIL</th>
                                        <th scope="col"
                                            class="px-5 py-3 text-start text-xs font-medium text-gray-500">
                                            VERIFIED</th>
                                        <th scope="col"
                                            class="pl-5 py-3 text-end text-xs font-medium text-gray-500">
                                            ACTIONS</th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-200">
                                    <RecipientRow  @onDeleteRecipient="deleteRecipient" @onVerifyRecipient="getList" v-for="recipient in list" :recipient="recipient" :key="rowKey" />
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-4">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { recipientApi } from '../api/recipient.ts'
import RecipientRow from './RecipientRow.vue'
import RecipientCreate from './RecipientCreate.vue'

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
    if (!confirm('Are you sure you want to delete recipient?')) return
    try {
        await recipientApi.delete(id)
        error.value = ''
        getList()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const renderRow = () => {
    rowKey.value++
}

onMounted(() => {
    getList()
})
</script>
<template>
    <div v-if="!list.length && loaded" class="flex flex-col my-14">
        <div class="flex flex-col items-center text-center">
            <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100">
                Add Recipients
            </h3>
            <p class="my-2 text-gray-500 dark:text-gray-400">
                To get started, add a recipient.
            </p>
            <div class="flex gap-4">
                <RecipientCreate @onCreateRecipient="getList" />
            </div>
        </div>
    </div>
    <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="flex flex-col p-5 pb-4 my-8 bg-white dark:bg-neutral-800">
        <h1 v-if="!isDashboard" class="text-3xl text-gray-800 dark:text-gray-100 font-semibold mb-5">Recipients</h1>
        <h1 v-if="isDashboard" class="text-3xl text-gray-800 dark:text-gray-100 font-semibold mb-5">Latest Recipients</h1>
        <div>
            <div class="flex items-center justify-between mb-6">
                <RecipientCreate @onCreateRecipient="getList" />
            </div>
            <div class="flex flex-col">
                <div class="-m-1.5 overflow-x-auto">
                    <div class="p-1.5 min-w-full inline-block align-middle">
                        <div class="overflow-hidden">
                            <table class="min-w-full divide-y divide-gray-200 dark:divide-neutral-600">
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
                                            class="pl-5 py-3 text-end text-xs font-medium text-gray-500 dark:text-gray-400">
                                            ACTIONS</th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-200">
                                    <RecipientRow  @onDeleteRecipient="deleteRecipient" @onVerifyRecipient="reload" v-for="recipient in list" :recipient="recipient" :key="rowKey" />
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <p v-if="isDashboard" class="text-sm text-gray-500 my-4">
            <a href="/recipients" class="text-bluish-500 hover:text-bluish-600 font-medium text-sm py-2"
                type="submit">All Recipients</a>
        </p>
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
    if (!confirm('Are you sure you want to delete recipient? Note that aliases with this recipient will be disabled.')) return
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

const reload = () => {
    if (!isDashboard) getList()
    if (isDashboard) location.reload()
}

const renderRow = () => {
    rowKey.value++
}

onMounted(() => {
    getList()
})
</script>
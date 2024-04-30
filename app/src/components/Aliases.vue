<template>
    <div v-if="!list.length && loaded" class="flex flex-col p-5 pb-4 my-8">
        <div class="flex flex-col items-center p-4 text-center py-20">
            <h3 class="text-lg font-bold text-gray-800">
                No aliases yet
            </h3>
            <p class="my-2 text-gray-500">
                To get started, create an alias.
            </p>
            <div class="flex gap-4">
                <AliasCreate v-if="recipients.length" @onCreateAlias="getList" :recipients.sync="recipients" />
            </div>
        </div>
    </div>
    <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="flex flex-col bg-white shadow-sm rounded-xl p-5 pb-4 my-8">
        <h1 class="text-2xl font-bold text-gray-800 mb-5">Aliases</h1>
        <div>
            <div class="flex items-center justify-between mb-6">
                <AliasCreate v-if="recipients.length && settings.recipient" @onCreateAlias="getList" :recipients.sync="recipients" :settings.sync="settings" />
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
                                            ALIAS</th>
                                        <th scope="col"
                                            class="px-5 py-3 text-start text-xs font-medium text-gray-500">
                                            RECIPIENTS</th>
                                        <th scope="col"
                                            class="px-5 py-3 text-start text-xs font-medium text-gray-500">
                                            COUNT
                                        </th>
                                        <th scope="col"
                                            class="px-5 py-3 text-start text-xs font-medium text-gray-500">
                                            ACTIVE</th>
                                        <th scope="col"
                                            class="pl-5 py-3 text-end text-xs font-medium text-gray-500">
                                            ACTIONS</th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-200">
                                    <AliasCard @onDeleteAlias="deleteAlias" v-for="alias in list" :alias="alias" :key="alias.id" :recipients.sync="recipients" />
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-4">{{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { aliasApi } from '../api/alias'
import { recipientApi } from '../api/recipient.ts'
import { settingsApi } from '../api/settings.ts'
import AliasCard from './AliasCard.vue'
import AliasCreate from './AliasCreate.vue'

const alias = {
    id: '',
    created_at: '',
    name: '',
    enabled: false,
    description: '',
    recipients: '',
    from_name: '',
    stats: {
        forwards: 0,
        blocks: 0,
        replies: 0,
        sends: 0
    }
}

const list = ref([] as typeof alias[])
const recipients = ref([])
const settings = ref({
    id: '',
    domain: '',
    recipient: '',
    from_name: ''
})
const error = ref('')
const loaded = ref(false)

const getList = async () => {
    try {
        const response = await aliasApi.getList()
        list.value = response.data
        loaded.value = true
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const getRecipients = async () => {
    try {
        const response = await recipientApi.getList()
        recipients.value = response.data.map((recipient: { email: string }) => recipient.email)
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const getSettings = async () => {
    try {
        const response = await settingsApi.get()
        settings.value = response.data
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const deleteAlias = async (id: string) => {
    if (!confirm('Are you sure you want to delete this alias?')) return
    try {
        await aliasApi.delete(id)
        error.value = ''
        getList()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

onMounted(() => {
    getList()
    getRecipients()
    getSettings()
})

</script>
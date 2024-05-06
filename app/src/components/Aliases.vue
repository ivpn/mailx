<template>
    <div v-if="!list.length && loaded" class="flex flex-col my-14">
        <div class="flex flex-col items-center text-center">
            <h3 class="text-lg font-bold text-gray-800">
                Create Aliases
            </h3>
            <p v-if="recipients.length && settings.recipient" class="my-2 text-gray-500">
                To get started, create an alias.
            </p>
            <p v-if="!recipients.length && loaded" class="my-2 text-gray-500">
                To get started, first add a recipient.
            </p>
            <div class="flex gap-4">
                <AliasCreate v-if="recipients.length && settings.recipient" @onCreateAlias="getList" :recipients.sync="recipients" :settings.sync="settings" />
            </div>
        </div>
    </div>
    <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="flex flex-col bg-white shadow-sm rounded-xl p-5 pb-4 my-8">
        <h1 v-if="!isDashboard" class="text-2xl font-bold text-gray-800 mb-5">Aliases</h1>
        <h1 v-if="isDashboard" class="text-2xl font-bold text-gray-800 mb-5">Latest Aliases</h1>
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
                                    <AliasRow @onDeleteAlias="deleteAlias" @onEditAlias="getList" v-for="alias in list" :alias="alias" :key="rowKey" :recipients.sync="recipients" />
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <p v-if="isDashboard" class="text-sm text-gray-500 my-4">
            <a href="/aliases" class="text-blue-600 hover:text-blue-700 font-medium text-sm py-2"
                type="submit">All Aliases</a>
        </p>
        <p v-if="error" class="text-red-600 text-sm mb-4">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { aliasApi } from '../api/alias'
import { recipientApi } from '../api/recipient.ts'
import { settingsApi } from '../api/settings.ts'
import AliasRow from './AliasRow.vue'
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

const props = defineProps(['dashboard'])
const isDashboard = props.dashboard
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
const rowKey = ref(0)

const getList = async () => {
    try {
        const response = await aliasApi.getList()
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
    if (!confirm('Are you sure you want to delete alias?')) return
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

const renderRow = () => {
    rowKey.value++
}

onMounted(() => {
    getList()
    getRecipients()
    getSettings()
})

</script>
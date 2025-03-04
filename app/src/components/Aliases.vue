<template>

    <!-- Standard Aliases -->
    <div>
        <div v-if="!list.length && loaded" class="flex flex-col my-14">
            <div class="flex flex-col items-center text-center">
                <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100">
                    Create Aliases
                </h3>
                <p v-if="recipients.length && settings.id" class="my-2 text-gray-500 dark:text-gray-400">
                    To get started, create an alias.
                </p>
                <p v-if="!recipients.length && loaded" class="my-2 text-gray-500 dark:text-gray-400">
                    To get started, first add a recipient.
                </p>
                <div class="flex gap-4">
                    <AliasCreate v-if="recipients.length && settings.id" :recipients.sync="recipients" :settings.sync="settings" :catchAll=false />
                </div>
            </div>
        </div>
        <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="flex flex-col p-5 pb-4 my-8 bg-white dark:bg-neutral-800">
            <h1>Aliases</h1>
            <div>
                <div class="flex items-center justify-between mb-6">
                    <AliasCreate v-if="recipients.length && settings.id" :recipients.sync="recipients" :settings.sync="settings" :catchAll=false />
                </div>
                <div class="flex flex-col">
                    <div class="-m-1.5 overflow-x-auto">
                        <div class="p-1.5 min-w-full inline-block align-middle">
                            <div class="overflow-x-auto">
                                <table class="table-auto w-full divide-y divide-gray-200 dark:divide-neutral-600">
                                    <thead>
                                        <tr>
                                            <th v-if="!isDashboard" scope="col" class="pr-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                <button
                                                @click="sort"
                                                data-sort="created_at"
                                                class="inline-flex justify-center items-center">
                                                    CREATED
                                                    <svg
                                                    data-sort="created_at"
                                                    v-bind:class="{ 'text-bluish-500': sortBy === 'created_at', 'rotate-180': sortOrder === 'ASC' && sortBy === 'created_at' }"
                                                    class="ms-1 flex-shrink-0 size-5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                                        stroke-linecap="round" stroke-linejoin="round">
                                                        <path d="m6 9 6 6 6-6" />
                                                    </svg>
                                                </button>
                                            </th>
                                            <th v-if="!isDashboard" scope="col" class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                <button
                                                @click="sort"
                                                data-sort="name"
                                                class="inline-flex justify-center items-center">
                                                    ALIAS
                                                    <svg
                                                    data-sort="name"
                                                    v-bind:class="{ 'text-bluish-500': sortBy === 'name', 'rotate-180': sortOrder === 'ASC' && sortBy === 'name' }"
                                                    class="ms-1 flex-shrink-0 size-5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                                        stroke-linecap="round" stroke-linejoin="round">
                                                        <path d="m6 9 6 6 6-6" />
                                                    </svg>
                                                </button>    
                                            </th>
                                            <th v-if="isDashboard" scope="col" class="pr-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                CREATED
                                            </th>
                                            <th v-if="isDashboard" scope="col" class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                ALIAS
                                            </th>
                                            <th scope="col" class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                COUNT
                                            </th>
                                            <th scope="col" class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                ACTIVE
                                            </th>
                                            <th scope="col" class="pl-5 py-3 text-end text-xs font-medium text-gray-500 dark:text-gray-400">
                                                ACTIONS
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody class="divide-y divide-gray-200 dark:divide-neutral-600">
                                        <AliasRow v-for="alias in list" :alias="alias" :key="rowKey" :recipients.sync="recipients" :catchAll=false />
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <p v-if="isDashboard" class="text-sm text-gray-500 dark:text-gray-300 my-4">
                <router-link to="/aliases" class="cta-plain"
                    type="submit">All Aliases</router-link>
            </p>
            <p v-if="error" class="text-red-600 text-sm mb-4">Error: {{ error }}</p>
            <Pagination v-if="list.length && !isDashboard" :list.sync="list" :limit="limit" :page="page" :total="total" :key="rowKey" @onUpdatePage="onUpdatePage" />
        </div>
    </div>

    <!-- Catch-all Aliases -->
    <div v-if="!isDashboard">
        <div v-if="!listCatchAll.length && loadedCatchAll" class="flex flex-col my-14">
            <div class="flex flex-col items-center text-center">
                <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100 mb-5">
                    Create Catch-all Aliases
                </h3>
                <p v-if="!recipients.length && loadedCatchAll" class="my-2 text-gray-500 dark:text-gray-400">
                    To get started, first add a recipient.
                </p>
                <div class="flex gap-4">
                    <AliasCreate v-if="recipients.length && settings.id" :recipients.sync="recipients" :settings.sync="settings" :catchAll=true />
                </div>
            </div>
        </div>
        <div v-bind:class="{ 'hidden': !listCatchAll.length || !loadedCatchAll }" class="flex flex-col p-5 pb-4 my-8 bg-white dark:bg-neutral-800">
            <h1>Catch-all Aliases</h1>
            <div>
                <div class="flex items-center justify-between mb-6">
                    <AliasCreate v-if="recipients.length && settings.id" :recipients.sync="recipients" :settings.sync="settings" :catchAll=true />
                </div>
                <div class="flex flex-col">
                    <div class="-m-1.5 overflow-x-auto">
                        <div class="p-1.5 min-w-full inline-block align-middle">
                            <div class="overflow-x-auto">
                                <table class="table-auto w-full divide-y divide-gray-200 dark:divide-neutral-600">
                                    <thead>
                                        <tr>
                                            <th scope="col" class="pr-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                CREATED
                                            </th>
                                            <th scope="col" class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                ALIAS
                                            </th>
                                            <th scope="col" class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                COUNT
                                            </th>
                                            <th scope="col" class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                                ACTIVE
                                            </th>
                                            <th scope="col" class="pl-5 py-3 text-end text-xs font-medium text-gray-500 dark:text-gray-400">
                                                ACTIONS
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody class="divide-y divide-gray-200 dark:divide-neutral-600">
                                        <AliasRow v-for="alias in listCatchAll" :alias="alias" :key="rowKey" :recipients.sync="recipients" :catchAll=true />
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
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
import Pagination from './Pagination.vue'
import events from '../events.ts'

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
const listCatchAll = ref([] as typeof alias[])
const recipients = ref([])
const settings = ref({
    id: '',
    domain: '',
    recipient: '',
    from_name: ''
})
const error = ref('')
const loaded = ref(false)
const loadedCatchAll = ref(false)
const rowKey = ref(0)
const limit = ref(25)
const page = ref(1)
const total = ref(0)
const sortBy = ref('created_at')
const sortOrder = ref('DESC')

const getList = async () => {
    try {
        const res = await aliasApi.getList({
            limit: limit.value,
            page: page.value,
            sort_by: sortBy.value,
            sort_order: sortOrder.value,
            catch_all: "false"
        })
        list.value = res.data.aliases
        total.value = res.data.total
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

const getListCatchAll = async () => {
    try {
        const res = await aliasApi.getList({
            limit: limit.value,
            page: page.value,
            sort_by: sortBy.value,
            sort_order: sortOrder.value,
            catch_all: "true"
        })
        listCatchAll.value = res.data.aliases
        loadedCatchAll.value = true
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const getRecipients = async () => {
    try {
        const res = await recipientApi.getList()
        const list = res.data.filter((item: { is_active: boolean }) => item.is_active)
        recipients.value = list.map((recipient: { email: string }) => recipient.email)
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const getSettings = async () => {
    try {
        const res = await settingsApi.get()
        settings.value = res.data
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const deleteAlias = async (payload: any) => {
    try {
        await aliasApi.delete(payload.id)
        error.value = ''
        fetch()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const renderRow = () => {
    rowKey.value++
}

const onUpdatePage = (obj: any) => {
    limit.value = obj.limit
    page.value = obj.page
    getList()
}

const onDeleteAlias = (payload: { id: string, catchAll: boolean }) => {
    deleteAlias(payload)
}

const sort = (e: any) => {
    const sort = e.target.dataset.sort
    if (sort === sortBy.value) {
        sortOrder.value = sortOrder.value === 'ASC' ? 'DESC' : 'ASC'
    } else {
        sortBy.value = sort
        sortOrder.value = 'DESC'
    }

    getList()
}

const fetch = () => {
    getList()
    getListCatchAll()
}

onMounted(async () => {
    await getRecipients()
    await getSettings()
    fetch()
    events.on('alias.create', fetch)
    events.on('alias.update', fetch)
    events.on('alias.delete', onDeleteAlias)
})

</script>
<template>
    <div class="card-container">
        <header class="head">
            <h2>Aliases</h2>
            <div class="flex gap-3 items-center justify-between">
                <div class="max-md:hidden relative">
                    <form v-if="loaded" @submit.prevent="getList" autocomplete="off">
                        <input class="search" type="text" v-model="search" placeholder="Search aliases...">
                    </form>
                    <button v-if="searchQuery" @click.prevent="clearSearch" class="absolute top-0 right-0 bottom-0 px-2 flex items-center justify-center">
                        <i class="icon close icon-tertiary text-base"></i>
                    </button>
                </div>
                <button v-if="recipients.length" class="cta text-nowrap" data-hs-overlay="#modal-create-alias-false">
                    New Alias
                </button>
            </div>
        </header>
        <div class="mb-7 tablet">
            <div class="relative">
                <form v-if="loaded" @submit.prevent="getList" autocomplete="off">
                    <input class="search" type="text" v-model="search" placeholder="Search aliases...">
                </form>
                <button v-if="searchQuery" @click.prevent="clearSearch" class="absolute top-0 right-0 bottom-0 px-2 flex items-center justify-center">
                    <i class="icon close icon-tertiary text-base"></i>
                </button>
            </div>
        </div>
        <div v-if="!list.length && loaded" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon at icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">
                <span v-if="!searchQuery && !loading">You have no aliases yet</span>
                <span v-if="searchQuery || loading">No aliases found</span>
            </h4>
            <p v-if="!recipients.length" class="text-tertiary mb-6">
                To get started, first <router-link to="/account">verify</router-link> your primary email address.
            </p>
             <button v-if="!searchQuery && !loading && recipients.length" class="cta" data-hs-overlay="#modal-create-alias-false">
                New Alias
            </button>
        </div>
        <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="card-primary">
            <div  class="table-container">
                <table>
                    <thead class="desktop">
                        <tr>
                            <th>Status</th>
                            <th>
                                <button
                                @click="sort"
                                data-sort="name"
                                class="sort">
                                    Alias
                                    <i
                                        data-sort="name"
                                        v-if="sortBy !== 'name'"
                                        v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'name' }"
                                        class="icon arrow-down text-xl icon-tertiary"
                                    ></i>
                                    <i
                                        data-sort="name"
                                        v-if="sortBy === 'name'"
                                        v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'name' }"
                                        class="icon arrow-down text-xl icon-accent"
                                    ></i>
                                </button>    
                            </th>
                            <th>Domain</th>
                            <th>Count</th>
                            <th>
                                <button
                                @click="sort"
                                data-sort="created_at"
                                class="sort">
                                    Created
                                    <i
                                        data-sort="created_at"
                                        v-if="sortBy !== 'created_at'"
                                        v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'created_at' }"
                                        class="icon arrow-down text-xl icon-tertiary"
                                    ></i>
                                    <i
                                        data-sort="created_at"
                                        v-if="sortBy === 'created_at'"
                                        v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'created_at' }"
                                        class="icon arrow-down text-xl icon-accent"
                                    ></i>
                                </button>
                            </th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <AliasRow v-for="alias in list" :alias="alias" :key="rowKey" :recipients.sync="recipients" :catchAll=false />
                    </tbody>
                </table>
            </div>
            <p v-if="error" class="error">Error: {{ error }}</p>
            <Pagination v-if="list.length" :list.sync="list" :limit="limit" :page="page" :total="total" :key="rowKey" @onUpdatePage="onUpdatePage" />
        </div>
    </div>
    <AliasCreate v-if="recipients.length && settings.id" :recipients.sync="recipients" :settings.sync="settings" :catchAll=false :label="'New Alias'" />
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
import { RouterLink } from 'vue-router'

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
const loading = ref(false)
const rowKey = ref(0)
const limit = ref(25)
const page = ref(1)
const total = ref(0)
const sortBy = ref('created_at')
const sortOrder = ref('DESC')
const search = ref('')
const searchQuery = ref('')

const getList = async () => {
    loading.value = true
    searchQuery.value = search.value.trim()
    if (searchQuery.value) {
        page.value = 1 // Reset to first page on search
    }

    try {
        const res = await aliasApi.getList({
            limit: limit.value,
            page: page.value,
            sort_by: sortBy.value,
            sort_order: sortOrder.value,
            catch_all: false,
            search: searchQuery.value
        })
        list.value = res.data.aliases
        total.value = res.data.total
        loaded.value = true
        loading.value = false
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
}

const clearSearch = () => {
    search.value = ''
    searchQuery.value = ''
    getList()
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
<template>
    <div>
        <div v-if="!list.length && loaded" class="flex flex-col my-14">
            <div class="flex flex-col items-center text-center">
                <h3>Create Aliases</h3>
                <p v-if="recipients.length && settings.id" class="my-2">
                    To get started, create an alias.
                </p>
                <p v-if="!recipients.length && loaded" class="my-2">
                    To get started, first add a recipient.
                </p>
                <div class="flex gap-4">
                    <AliasCreateDropdown />
                </div>
            </div>
        </div>
        <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="card-container">
            <div class="flex flex-row justify-between">
                <h2>Aliases</h2>
                <div class="flex items-center justify-between mb-6">
                    <AliasCreateDropdown />
                </div>
            </div>
            <div class="card-primary">
                <div class="table-container">
                    <table>
                        <thead>
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
                                <th>
                                    <button
                                    @click="sort"
                                    data-sort="catch_all"
                                    class="sort">
                                        Type
                                        <i
                                            data-sort="catch_all"
                                            v-if="sortBy !== 'catch_all'"
                                            v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'catch_all' }"
                                            class="icon arrow-down text-xl icon-tertiary"
                                        ></i>
                                        <i
                                            data-sort="catch_all"
                                            v-if="sortBy === 'catch_all'"
                                            v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'catch_all' }"
                                            class="icon arrow-down text-xl icon-accent"
                                        ></i>
                                    </button>    
                                </th>
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
    </div>
    <AliasCreate v-if="recipients.length && settings.id" :recipients.sync="recipients" :settings.sync="settings" :catchAll=false :label="'New Alias'" />
    <AliasCreate v-if="recipients.length && settings.id" :recipients.sync="recipients" :settings.sync="settings" :catchAll=true :label="'New Catch-all Alias'" />
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
import AliasCreateDropdown from './AliasCreateDropdown.vue'
import events from '../events.ts'
import dropdown from '@preline/dropdown'

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
            sort_order: sortOrder.value
        })
        list.value = res.data.aliases
        total.value = res.data.total
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

onMounted(async () => {
    dropdown.autoInit()
    await getRecipients()
    await getSettings()
    fetch()
    events.on('alias.create', fetch)
    events.on('alias.update', fetch)
    events.on('alias.delete', onDeleteAlias)
})
</script>
<template>
    <div class="card-container">
        <header class="head">
            <h2>Wildcard Aliases</h2>
            <div class="flex items-center justify-between">
                <button v-if="recipients.length && loaded" class="cta" data-hs-overlay="#modal-create-alias-true">
                    New Wildcard
                </button>
            </div>
        </header>
        <div v-if="showEmptyCard" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon at icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">You have no wildcards yet</h4>
            <p v-if="!recipients.length && loaded" class="text-tertiary mb-6">
                To get started, first <router-link to="/account/profile">verify</router-link> your primary email address.
            </p>
             <button v-if="recipients.length && loaded" class="cta" data-hs-overlay="#modal-create-alias-true">
                New Wildcard
            </button>
        </div>
        <div v-bind:class="{ 'hidden': showEmptyCard || !loaded }">
            <div class="tablet-lg">
                <div class="hs-dropdown [--placement:bottom-left] mb-2">
                    <button id="hs-dropdown-wildcard-status-mobile" class="sort">
                        <span class="flex mr-1">Show:</span>
                        <span class="text-accent flex items-center">
                            {{ statusLabel }}
                        </span>
                    </button>
                    <div
                        class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                        aria-labelledby="hs-dropdown-wildcard-status-mobile"
                    >
                        <button @click="setStatus('active_inactive')">All</button>
                        <button @click="setStatus('active')">Enabled</button>
                        <button @click="setStatus('inactive')">Disabled</button>
                        <button @click="setStatus('deleted')">Deleted</button>
                    </div>
                </div>
            </div>
            <div class="card-primary">
                <div class="table-container">
                    <table>
                        <thead class="desktop-lg">
                            <tr>
                                <th>
                                    <div class="hs-dropdown [--placement:bottom-left]">
                                        <button id="hs-dropdown-wildcard-status" class="sort">
                                            <span class="flex mr-1">Show:</span>
                                            <span class="text-accent flex items-center">
                                                {{ statusLabel }}
                                            </span>
                                        </button>
                                        <div
                                            class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                                            aria-labelledby="hs-dropdown-wildcard-status"
                                        >
                                            <button @click="setStatus('active_inactive')">All</button>
                                            <button @click="setStatus('active')">Enabled</button>
                                            <button @click="setStatus('inactive')">Disabled</button>
                                            <button @click="setStatus('deleted')">Deleted</button>
                                        </div>
                                    </div>
                                </th>
                                <th>Description</th>
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
                            <AliasRow v-for="alias in list" :alias="alias" :key="alias.id" :recipients.sync="recipients" :wildcard=true @onEdit="onEditAlias" @onSend="onSendAlias" />
                            <!-- Also keeps <tbody> from collapsing to its bare divide-y border, which
                                 Chrome reads as a 1px overflow and answers with a stray scrollbar. -->
                            <tr v-if="!list.length && loaded">
                                <td colspan="6" class="text-center text-tertiary py-10">No wildcards found</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <p v-if="error" class="error">Error: {{ error }}</p>
                <Pagination v-if="list.length" :list.sync="list" :limit="limit" :page="page" :total="total" :key="limit + '-' + page + '-' + total" @onUpdatePage="onUpdatePage" />
            </div>
        </div>
    </div>
    <AliasCreate v-if="recipients.length && settings.id && loaded" :recipients.sync="recipients" :settings.sync="settings" :wildcard=true :label="'New Wildcard Alias'" />
    <AliasEdit v-if="recipients.length" ref="editModal" :recipients="recipients" :key="recipients.join(',')" />
    <AliasSend ref="sendModal" />
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed, nextTick } from 'vue'
import axios from 'axios'
import { aliasApi } from '../api/alias'
import { settingsApi } from '../api/settings.ts'
import AliasRow from './AliasRow.vue'
import AliasCreate from './AliasCreate.vue'
import AliasEdit from './AliasEdit.vue'
import AliasSend from './AliasSend.vue'
import Pagination from './Pagination.vue'
import events from '../events.ts'
import { RouterLink } from 'vue-router'
import { closeDropdowns, initDropdowns, initOverlays, initTooltips } from '../lib/preline.ts'

const alias = {
    id: '',
    created_at: '',
    deleted_at: null as string | null,
    name: '',
    enabled: false,
    description: '',
    recipients: '',
    from_name: '',
    pinned: false,
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
    domains: [],
    custom_domains: [],
    recipient: '',
    from_name: ''
})
const error = ref('')
const loaded = ref(false)
const loading = ref(false)
const limit = ref(25)
const page = ref(1)
const total = ref(0)
const sortBy = ref('created_at')
const sortOrder = ref('DESC')
const status = ref('active_inactive')
const editModal = ref<InstanceType<typeof AliasEdit> | null>(null)
const sendModal = ref<InstanceType<typeof AliasSend> | null>(null)
const statusLabel = computed(() => {
    if (status.value === 'active') return 'Enabled'
    if (status.value === 'inactive') return 'Disabled'
    if (status.value === 'deleted') return 'Deleted'
    return 'All'
})

// Only the default filter gets the full-page empty card; every other filter keeps the table (and
// with it the status dropdown) on screen so the user can pick their way back out of it. Holding
// off while a fetch is in flight stops the card from flashing over a list that is about to fill.
const showEmptyCard = computed(() => loaded.value && !loading.value && !list.value.length && status.value === 'active_inactive')

const getList = async () => {
    loading.value = true
    try {
        const res = await aliasApi.getList({
            limit: limit.value,
            page: page.value,
            sort_by: sortBy.value,
            sort_order: sortOrder.value,
            wildcard: true,
            status: status.value,
        })
        list.value = res.data.aliases
        total.value = res.data.total
        loaded.value = true
        error.value = ''
        // Removing the last row of a page (by filtering it out, deleting it, ...) would otherwise
        // strand the user on an empty page, with the pagination they'd leave it by hidden too.
        if (!list.value.length && page.value > 1) {
            page.value = 1
            return getList()
        }
        loading.value = false
        await nextTick()
        bindPreline()
    } catch (err) {
        loading.value = false
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const getSettings = async () => {
    try {
        const res = await settingsApi.getDefaults()
        settings.value = res.data
        recipients.value = res.data.recipients
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

// Bound once per list render; the row components deliberately do not init Preline themselves.
const bindPreline = () => {
    initTooltips()
    initDropdowns()
    initOverlays()
}

const onUpdatePage = (obj: any) => {
    limit.value = obj.limit
    page.value = obj.page
    getList()
}

const onDeleteAlias = (payload: { id: string, wildcard: boolean }) => {
    deleteAlias(payload)
}

// Flipping a row's switch changes which aliases belong in the list, but only while the list is
// filtered on that state.
const onAliasEnabled = () => {
    if (status.value === 'active' || status.value === 'inactive') {
        getList()
    }
}

const onEditAlias = (alias: any) => {
    editModal.value?.open(alias)
}

const onSendAlias = (alias: any) => {
    sendModal.value?.open(alias)
}

const setStatus = (value: string) => {
    if (value === status.value) return
    // The dropdown sits inside the block the new status may hide, and Preline jams for good if
    // its menu is hidden before the closing transition ends, so close it outright first.
    closeDropdowns(false)
    status.value = value
    page.value = 1
    getList()
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
    await getSettings()
    fetch()
    initDropdowns()
    events.on('alias.create', fetch)
    events.on('alias.update', fetch)
    events.on('alias.enabled', onAliasEnabled)
    events.on('alias.delete', onDeleteAlias)
})

onUnmounted(() => {
    events.off('alias.create', fetch)
    events.off('alias.update', fetch)
    events.off('alias.enabled', onAliasEnabled)
    events.off('alias.delete', onDeleteAlias)
})
</script>
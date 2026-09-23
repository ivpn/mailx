<template>
    <div class="card-container">
        <header class="head">
            <h2>Blocklist</h2>
        </header>
        <div class="card-primary grow-0 mb-6">
            <h4>Add to blocklist</h4>
            <p>
                Emails from blocked senders are discarded before they reach your recipients. Blocking a domain also blocks its subdomains. With email delivery diagnostics enabled in <router-link to="/account/settings">Settings</router-link>, blocked emails are listed in <router-link to="/account/diagnostics">Diagnostics</router-link>.
            </p>
            <form @submit.prevent="addItem" autocomplete="off" class="flex flex-col md:flex-row gap-3">
                <select v-model="req.type" class="md:w-40 md:shrink-0 mb-0" aria-label="Type">
                    <option value="email">Email</option>
                    <option value="domain">Domain</option>
                </select>
                <input
                    v-model="req.value"
                    v-bind:class="{ 'error': valueError }"
                    v-bind:placeholder="req.type === 'domain' ? 'example.com' : 'name@example.com'"
                    id="blocklist_value"
                    aria-label="Value"
                    type="text"
                >
                <button type="submit" class="cta justify-center text-nowrap" v-bind:disabled="submitting">
                    Add to blocklist
                </button>
            </form>
            <p v-if="valueError" class="error">Required</p>
            <p v-if="formError" class="error">Error: {{ formError }}</p>
        </div>
        <div v-if="!list.length && loaded" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon block icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">Your blocklist is empty</h4>
            <p class="text-tertiary mb-6">Blocked email addresses and domains will be listed here.</p>
        </div>
        <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="card-primary">
            <div class="table-container">
                <table>
                    <thead class="desktop-lg">
                        <tr>
                            <th>Created</th>
                            <th>Value</th>
                            <th>Type</th>
                            <th>Blocked</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <BlocklistRow v-for="item in list" :item="item" :key="item.id" @onDelete="deleteItem" />
                    </tbody>
                </table>
            </div>
            <p v-if="error" class="error">Error: {{ error }}</p>
            <Pagination v-if="list.length" :limit="limit" :page="page" :total="total" :key="limit + '-' + page + '-' + total" @onUpdatePage="onUpdatePage" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue'
import axios from 'axios'
import { blocklistApi } from '../api/blocklist.ts'
import BlocklistRow from './BlocklistRow.vue'
import Pagination from './Pagination.vue'
import { initDropdowns } from '../lib/preline.ts'
import { useBreakpoint } from '../lib/useBreakpoint.ts'

const blocklistItem = {
    id: '',
    created_at: '',
    type: '',
    value: '',
    blocked: 0,
}

const list = ref([] as typeof blocklistItem[])
const req = ref({
    type: 'email',
    value: '',
})
const valueError = ref(false)
const formError = ref('')
const submitting = ref(false)
const error = ref('')
const loaded = ref(false)
const limit = ref(25)
const page = ref(1)
const total = ref(0)
const { isDesktop } = useBreakpoint()

const getList = async () => {
    try {
        const res = await blocklistApi.getList({
            limit: limit.value,
            page: page.value,
        })
        list.value = res.data.items
        total.value = res.data.total
        error.value = ''
        // Deleting the last row of a page would otherwise strand the user on an empty page.
        if (!list.value.length && page.value > 1) {
            page.value = 1
            return getList()
        }
        await nextTick()
        initDropdowns()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    } finally {
        loaded.value = true
    }
}

const addItem = async () => {
    valueError.value = !req.value.value.trim()
    if (valueError.value) return

    submitting.value = true
    try {
        await blocklistApi.create({
            type: req.value.type,
            value: req.value.value.trim(),
        })
        req.value.value = ''
        formError.value = ''
        page.value = 1
        await getList()
        document.getElementById('blocklist_value')?.focus()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            formError.value = err.response?.data.error || err.message
        }
    } finally {
        submitting.value = false
    }
}

const deleteItem = async (item: typeof blocklistItem) => {
    try {
        await blocklistApi.delete(item.id)
        error.value = ''
        getList()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}

const onUpdatePage = (obj: any) => {
    limit.value = obj.limit
    page.value = obj.page
    getList()
}

// Rows swap between their desktop and mobile markup at the breakpoint, and the new dropdowns
// need binding just like freshly fetched ones.
watch(isDesktop, () => {
    nextTick(() => initDropdowns())
})

onMounted(() => {
    getList()
})
</script>

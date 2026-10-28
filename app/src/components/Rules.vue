<template>
    <div class="card-container">
        <header class="head">
            <h2>Rules</h2>
            <div class="flex items-center justify-between">
                <button @click="openCreate" v-bind:disabled="readOnly" class="cta" data-test="new-rule">
                    New Rule
                </button>
            </div>
        </header>
        <p class="mt-0 mb-6">
            Rules automatically block, forward or modify emails sent to your aliases, based on conditions you set. They run from top to bottom, and every matching rule applies, until one marked "Stops" matches. Drag rules to change the order.
        </p>
        <div v-if="readOnly" class="card-tertiary mb-6" data-test="read-only-notice">
            <footer class="mt-0">
                <div>
                    <i class="icon info icon-primary"></i>
                </div>
                <div>
                    <p>
                        Your subscription is not active. Your existing rules are still applied, but you can't create, edit, delete or reorder rules until you <a target="_blank" :href="activateUrl">renew</a>.
                    </p>
                </div>
            </footer>
        </div>
        <p v-if="error" class="error" data-test="rules-error">Error: {{ error }}</p>
        <div v-if="!list.length && loaded" class="card-empty" data-test="empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon filter icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">You don't have any rules yet</h4>
            <button @click="openCreate" v-bind:disabled="readOnly" class="cta">
                New Rule
            </button>
        </div>
        <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="card-primary">
            <div class="table-container">
                <table>
                    <thead class="desktop-lg">
                        <tr>
                            <th><span class="sr-only">Order</span></th>
                            <th>Created</th>
                            <th>Name</th>
                            <th>Active</th>
                            <th>Applied</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <VueDraggable
                        v-model="list"
                        v-bind:disabled="readOnly"
                        v-bind:animation="150"
                        ghost-class="opacity-40"
                        @start="onDragStart"
                        @end="onDragEnd"
                        tag="tbody"
                        handle=".rule-handle"
                    >
                        <RuleRow
                            v-for="rule in list"
                            :key="rule.id"
                            :rule="rule"
                            :read-only="readOnly"
                            @onToggle="toggleRule"
                            @onEdit="openEdit"
                            @onDuplicate="openDuplicate"
                            @onDelete="confirmDelete"
                            @onMove="moveRule"
                        />
                    </VueDraggable>
                </table>
            </div>
        </div>
        <RuleEdit ref="ruleEdit" :recipients="recipients" @saved="onSaved" @inactiveSub="reload" />
        <!-- The modal is opened from code, but Preline only binds a modal that has a toggle -->
        <button data-hs-overlay="#modal-delete-rule" type="button" class="hidden" tabindex="-1" aria-hidden="true"></button>
        <div id="modal-delete-rule" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="closeDelete" type="button" class="close" aria-label="Close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>DELETE RULE</h4>
                    </header>
                    <article>
                        <p class="break-all">
                            Are you sure you want to delete the rule <strong>{{ deleting?.name }}</strong>? This can't be undone.
                        </p>
                    </article>
                    <footer>
                        <nav>
                            <button @click="deleteRule" type="button" class="cta delete" data-test="confirm-delete">
                                Delete Rule
                            </button>
                            <button @click="closeDelete" type="button" class="cancel">
                                Cancel
                            </button>
                        </nav>
                        <p v-if="deleteError" class="error px-5">Error: {{ deleteError }}</p>
                    </footer>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import axios from 'axios'
import overlay from '@preline/overlay'
import { VueDraggable } from 'vue-draggable-plus'
import { ruleApi } from '../api/rule.ts'
import { recipientApi } from '../api/recipient.ts'
import { subscriptionApi } from '../api/subscription.ts'
import RuleEdit from './RuleEdit.vue'
import RuleRow from './RuleRow.vue'
import { initDropdowns, initOverlays } from '../lib/preline.ts'
import { useBreakpoint } from '../lib/useBreakpoint.ts'
import { type Rule, isInactiveSubError, isStaleRulesError } from '../lib/rules.ts'

// Matches ActiveStatus() on the server
const activeStatuses = ['active', 'grace_period']
const activateUrl = import.meta.env.VITE_RESYNC_URL

const list = ref<Rule[]>([])
const recipients = ref<{ id: string, email: string }[]>([])
const status = ref('')
const loaded = ref(false)
const error = ref('')
const deleting = ref<Rule | null>(null)
const deleteError = ref('')
const ruleEdit = ref<InstanceType<typeof RuleEdit> | null>(null)
const { isDesktop } = useBreakpoint()

const readOnly = computed(() => status.value !== '' && !activeStatuses.includes(status.value))

const errorMessage = (err: unknown) => {
    if (axios.isAxiosError(err)) return err.response?.data?.error || err.message
    return String(err)
}

const getList = async () => {
    try {
        const res = await ruleApi.getList()
        list.value = res.data
        await nextTick()
        initDropdowns()
    } catch (err) {
        error.value = errorMessage(err)
    } finally {
        loaded.value = true
    }
}

const getSubscription = async () => {
    try {
        const res = await subscriptionApi.get()
        status.value = res.data.status || ''
    } catch {}
}

const getRecipients = async () => {
    try {
        const res = await recipientApi.getList()
        recipients.value = res.data
            .filter((item: { is_active: boolean }) => item.is_active)
            .map((item: { id: string, email: string }) => ({ id: item.id, email: item.email }))
    } catch {}
}

// The subscription, or the rules themselves, changed since the page was loaded
const reload = () => {
    getSubscription()
    getList()
}

const handleError = (err: unknown) => {
    error.value = errorMessage(err)
    if (isInactiveSubError(error.value)) reload()
}

const onSaved = () => {
    error.value = ''
    getList()
}

const openCreate = () => {
    if (readOnly.value) return
    ruleEdit.value?.open(null)
}

const openEdit = (rule: Rule) => {
    if (readOnly.value) return
    ruleEdit.value?.open(rule)
}

const openDuplicate = (rule: Rule) => {
    if (readOnly.value) return
    ruleEdit.value?.open(rule, true)
}

const toggleRule = async (rule: Rule) => {
    if (readOnly.value) return

    const previous = rule.active
    rule.active = !previous
    try {
        await ruleApi.updateActive(rule.id, rule.active)
        error.value = ''
    } catch (err) {
        rule.active = previous
        handleError(err)
    }
}

const saveOrder = async (previous: Rule[]) => {
    try {
        await ruleApi.reorder(list.value.map(rule => rule.id))
        error.value = ''
    } catch (err) {
        list.value = previous
        handleError(err)
        if (isStaleRulesError(error.value)) getList()
    }
}

let orderBeforeDrag: Rule[] = []

const onDragStart = () => {
    orderBeforeDrag = [...list.value]
}

const onDragEnd = () => {
    const previous = orderBeforeDrag
    orderBeforeDrag = []
    if (previous.length === list.value.length && previous.every((rule, i) => rule.id === list.value[i].id)) return

    saveOrder(previous)
}

const moveRule = (rule: Rule, delta: number) => {
    if (readOnly.value) return

    const from = list.value.findIndex(item => item.id === rule.id)
    const to = from + delta
    if (from < 0 || to < 0 || to >= list.value.length) return

    const previous = [...list.value]
    const next = [...list.value]
    next.splice(to, 0, next.splice(from, 1)[0])
    list.value = next
    saveOrder(previous)

    // Moving the row's element drops its focus, so the handle is focused again for the next key press
    nextTick(() => {
        document.querySelector<HTMLElement>(`.rule-handle[data-rule-id="${rule.id}"]`)?.focus()
    })
}

const deleteModal = () => document.getElementById('modal-delete-rule') as HTMLElement

const confirmDelete = (rule: Rule) => {
    if (readOnly.value) return

    deleting.value = rule
    deleteError.value = ''
    overlay.open(deleteModal())
}

const closeDelete = () => {
    overlay.close(deleteModal())
}

const deleteRule = async () => {
    if (!deleting.value) return

    try {
        await ruleApi.delete(deleting.value.id)
        error.value = ''
        closeDelete()
        getList()
    } catch (err) {
        deleteError.value = errorMessage(err)
        if (isInactiveSubError(deleteError.value)) reload()
    }
}

// Rows swap between their desktop and mobile markup at the breakpoint, and the new dropdowns
// need binding just like freshly fetched ones.
watch(isDesktop, () => {
    nextTick(() => initDropdowns())
})

onMounted(async () => {
    initOverlays()
    await Promise.all([getSubscription(), getRecipients()])
    getList()
})
</script>

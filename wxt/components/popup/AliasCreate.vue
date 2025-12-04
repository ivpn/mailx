<template>
    <div>
        <div id="modal-create-alias" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4 class="uppercase">
                            New Alias
                        </h4>
                    </header>
                    <article>
                        <div>
                            <div class="pb-3">
                                <label for="alias_description">
                                    Description (optional)
                                </label>
                                <input
                                    v-model="alias.description"
                                    id="alias_description"
                                    type="text"
                                >
                            </div>
                            <div class="pb-3">
                                <label for="create-alias-recipient">
                                    Recipient(s)
                                </label>
                                <select
                                    id="create-alias-recipient"
                                    v-model="selectRecipients"
                                    :multiple="true"
                                    data-hs-select='{
                                    "placeholder": "Select recipient(s)",
                                    "toggleTag": "<button type=\"button\" aria-expanded=\"false\"></button>",
                                    "toggleClasses": "hs-select-disabled:pointer-events-none hs-select-disabled:opacity-50 relative py-2.5 ps-4 pe-9 flex gap-x-2 text-nowrap w-full cursor-pointer border border-primary text-secondary bg-secondary leading-tight focus:border-accent",
                                    "dropdownClasses": "mt-2 z-50 w-full max-h-72 p-1 space-y-0.5 bg-primary border border-tertiary overflow-hidden overflow-y-auto",
                                    "optionClasses": "py-2 px-4 w-full text-secondary cursor-pointer hover:bg-secondary bg-primary",
                                    "optionTemplate": "<div class=\"flex justify-between items-center w-full\"><span data-title></span><span class=\"hidden hs-selected:block\"><svg class=\"shrink-0 size-3.5 text-accent \" xmlns=\"http:.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><polyline points=\"20 6 9 17 4 12\"/></svg></span></div>",
                                    "extraMarkup": "<div class=\"absolute top-1/2 end-3 -translate-y-1/2\"><svg class=\"shrink-0 size-3.5 text-secondary \" xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"m7 15 5 5 5-5\"/><path d=\"m7 9 5-5 5 5\"/></svg></div>"
                                    }' class="hidden">
                                    <option v-for="(recipient, index) in defaults.recipients"
                                        v-bind:value=recipient
                                        :selected="recipient == defaults.recipient || (!defaults.recipient && index === 0)"
                                        :key="recipient">
                                        {{ recipient }}
                                    </option>
                                </select>
                                <p v-if="errorRecipients" class="error pt-3">{{ errorRecipients }}</p>
                            </div>
                            <div class="hs-accordion-group">
                                <div class="hs-accordion" id="alias-accordion-one">
                                    <button class="hs-accordion-toggle plain font-bold flex flex-row gap-2 mb-2 items-center" aria-expanded="false" aria-controls="alias-accordion-collapse-one">
                                        <svg class="hs-accordion-active:hidden inline-flex size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                            <path d="M5 12h14"></path>
                                            <path d="M12 5v14"></path>
                                        </svg>
                                        <svg class="hs-accordion-active:inline-flex hidden size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                            <path d="M5 12h14"></path>
                                        </svg>
                                        Custom settings
                                    </button>
                                    <div id="alias-accordion-collapse-one" class="hs-accordion-content hidden overflow-hidden transition-[height] duration-300" role="region" aria-labelledby="alias-accordion-one">
                                        <div class="pb-3">
                                            <label for="alias_format">
                                                Format
                                            </label>
                                            <select id="alias_format">
                                                <option v-for="(format, index) in formats" v-bind:value="format.value"
                                                    :selected="format.value == alias.format || index === 0" :key="format.value">
                                                    {{ format.name }}
                                                </option>
                                            </select>
                                        </div>
                                        <div class="pb-3">
                                            <label for="alias_domain">
                                                Domain
                                            </label>
                                            <select id="alias_domain" :disabled="!defaults.domains.length">
                                                <option v-for="(domain, index) in defaults.domains" v-bind:domain
                                                    :selected="domain == alias.domain || index === 0" :key="domain">
                                                    {{ domain }}
                                                </option>
                                            </select>
                                        </div>
                                        <div class="pb-3">
                                            <label for="alias_from_name" class="flex gap-2">
                                                From name
                                                <span class="hs-tooltip [--strategy:absolute] flex">
                                                    <i class="icon info icon-primary hs-tooltip-toggle"></i>
                                                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible text-center" role="tooltip">
                                                        Leave blank to use<br>
                                                        alias email address or<br>
                                                        default From Name defined<br>
                                                        in Settings (if set)
                                                    </span>
                                                </span>
                                            </label>
                                            <input
                                                v-model="alias.from_name"
                                                id="alias_from_name"
                                                type="text"
                                            >
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button
                                v-bind:disabled="errorRecipients.length > 0"
                                @click="postAlias"
                                class="cta">
                                Create and copy to clipboard
                            </button>
                            <button @click="close" class="cancel">
                                Cancel
                            </button>
                        </nav>
                        <p v-if="error" class="error px-5">Error: {{ error }}</p>
                    </footer>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import overlay from '@preline/overlay'
import select from '@preline/select'
import { api } from '@/lib/api'
import events from '@/lib/events'
import tooltip from '@preline/tooltip'
import accordion from '@preline/accordion'
import { Defaults, Alias } from '@/lib/types'

const props = defineProps<{
    apiToken: string
    defaults: Defaults
}>()
const alias = ref({} as Alias)
const selectRecipients = ref([] as string[])
const formats = ref([{
    name: 'Words',
    value: 'words'
}, {
    name: 'Random',
    value: 'random'
}, {
    name: 'UUID',
    value: 'uuid'
}])
const error = ref('')
const errorRecipients = ref('')
const loading = ref(false)

const postAlias = async () => {
    if (loading.value) return

    alias.value.domain = (document.getElementById('alias_domain') as HTMLInputElement).value
    alias.value.recipients = selectRecipients.value.join(',')
    alias.value.enabled = true
    
    const formatElement = document.getElementById('alias_format') as HTMLInputElement;
    if (formatElement) {
        alias.value.format = formatElement.value;
    }

    if (!validate(alias.value.recipients)) return

    try {
        loading.value = true
        const res = await api.createAlias(props.apiToken, alias.value)
        console.log('Created alias:', res)
        copyAlias(res.name)
        error.value = ''
        events.emit('alias.create', {})
        close()
    } catch (err) {
        error.value = 'An unexpected error occurred'
        console.error('Create alias error:', err)
    } finally {
        loading.value = false
    }
}

const close = () => {
    resetAlias()
    error.value = ''

    document.removeEventListener('keydown', handleKeydown)

    const accordionInstance = accordion.getInstance('#alias-accordion-one', true);
    if (accordionInstance && 'element' in accordionInstance) {
        accordionInstance.element.hide()
    }

    const modal = document.querySelector('#modal-create-alias') as any
    overlay.close(modal)

    const multiselect = select.getInstance('#create-alias-recipient' as any, true) as any
    multiselect.element.setValue([props.defaults.recipient])
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-create-alias' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
    modal.element.on('open', () => {
        document.addEventListener('keydown', handleKeydown)
        focusFirstInput()
    })

    const multiselect = select.getInstance('#create-alias-recipient' as any, true) as any
    multiselect.element.on('change', (val: any) => {
        errorRecipients.value = val.length === 0 ? 'Select one or more recipients' : ''
    })
}

const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Enter') {
        event.preventDefault()
        postAlias()
    }
}

const focusFirstInput = () => {
    const input = document.getElementById('alias_description')
    input?.focus()
}

const validate = (rcps: string) => {
    if (rcps.length === 0) {
        error.value = 'Select one or more recipients'
        return false
    }

    return true
}

const resetAlias = () => {
    const multiselect = select.getInstance('#create-alias-recipient' as any, true) as any
    multiselect.element.setValue([props.defaults.recipient])
    selectRecipients.value = [props.defaults.recipient]
    alias.value = {} as Alias
    alias.value.domain = props.defaults.domain
}

const copyAlias = (alias: string) => {
    navigator.clipboard.writeText(alias)
}

onMounted(async () => {
    overlay.autoInit()
    select.autoInit()
    tooltip.autoInit()
    accordion.autoInit()
    addEvents()
    resetAlias()
})
</script>

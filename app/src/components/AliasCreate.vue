<template>
    <div>
        <div v-bind:id="'modal-create-alias-' + props.catchAll" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4 class="uppercase">
                            {{ props.label }}
                        </h4>
                        <div v-if="props.catchAll" class="hs-tooltip [--strategy:absolute]">
                            <i class="icon info icon-primary hs-tooltip-toggle"></i>
                            <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">Limited to 2 wildcard aliases per domain</span>
                        </div>
                    </header>
                    <article>
                        <div v-if="props.catchAll">
                            <div class="mb-3">
                                <label for="alias_catch_all_suffix">
                                    Alias suffix (6-12 alphanumeric chars.):
                                </label>
                                <input 
                                    v-model="alias.catch_all_suffix"
                                    v-bind:class="{ 'error': errorCatchAllSuffix }"
                                    id="alias_catch_all_suffix"
                                    type="text"
                                >
                                <p v-if="errorCatchAllSuffix" class="error">Wildcard suffix must be between 6 and 12 characters</p>
                                <p class="text-primary mb-1">
                                    *+{{ alias.catch_all_suffix }}@{{ alias.domain }}
                                </p>
                            </div>
                        </div>
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
                                    :disabled="!recipients.length"
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
                                    <option v-for="(recipient, index) in recipients"
                                        v-bind:value=recipient
                                        :selected="recipient == settings.recipient || (!settings.recipient && index === 0)"
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
                                        <div v-if="!props.catchAll" class="pb-3">
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
                                            <select id="alias_domain" :disabled="!domains.length">
                                                <option v-for="(domain, index) in domains" v-bind:domain
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
                                                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">Leave blank to use alias email address or default From Name defined in Settings (if set)</span>
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
import axios from 'axios'
import { aliasApi } from '../api/alias.ts'
import events from '../events.ts'
import tooltip from '@preline/tooltip'
import accordion from '@preline/accordion'

const envDomains = import.meta.env.VITE_DOMAINS.split(',')
const props = defineProps(['recipients', 'settings', 'catchAll', 'label'])
const alias = ref({
    description: '',
    enabled: true,
    format: '',
    from_name: '',
    recipients: '',
    domain: envDomains[0],
    catch_all: props.catchAll ? 'true' : 'false',
    catch_all_suffix: ''
})
const recipients = ref(props.recipients)
const settings = ref(props.settings)
const selectRecipients = ref([settings.value.recipient ? settings.value.recipient : props.recipients[0]])
const domains = ref(envDomains)
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
const errorCatchAllSuffix = ref(false)
const loading = ref(false)

const postAlias = async () => {
    if (loading.value) return

    alias.value.domain = (document.getElementById('alias_domain') as HTMLInputElement).value
    alias.value.recipients = selectRecipients.value.join(',')
    alias.value.enabled = true
    
    if (props.catchAll) {
        alias.value.format = 'catch_all'
    } else {
        const formatElement = document.getElementById('alias_format') as HTMLInputElement;
        if (formatElement) {
            alias.value.format = formatElement.value;
        }
    }

    if (!validate(alias.value.recipients)) return

    try {
        loading.value = true
        const res = await aliasApi.create(alias.value)
        copyAlias(res.data.name)
        events.emit('alias.create', {})
        error.value = ''
        loading.value = false
        close()
    } catch (err) {
        loading.value = false
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const close = () => {
    resetAlias()
    error.value = ''
    errorCatchAllSuffix.value = false

    document.removeEventListener('keydown', handleKeydown)

    const accordionInstance = accordion.getInstance('#alias-accordion-one', true);
    if (accordionInstance && 'element' in accordionInstance) {
        accordionInstance.element.hide()
    }

    const modal = document.querySelector('#modal-create-alias-' + props.catchAll) as any
    overlay.close(modal)

    const multiselect = select.getInstance('#create-alias-recipient' as any, true) as any
    multiselect.element.setValue([settings.value.recipient ? settings.value.recipient : props.recipients[0]])
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-create-alias-' + props.catchAll as any, true) as any
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
    const input = props.catchAll ? document.getElementById('alias_catch_all_suffix') : document.getElementById('alias_description')
    input?.focus()
}

const validate = (rcps: string) => {
    if (rcps.length === 0) {
        error.value = 'Select one or more recipients'
        return false
    }

    if (props.catchAll) {
        errorCatchAllSuffix.value = alias.value.catch_all_suffix.length < 6 || alias.value.catch_all_suffix.length > 12
    } else {
        errorCatchAllSuffix.value = false
    }

    return !errorCatchAllSuffix.value
}

const resetAlias = () => {
    alias.value = {
        description: '',
        enabled: true,
        format: props.settings.alias_format || 'words',
        from_name: '',
        recipients: '',
        domain: props.settings.domain || envDomains[0],
        catch_all: props.catchAll ? 'true' : 'false',
        catch_all_suffix: ''
    }
}

const copyAlias = (alias: string) => {
    navigator.clipboard.writeText(alias)
}

onMounted(() => {
    overlay.autoInit()
    select.autoInit()
    tooltip.autoInit()
    accordion.autoInit()
    addEvents()
    resetAlias()
})
</script>
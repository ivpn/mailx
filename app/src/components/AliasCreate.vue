<template>
    <div>
        <div v-bind:id="'modal-create-alias-' + props.catchAll" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="prev" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4 class="uppercase">
                            {{ props.label }}
                            <span v-if="step == 2">| {{ alias.format }}</span>
                        </h4>
                        <span v-if="!props.catchAll">{{ step }}/2</span>
                    </header>
                    <article>
                        <div v-if="!props.catchAll" v-bind:class="{ 'hidden': step !== 1 }">
                            <h5>Format</h5>
                            <div class="select">
                                <button
                                    v-bind:class="{ 'active': alias.format === 'words' }"
                                    @click="selectFormat('words')">
                                    <div>
                                        <strong>Words</strong>
                                        <span>e.g.: <i>quiet.haze16</i></span>
                                    </div>
                                    <i class="icon arrow-left-line icon-primary" @click.stop="next"></i>
                                </button>
                                <button
                                    v-bind:class="{ 'active': alias.format === 'random' }"
                                    @click="selectFormat('random')">
                                    <div>
                                        <strong>Random</strong>
                                        <span>e.g.: <i>uf1h0hxi</i></span>
                                    </div>
                                    <i class="icon arrow-left-line icon-primary" @click.stop="next"></i>
                                </button>
                                <button
                                    v-bind:class="{ 'active': alias.format === 'uuid' }"
                                    @click="selectFormat('uuid')">
                                    <div>
                                        <strong>UUID</strong>
                                        <span>e.g.: <i>550e8400-e29b-41d4-a716-446655440000</i></span>
                                    </div>
                                    <i class="icon arrow-left-line icon-primary" @click.stop="next"></i>
                                </button>
                            </div>
                        </div>

                        <div v-if="props.catchAll">
                            <div class="mb-5">
                                <label for="alias_catch_all_suffix">
                                    Alias sufix (6-12 alphanumeric characters):
                                </label>
                                <input 
                                    v-model="alias.catch_all_suffix"
                                    v-bind:class="{ 'error': errorCatchAllSuffix }"
                                    id="alias_catch_all_suffix"
                                    type="text"
                                >
                                <p v-if="errorCatchAllSuffix" class="error">Catch-all suffix must be between 6 and 12 characters</p>
                                <p class="text-primary mb-1">
                                    *+{{ alias.catch_all_suffix }}@{{ alias.domain }}
                                </p>
                            </div>
                        </div>

                        <div v-bind:class="{ 'hidden': step !== 2 && !props.catchAll }">
                            <div class="mb-5">
                                <label for="alias_description">
                                    Description
                                </label>
                                <input
                                    v-model="alias.description"
                                    id="alias_description"
                                    type="text"
                                >
                            </div>
                            <div class="mb-5">
                                <label for="alias_from_name">
                                    From name
                                </label>
                                <input
                                    v-model="alias.from_name"
                                    id="alias_from_name"
                                    type="text"
                                >
                            </div>
                            <div class="mb-6">
                                <label for="create-alias-recipient" class="required">
                                    Recipients
                                </label>
                                <select
                                    id="create-alias-recipient"
                                    v-model="selectRecipients"
                                    :disabled="!recipients.length"
                                    :multiple="true"
                                    data-hs-select='{
                                    "placeholder": "Select recipient",
                                    "toggleTag": "<button type=\"button\" aria-expanded=\"false\"></button>",
                                    "toggleClasses": "hs-select-disabled:pointer-events-none hs-select-disabled:opacity-50 relative py-2.5 ps-4 pe-9 flex gap-x-2 text-nowrap w-full cursor-pointer border border-primary text-secondary bg-secondary leading-tight focus:border-accent",
                                    "dropdownClasses": "mt-2 z-50 w-full max-h-72 p-1 space-y-0.5 bg-primary border border-primary overflow-hidden overflow-y-auto",
                                    "optionClasses": "py-2 px-4 w-full text-secondary cursor-pointer hover:bg-secondary bg-primary",
                                    "optionTemplate": "<div class=\"flex justify-between items-center w-full\"><span data-title></span><span class=\"hidden hs-selected:block\"><svg class=\"shrink-0 size-3.5 text-accent \" xmlns=\"http:.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><polyline points=\"20 6 9 17 4 12\"/></svg></span></div>",
                                    "extraMarkup": "<div class=\"absolute top-1/2 end-3 -translate-y-1/2\"><svg class=\"shrink-0 size-3.5 text-secondary \" xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"m7 15 5 5 5-5\"/><path d=\"m7 9 5-5 5 5\"/></svg></div>"
                                    }' class="hidden">
                                    <option v-for="(recipient, _) in recipients"
                                        v-bind:value=recipient
                                        :selected="recipient == settings.recipient"
                                        :key="recipient">
                                        {{ recipient }}
                                    </option>
                                </select>
                                <p v-if="errorRecipients" class="error pt-3">{{ errorRecipients }}</p>
                            </div>
                            <div class="mb-6">
                                <label for="alias_domain" class="required">
                                    Domain
                                </label>
                                <select id="alias_domain" :disabled="!domains.length">
                                    <option v-for="(domain, index) in domains" v-bind:domain
                                        :selected="domain == alias.domain || index === 0" :key="domain">
                                        {{ domain }}
                                    </option>
                                </select>
                            </div>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button @click="close" class="cancel">
                                Cancel
                            </button>
                            <button
                                v-if="step == 2 || props.catchAll"
                                v-bind:disabled="errorRecipients.length > 0"
                                @click="postAlias"
                                class="cta">
                                Create and copy to clipboard
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
const error = ref('')
const errorRecipients = ref('')
const errorCatchAllSuffix = ref(false)
const step = ref(1)

const postAlias = async () => {
    if (!validate()) return

    alias.value.domain = (document.getElementById('alias_domain') as HTMLInputElement).value
    alias.value.recipients = selectRecipients.value.join(',')
    alias.value.enabled = true

    if (props.catchAll) {
        alias.value.format = 'catch_all'
    }

    try {
        const res = await aliasApi.create(alias.value)
        copyAlias(res.data.name)
        events.emit('alias.create', {})
        error.value = ''
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        }
    }
}

const close = () => {
    resetAlias()
    error.value = ''
    const modal = document.querySelector('#modal-create-alias-' + props.catchAll) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-create-alias-' + props.catchAll as any, true) as any
    modal.element.on('close', () => {
        close()
    })

    const multiselect = select.getInstance('#create-alias-recipient' as any, true) as any
    multiselect.element.on('change', (val: any) => {
        errorRecipients.value = val.length === 0 ? 'Select one or more recipients' : ''
    })
}

const validate = () => {
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
        domain: envDomains[0],
        catch_all: props.catchAll ? 'true' : 'false',
        catch_all_suffix: ''
    }
    step.value = 1
}

const selectFormat = (format: string) => {
    alias.value.format = format
}

const next = () => {
    step.value = 2
}

const prev = () => {
    if (step.value === 2) {
        step.value = 1
        return
    }

    close()
}

const copyAlias = (alias: string) => {
    navigator.clipboard.writeText(alias)
}

onMounted(() => {
    overlay.autoInit()
    select.autoInit()
    addEvents()
    resetAlias()
})
</script>
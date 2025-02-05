<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-create-alias' + props.catchAll"
            class="mt-3 py-2 pl-2 pr-3 inline-flex justify-center items-center gap-x-2 font-medium text-base bg-bluish-500 text-white hover:bg-bluish-600">
            <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"
                stroke-linejoin="round">
                <path d="M5 12h14"></path>
                <path d="M12 5v14"></path>
            </svg>
            Create Alias
        </button>
        <div v-bind:id="'modal-create-alias' + props.catchAll"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white dark:bg-neutral-800 border dark:border-neutral-600 shadow-sm rounded pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b dark:border-neutral-600">
                        <h3 class="text-xl text-gray-800 dark:text-gray-100 font-semibold">
                            <span v-if="!props.catchAll">Create Alias</span>
                            <span v-if="props.catchAll">Create Catch-all Alias</span>
                        </h3>
                        <button @click="close" type="button"
                            class="flex justify-center items-center size-7 text-sm font-semibold rounded-full border border-transparent text-gray-800 dark:text-gray-100 hover:bg-gray-100 dark:hover:bg-neutral-700 disabled:opacity-50 disabled:pointer-events-none">
                            <span class="sr-only">Close</span>
                            <svg class="flex-shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                stroke-linecap="round" stroke-linejoin="round">
                                <path d="M18 6 6 18"></path>
                                <path d="m6 6 12 12"></path>
                            </svg>
                        </button>
                    </div>
                    <div class="p-4 whitespace-normal text-left text-base">
                        <div v-if="!props.catchAll" class="grid space-y-3 mb-5">
                            <p class="text-gray-500 dark:text-gray-400 mb-1">
                                Alias format:
                            </p>
                            <div class="relative flex items-start">
                                <div class="flex items-center h-5 mt-1">
                                    <input v-model="alias.format" value="word" id="hs-radio-word" name="hs-radio-with-description" type="radio"
                                        class="form-radio border-gray-200 rounded-full text-bluish-500 focus:ring-white dark:focus:ring-transparent dark:bg-neutral-800 dark:border-neutral-600 dark:checked:border-transparent dark:focus:ring-offset-gray-800"
                                        aria-describedby="hs-radio-word-description" checked>
                                </div>
                                <label for="hs-radio-word" class="ms-3">
                                    <span class="block text-sm font-semibold text-gray-800 dark:text-gray-100">Word</span>
                                    <span id="hs-radio-word-description" class="block text-sm text-gray-500 dark:text-gray-400">e.g.
                                        quiet.haze16@{{ alias.domain }}</span>
                                </label>
                            </div>                            
                            <div class="relative flex items-start">
                                <div class="flex items-center h-5 mt-1">
                                    <input v-model="alias.format" value="chars" id="hs-radio-chars" name="hs-radio-with-description" type="radio"
                                        class="form-radio border-gray-200 rounded-full text-bluish-500 focus:ring-white dark:focus:ring-transparent dark:bg-neutral-800 dark:border-neutral-600 dark:checked:border-transparent dark:focus:ring-offset-gray-800"
                                        aria-describedby="hs-radio-chars-description">
                                </div>
                                <label for="hs-radio-chars" class="ms-3">
                                    <span class="block text-sm font-semibold text-gray-800 dark:text-gray-100">Random</span>
                                    <span id="hs-radio-chars-description" class="block text-sm text-gray-500 dark:text-gray-400">e.g.
                                        uf1h0hxi@{{ alias.domain }}</span>
                                </label>
                            </div>
                            <div class="relative flex items-start">
                                <div class="flex items-center h-5 mt-1">
                                    <input v-model="alias.format" value="uuid" id="hs-radio-uuid" name="hs-radio-with-description" type="radio"
                                        class="form-radio border-gray-200 rounded-full text-bluish-500 focus:ring-white dark:focus:ring-transparent dark:bg-neutral-800 dark:border-neutral-600 dark:checked:border-transparent dark:focus:ring-offset-gray-800"
                                        aria-describedby="hs-radio-uuid-description">
                                </div>
                                <label for="hs-radio-uuid" class="ms-3">
                                    <span class="block text-sm font-semibold text-gray-800 dark:text-gray-100">UUID</span>
                                    <span id="hs-radio-uuid-description" class="block text-sm text-gray-500 dark:text-gray-400">e.g.
                                        550e8400-e29b-41d4-a716-446655440000@{{ alias.domain }}</span>
                                </label>
                            </div>
                        </div>
                        <div v-if="props.catchAll">
                            <div class="mb-5">
                                <label for="alias_catch_all_suffix" class="block text-gray-500 dark:text-gray-400 mb-3">
                                    Alias sufix (6-12 alphanumeric characters):
                                </label>
                                <input id="alias_catch_all_suffix" v-model="alias.catch_all_suffix"
                                    v-bind:class="{ 'border-gray-500': !errorCatchAllSuffix, 'border-red-600 dark:border-red-600': errorCatchAllSuffix }"
                                    class="appearance-none outline-none border border-gray-500 w-full py-3 px-4 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 leading-tight focus:border-bluish-500 mb-3"
                                    type="text">
                                <p v-if="errorCatchAllSuffix" class="text-red-600 text-sm mb-3">Catch-all suffix must be between 6 and 12 characters</p>
                                <p class="text-white dark:text-gray-100 mb-1">
                                    *+{{ alias.catch_all_suffix }}@{{ alias.domain }}
                                </p>
                            </div>
                        </div>
                        <div class="mb-5">
                            <label for="alias_description" class="block text-gray-500 dark:text-gray-400 mb-3">
                                Description:
                            </label>
                            <input id="alias_description" v-model="alias.description"
                                class="appearance-none outline-none border border-gray-500 w-full py-3 px-4 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 leading-tight focus:border-bluish-500 mb-2"
                                type="text">
                        </div>
                        <div class="mb-5">
                            <label for="alias_from_name" class="block text-gray-500 dark:text-gray-400 mb-3">
                                From name:
                            </label>
                            <input id="alias_from_name" v-model="alias.from_name"
                                class="appearance-none outline-none border border-gray-500 w-full py-3 px-4 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 leading-tight focus:border-bluish-500 mb-2"
                                type="text">
                        </div>
                        <div class="mb-6">
                            <label
                                for="create-alias-recipient"
                                class="block text-gray-500 dark:text-gray-400 mb-3">
                                Recipients:
                            </label>
                            <select
                                id="create-alias-recipient"
                                v-model="selectRecipients"
                                :disabled="!recipients.length"
                                :multiple="true"
                                data-hs-select='{
                                "placeholder": "Select recipient",
                                "toggleTag": "<button type=\"button\" aria-expanded=\"false\"></button>",
                                "toggleClasses": "hs-select-disabled:pointer-events-none hs-select-disabled:opacity-50 relative py-3 ps-4 pe-9 flex gap-x-2 text-nowrap w-full cursor-pointer border border-gray-500 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 leading-tight focus:border-bluish-500",
                                "dropdownClasses": "mt-2 z-50 w-full max-h-72 p-1 space-y-0.5 bg-white border border-gray-200 overflow-hidden overflow-y-auto [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-gray-100 [&::-webkit-scrollbar-thumb]:bg-gray-300 dark:[&::-webkit-scrollbar-track]:bg-neutral-700 dark:[&::-webkit-scrollbar-thumb]:bg-neutral-500 dark:bg-neutral-900 dark:border-neutral-700",
                                "optionClasses": "py-2 px-4 w-full text-gray-800 cursor-pointer hover:bg-gray-100 focus:outline-none focus:bg-gray-100 dark:bg-neutral-900 dark:hover:bg-neutral-800 dark:text-neutral-200 dark:focus:bg-neutral-800",
                                "optionTemplate": "<div class=\"flex justify-between items-center w-full\"><span data-title></span><span class=\"hidden hs-selected:block\"><svg class=\"shrink-0 size-3.5 text-bluish-600 dark:text-bluish-500 \" xmlns=\"http:.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><polyline points=\"20 6 9 17 4 12\"/></svg></span></div>",
                                "extraMarkup": "<div class=\"absolute top-1/2 end-3 -translate-y-1/2\"><svg class=\"shrink-0 size-3.5 text-gray-500 dark:text-neutral-500 \" xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"m7 15 5 5 5-5\"/><path d=\"m7 9 5-5 5 5\"/></svg></div>"
                                }' class="hidden">
                                <option v-for="(recipient, _) in recipients"
                                    v-bind:value=recipient
                                    :selected="recipient == settings.recipient"
                                    :key="recipient">
                                    {{ recipient }}
                                </option>
                            </select>
                            <p v-if="errorRecipients" class="pt-3 text-red-600 text-sm">{{ errorRecipients }}</p>
                        </div>

                        <div class="mb-6">
                            <label class="block text-gray-500 dark:text-gray-400 mb-3" for="alias_domain">
                                Domain:
                            </label>
                            <select id="alias_domain" :disabled="!domains.length"
                                class="form-select py-2.5 px-4 pe-9 block w-full border border-gray-500 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 focus:border-bluish-500 disabled:opacity-50 disabled:pointer-events-none outline-none focus:ring-transparent">
                                <option v-for="(domain, index) in domains" v-bind:domain
                                    :selected="domain == alias.domain || index === 0" :key="domain">
                                    {{ domain }}
                                </option>
                            </select>
                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-3 py-4 px-4 border-t dark:border-neutral-600">
                        <button
                            v-bind:disabled="errorRecipients.length > 0"
                            @click="postAlias"
                            class="py-2 px-3 inline-flex items-center gap-x-2 font-medium text-base bg-bluish-500 text-white hover:bg-bluish-600 disabled:opacity-50 disabled:pointer-events-none">
                            Create Alias
                        </button>
                        <button @click="close"
                            class="text-gray-500 bg-gray-100 hover:bg-gray-200 dark:text-gray-300 dark:bg-neutral-600 dark:hover:bg-neutral-700 font-medium text-base py-2 px-3 focus:outline-none focus:shadow-outline">
                            Cancel
                        </button>
                    </div>
                    <div class="flex items-start">
                        <p v-if="error" class="px-5 text-red-600 text-sm mb-3">Error: {{ error }}</p>
                    </div>
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
const props = defineProps(['recipients', 'settings', 'catchAll'])
const alias = ref({
    description: '',
    enabled: true,
    format: 'word',
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

const postAlias = async () => {
    if (!validate()) return

    const domainInput = document.getElementById('alias_domain') as HTMLInputElement
    alias.value.domain = domainInput.value
    alias.value.recipients = selectRecipients.value.join(',')
    alias.value.enabled = true

    if (props.catchAll) {
        alias.value.format = 'catch_all'
    }

    try {
        await aliasApi.create(alias.value)
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
    alias.value = { format: 'word' } as any
    error.value = ''
    const modal = document.querySelector('#modal-create-alias' + props.catchAll) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-create-alias' + props.catchAll as any, true) as any
    modal.element.on('close', () => {
        close()
    })

    const multiselect = select.getInstance('#create-alias-recipient' as any, true) as any
    multiselect.element.on('change', (val: any) => {
        errorRecipients.value = val.length === 0 ? 'Select one or more recipients' : ''
    })
}

const validate = () => {
    if (props.catchAll && (alias.value.catch_all_suffix.length < 6 || alias.value.catch_all_suffix.length > 12)) {
        errorCatchAllSuffix.value = true
    } else {
        errorCatchAllSuffix.value = false
    }

    return !errorCatchAllSuffix.value
}

onMounted(() => {
    overlay.autoInit()
    select.autoInit()
    addEvents()
})
</script>
<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-alias-edit' + alias.id"
            class="text-bluish-500 hover:text-bluish-600 font-medium text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
            type="submit">
            Edit
        </button>
        <div v-bind:id="'modal-alias-edit' + alias.id"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white dark:bg-neutral-800 border dark:border-neutral-600 shadow-sm rounded pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b dark:border-neutral-600">
                        <h3 class="text-xl text-gray-800 dark:text-gray-100 font-semibold">
                            Edit alias
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
                        <h1 class="text-xl text-gray-800 dark:text-gray-100 font-semibold mb-5">{{ alias.name }}</h1>
                        <div class="mb-5">
                            <label v-bind:for="'description_' + alias.id"
                                class="block text-gray-500 dark:text-gray-400 mb-3">
                                Description:
                            </label>
                            <input v-bind:id="'description_' + alias.id" v-model="alias.description"
                                class="appearance-none outline-none border border-gray-500 w-full py-3 px-4 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 leading-tight focus:border-bluish-500 mb-2"
                                type="text">
                        </div>
                        <div class="mb-5">
                            <label v-bind:for="'from_' + alias.id"
                                class="block text-gray-500 dark:text-gray-400 mb-3">
                                From name:
                            </label>
                            <input v-bind:id="'from_' + alias.id" v-model="alias.from_name"
                                class="appearance-none outline-none w-full py-3 px-4 mb-2 border border-gray-500 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 dark:border-neutral-400 leading-tight focus:border-bluish-500"
                                type="text">
                        </div>
                        <div class="mb-6">
                            <label v-bind:for="'recipient_' + alias.id"
                                class="block text-gray-500 dark:text-gray-400 mb-3">
                                Recipients:
                            </label>
                            <select
                                v-model="selectRecipients"
                                v-bind:id="'recipient_' + alias.id"
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
                                <option v-for="recipient in recipients"
                                    v-bind:value=recipient
                                    :selected="alias.recipients.includes(recipient)"
                                    :key="recipient">
                                    {{ recipient }}
                                </option>
                            </select>
                            <p v-if="errorRecipients" class="pt-3 text-red-600 text-sm">{{ errorRecipients }}</p>
                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-3 py-4 px-4 border-t dark:border-neutral-600">
                        <button
                            v-if="!success"
                            @click="updateAlias"
                            v-bind:disabled="errorRecipients.length > 0"
                            class="py-2 px-3 inline-flex items-center gap-x-2 font-medium text-base bg-bluish-500 text-white hover:bg-bluish-600 disabled:opacity-50 disabled:pointer-events-none">
                            Save
                        </button>
                        <button
                            @click="close"
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

const props = defineProps(['alias', 'recipients'])
let alias = ref(Object.assign({}, props.alias))
const recipients = ref(props.recipients)
const selectRecipients = ref(props.alias.recipients)
const success = ref('')
const error = ref('')
const errorRecipients = ref('')

const updateAlias = async () => {
    alias.value.recipients = selectRecipients.value.toString()

    try {
        const res = await aliasApi.update(alias.value.id, alias.value)
        success.value = res.data.message
        error.value = ''
        events.emit('alias.update', {})
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            success.value = ''
            error.value = err.message
        }
    }
}

const close = () => {
    alias.value.description = props.alias.description
    alias.value.from_name = props.alias.from_name
    alias.value.recipients = props.alias.recipients
    success.value = ''
    error.value = ''
    const modal = document.querySelector('#modal-alias-edit' + alias.value.id) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-alias-edit' + alias.value.id as any, true) as any
    modal.element.on('close', () => {
        close()
    })

    const multiselect = select.getInstance('#recipient_' + alias.value.id as any, true) as any
    multiselect.element.on('change', (val: any) => {
        errorRecipients.value = val.length === 0 ? 'Select one or more recipients' : ''
    })
}

onMounted(() => {
    overlay.autoInit()
    select.autoInit()
    addEvents()
})
</script>
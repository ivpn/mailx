<template>
    <div>
        <button data-hs-overlay="#hs-basic-modal"
            class="text-violet-600 hover:text-violet-700 font-semibold text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
            type="submit">
            Edit
        </button>
        <div id="hs-basic-modal"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white border shadow-sm rounded-xl pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b">
                        <h3 class="font-bold text-gray-800">
                            Edit alias
                        </h3>
                        <button type="button"
                            class="flex justify-center items-center size-7 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none"
                            data-hs-overlay="#hs-basic-modal">
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
                        <h1 class="text-xl font-bold text-gray-800 mb-5">{{ alias.name }}</h1>
                        <div class="max-w-xs mb-5">
                            <label
                                v-bind:for="'description_' + alias.id"
                                class="block text-gray-500 text-sm font-semibold mb-3">
                                Description:
                            </label>
                            <input
                                v-bind:id="'description_' + alias.id"
                                v-model="alias.description"
                                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-violet-600 mb-2" type="text">
                        </div>
                        <div class="max-w-xs mb-5">
                            <label
                                v-bind:for="'from_' + alias.id"
                                class="block text-gray-500 text-sm font-semibold mb-3">
                                From name:
                            </label>
                            <input
                                v-bind:id="'from_' + alias.id"
                                v-model="alias.from_name"
                                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-violet-600 mb-2" type="text">
                        </div>
                        <div class="max-w-xs mb-6">
                            <label
                                v-bind:for="'recipient_' + alias.id"
                                class="block text-gray-500 text-sm font-semibold mb-3">
                                Select default recipient:
                            </label>
                            <select
                                v-bind:id="'recipient_' + alias.id"
                                :disabled="!recipients.length"
                                class="form-select py-2.5 px-4 pe-9 block w-full border-2 border-gray-200 rounded-lg text-gray-700 focus:border-violet-600 disabled:opacity-50 disabled:pointer-events-none outline-none">
                                <option v-for="recipient in recipients" v-bind:value=recipient
                                    :selected="recipient == alias.recipients" :key="recipient">
                                    {{ recipient }}
                                </option>
                            </select>
                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-2 py-3 px-4 border-t">
                        <button
                            @click="updateAlias"
                            class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-semibold rounded-md border border-transparent bg-violet-600 text-white hover:bg-violet-700 disabled:opacity-50 disabled:pointer-events-none">
                            Save
                        </button>
                        <button type="button"
                            class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-medium rounded-md border border-gray-200 bg-white text-gray-800 shadow-sm hover:bg-gray-50 disabled:opacity-50 disabled:pointer-events-none"
                            data-hs-overlay="#hs-basic-modal">
                            Close
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import overlay from '@preline/overlay'
import axios from 'axios'
import { aliasApi } from '../api/alias.ts'

const props = defineProps(['alias', 'recipients'])
const alias = ref(props.alias)
const recipients = ref(props.recipients)
const error = ref('')

const updateAlias = async () => {
    try {
        await aliasApi.update(alias.value.id, alias.value)
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

onMounted(() => {
    overlay.autoInit()
})
</script>
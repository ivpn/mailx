<template>
    <div>
        <hr class="mb-7 dark:border-gray-600">
        <div class="flex flex-row justify-between pb-3">
            <div>
                <select
                v-model="limit"
                @change="updateLimit"
                class="form-select py-2 px-4 pe-9 block w-full border border-gray-200 dark:border-neutral-600 bg-transparent text-gray-500 dark:text-gray-300 focus:border-bluish-500 disabled:opacity-50 disabled:pointer-events-none outline-none focus:ring-transparent">
                <option>25</option>
                <option>50</option>
                <option>75</option>
            </select>
            </div>
            <nav class="flex items-center gap-x-1">
                <button type="button"
                    @click="prev"
                    class="min-h-[38px] min-w-[38px] py-2 px-2.5 inline-flex justify-center items-center gap-x-2 text-sm text-gray-800 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-neutral-700 focus:outline-none disabled:opacity-50 disabled:pointer-events-none">
                    <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                        stroke-linejoin="round">
                        <path d="m15 18-6-6 6-6"></path>
                    </svg>
                    <span aria-hidden="true" class="sr-only">Previous</span>
                </button>
                <div class="flex items-center gap-x-1">
                    <span
                        class="min-h-[38px] min-w-[38px] flex justify-center items-center border border-gray-200 dark:border-neutral-600 text-gray-800 dark:text-gray-400 py-2 px-3 text-sm focus:outline-none disabled:opacity-50 disabled:pointer-events-none">{{ page }}</span>
                    <span
                        class="min-h-[38px] flex justify-center items-center text-gray-500 dark:text-gray-400 py-2 px-1.5 text-sm">of</span>
                    <span
                        class="min-h-[38px] flex justify-center items-center text-gray-500 dark:text-gray-400 py-2 px-1.5 text-sm">{{ pages }}</span>
                </div>
                <button type="button"
                    @click="next"
                    class="min-h-[38px] min-w-[38px] py-2 px-2.5 inline-flex justify-center items-center gap-x-2 text-sm text-gray-800 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-neutral-700 focus:outline-none disabled:opacity-50 disabled:pointer-events-none">
                    <span aria-hidden="true" class="sr-only">Next</span>
                    <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                        stroke-linejoin="round">
                        <path d="m9 18 6-6-6-6"></path>
                    </svg>
                </button>
            </nav>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

const emit = defineEmits(['onUpdatePage'])
const props = defineProps(['limit', 'page', 'total'])
const limit = ref(props.limit)
const page = ref(props.page)
const total = ref(props.total)
const pages = ref(1)

const next = () => {
    if (page.value * props.limit >= total.value) return
    page.value++
    emit('onUpdatePage', { limit: limit.value, page: page.value })
}

const prev = () => {
    if (page.value === 1) return
    page.value--
    emit('onUpdatePage', { limit: limit.value, page: page.value })
}

const updateLimit = () => {
    page.value = 1
    updatePages()
    emit('onUpdatePage', { limit: limit.value, page: page.value })
}

const updatePages = () => {
    pages.value = Math.ceil(total.value / limit.value)
}   

onMounted(() => {
    updatePages()
})
</script>
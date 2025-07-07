<template>
    <div class="flex flex-row justify-between pt-5 pb-3">
        <div class="flex items-center gap-x-3">
            <select v-model="limit" @change="updateLimit" class="py-2 mb-0">
                <option>25</option>
                <option>50</option>
                <option>75</option>
            </select>
            <p class="text-nowrap m-0 desktop">per page</p>
        </div>
        <nav class="flex items-center gap-x-2">
            <div class="flex items-center gap-x-3">
                <span class="min-h-[38px] min-w-[38px] inline-flex justify-center items-center border border-secondary text-secondary text-sm">{{ page }}</span>
                <span class="min-h-[38px] min-w-[38px] inline-flex justify-center items-center text-secondary text-sm">of {{ pages }}</span>
            </div>
            <button
                @click="prev"
                class="h-[38px] w-[38px] inline-flex justify-center items-center hover:bg-secondary">
                <i class="icon arrow-down rotate-90 text-xl icon-secondary"></i>
            </button>
            <button
                @click="next"
                class="h-[38px] w-[38px] inline-flex justify-center items-center hover:bg-secondary">
                <i class="icon arrow-down -rotate-90 text-xl icon-secondary"></i>
            </button>
        </nav>
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
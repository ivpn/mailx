<template>
    <div class="grid grid-cols-6 gap-2 text-sm border-b border-gray-200 items-center py-3">
        <div>{{ new Date(alias.created).toDateString() }}</div>
        <div>{{ alias.name }}</div>
        <div>{{ alias.recipients.split(',').length }}</div>
        <div>
            {{ alias.stats.forwards }}/{{ alias.stats.blocks }}/{{ alias.stats.replies }}/{{ alias.stats.sends }}
        </div>
        <div>{{ alias.enabled }}</div>
        <div class="flex gap-2 justify-items-end place-content-end">
            <button
                class="bg-violet-600 hover:bg-violet-700 text-white font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                type="submit">
                Edit
            </button>
            <button
                class="bg-red-600 hover:bg-red-700 text-white font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                type="submit">
                Delete
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps(['alias'])
</script>
<template>
    <tr>
        <td class="pr-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            {{ new Date(alias.created).toDateString() }}
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            {{ alias.name }}
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            {{ alias.recipients.split(',').length }}
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm">
            {{ alias.stats.forwards }}/{{ alias.stats.blocks }}/{{ alias.stats.replies }}/{{ alias.stats.sends }}
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm">
            {{ alias.enabled }}
        </td>
        <td class="pl-5 py-4 whitespace-nowrap text-end text-sm">
            <div class="flex gap-5 justify-end">
                <button
                    class="text-violet-600 hover:text-violet-700 font-semibold text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
                    type="submit">
                    Edit
                </button>
                <button
                    class="text-red-600 hover:text-red-700 font-semibold text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
                    type="submit">
                    Delete
                </button>
            </div>
        </td>
    </tr>
</template>

<script setup lang="ts">
defineProps(['alias'])
</script>
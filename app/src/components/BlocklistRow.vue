<template>
    <tr v-if="isDesktop" class="desktop-lg">
        <td>
            <p>{{ new Date(item.created_at).toDateString() }}</p>
        </td>
        <td>
            <p class="text-wrap break-all">{{ item.value }}</p>
        </td>
        <td>
            <p>{{ typeLabel }}</p>
        </td>
        <td>
            <p>{{ item.blocked }}</p>
        </td>
        <td>
            <div class="hs-dropdown [--offset:0]">
                <button v-bind:id="'hs-dropdown-blocklist-' + item.id">
                    <i class="icon icon-secondary more text-lg"></i>
                </button>
                <div
                    class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                    v-bind:aria-labelledby="'hs-dropdown-blocklist-' + item.id"
                    >
                    <button @click.stop="deleteItem" class="delete">
                        <i class="icon icon-error trash text-xs"></i>
                        Delete
                    </button>
                </div>
            </div>
        </td>
    </tr>
    <tr v-else class="tablet-lg">
        <td>
            <div class="flex gap-3 justify-between">
                <div class="text-start">
                    <div>
                        <p class="mb-3">{{ new Date(item.created_at).toDateString() }}</p>
                    </div>
                    <div class="break-all">
                        <p class="plain text-base text-primary text-wrap">
                            <span class="text-sm text-secondary">{{ typeLabel }}:</span><br>
                            {{ item.value }}
                        </p>
                    </div>
                    <div>
                        <p class="plain text-base text-primary">
                            <span class="text-sm text-secondary">Blocked:</span><br>
                            {{ item.blocked }}
                        </p>
                    </div>
                </div>
                <div class="text-end">
                    <div class="hs-dropdown">
                        <button class="py-0" v-bind:id="'hs-dropdown-blocklist-' + item.id">
                            <i class="icon icon-secondary more text-lg"></i>
                        </button>
                        <div
                            class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                            v-bind:aria-labelledby="'hs-dropdown-blocklist-' + item.id"
                            >
                            <button @click.stop="deleteItem" class="delete">
                                <i class="icon icon-error trash text-xs"></i>
                                Delete
                            </button>
                        </div>
                    </div>
                </div>
            </div>
            <hr>
        </td>
    </tr>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useBreakpoint } from '../lib/useBreakpoint.ts'
import { closeDropdowns } from '../lib/preline.ts'

const props = defineProps(['item'])
const emit = defineEmits(['onDelete'])
const { isDesktop } = useBreakpoint()

const typeLabel = computed(() => props.item.type === 'domain' ? 'Domain' : 'Email')

const deleteItem = () => {
    closeDropdowns()
    if (!confirm('Are you sure you want to remove ' + props.item.value + ' from your blocklist?')) return

    emit('onDelete', props.item)
}
</script>

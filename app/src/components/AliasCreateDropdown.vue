<template>
    <div class="hs-dropdown">
        <button id="hs-dropdown-new-alias" class="cta">
            New Alias
            <i class="icon arrow-down text-xl icon-white"></i>
        </button>
        <div class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden" aria-labelledby="hs-dropdown-new-alias">
            <button data-hs-overlay="#modal-create-alias-false">
                <i class="icon at icon-accent text-sm"></i>
                Alias
            </button>
            <button data-hs-overlay="#modal-create-alias-true">
                <i class="icon scan icon-accent text-sm"></i>
                Catch-all Alias
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import dropdown from '@preline/dropdown'
import overlay from '@preline/overlay'

onMounted(async () => {
    dropdown.autoInit()
    overlay.autoInit()
})
</script>
<template>
    <div class="hs-dropdown">
        <button id="hs-dropdown-new-alias" class="cta">
            New Alias
            <span class="icon arrow-down-fill"></span>
        </button>
        <div class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden" aria-labelledby="hs-dropdown-new-alias">
            <button data-hs-overlay="#modal-create-alias-false">
                <span class="icon at bg-bd-9 text-sm"></span>
                Alias
            </button>
            <button data-hs-overlay="#modal-create-alias-true">
                <span class="icon scan bg-bd-9 text-sm"></span>
                Catch-all Alias
            </button>
        </div>
    </div>
</template>
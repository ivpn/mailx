<template>
    <div class="mb-5">
        <h2>Alias Import</h2>
        <p>
            Import a list of your aliases from a CSV file. Only aliases with your verified domain will be imported.
        </p>
        <p>
            CSV file format:
            <br>
            <code class="text-sm">
                alias,description,enabled,recipients<br>
                some.alias@example.net,A description,true,recipient1@provider.net recipient2@provider.net
            </code>
        </p>
        <div class="mb-6">
            <input type="file" name="csv" id="csv" accept=".csv">
            <p class="text-sm">Only .csv files are supported.</p>
        </div>
        <div>
            <button
                v-if="!importing"
                @click="importAliases"
                class="cta mb-4">
                Import Aliases
            </button>
            <button
                v-if="importing"
                disabled
                class="cta mb-4">
                Importing...
            </button>
        </div>
        <p v-if="error" class="error">Error: {{ error }}</p>
        <p v-if="success.count > 0" class="success">Successfully imported {{ success.count }} aliases.</p>
        <p v-if="success.count === 0 && success.message" class="success">0 aliases imported.</p>
    </div>
</template>

<script setup lang="ts">
import { aliasApi } from '../api/alias'
import axios from 'axios'
import { ref } from 'vue'

const importing = ref(false)
const error = ref('')
const success = ref({
    count: 0,
    message: '',
})

const validateFile = (file: File | null): string | null => {
    if (!file) {
        return 'Please select a CSV file to import.'
    }
    if (file.type !== 'text/csv' && !file.name.endsWith('.csv')) {
        return 'Invalid file type. Please select a CSV file.'
    }
    return null
}

const importAliases = async () => {
    importing.value = true
    try {
        const fileInput = document.getElementById('csv') as HTMLInputElement
        if (!fileInput.files || fileInput.files.length === 0) {
            error.value = 'Please select a CSV file to import.'
            return
        }

        const validationError = validateFile(fileInput.files[0])
        if (validationError) {
            error.value = validationError
            return
        }

        const formData = new FormData()
        formData.append('file', fileInput.files[0])
        const res = await aliasApi.import(formData)
        success.value.count = res.data.count
        success.value.message = res.data.message
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    } finally {
        importing.value = false
    }
}
</script>
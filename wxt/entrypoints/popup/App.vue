<template>
  <div class="max-w-[420px] min-w-[420px]">
    <Login v-if="!apiToken" />
    <Aliases v-if="apiToken && defaults && route === 'aliases'" :apiToken="apiToken" :defaults="defaults" />
    <Settings v-if="apiToken && defaults && route === 'settings'" :apiToken="apiToken" :defaults="defaults" />
    <header v-if="apiToken && defaults" class="bg-secondary fixed bottom-0 left-0 right-0 z-10">
      <nav>
          <div class="flex flex-row items-center">
              <button @click="updateRoute('aliases')" v-bind:class="{ 'active': route === 'aliases' }">
                  Aliases
              </button>
              <button @click="updateRoute('settings')" v-bind:class="{ 'active': route === 'settings' }">
                  Settings
              </button>
          </div>
      </nav>
    </header>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { api } from '@/lib/api'
import { store } from '@/lib/store'
import { Defaults } from '@/lib/types'
import Login from '@/components/popup/Login.vue'
import Aliases from '@/components/popup/Aliases.vue'
import Settings from '@/components/popup/Settings.vue'

const apiToken = ref<string | undefined>()
const defaults = ref<Defaults | undefined>()
const route = ref<string>('aliases')

const updateRoute = (newRoute: string) => {
    route.value = newRoute
}

onMounted(async () => {
  apiToken.value = await store.getApiToken()
  defaults.value = await store.getDefaults()

  store.onApiTokenChange((newToken) => {
    console.log('API token changed:', newToken)
    apiToken.value = newToken
  })

  store.onDefaultsChange((newDefaults) => {
    console.log('Defaults changed:', newDefaults)
    defaults.value = newDefaults
  })

  try {
    const res = await api.livez()
    console.log('/livez response:', res)
  } catch (err) {
    console.error('Error calling /livez:', err)
  }
})
</script>

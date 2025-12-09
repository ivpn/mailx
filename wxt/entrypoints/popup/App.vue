<template>
  <div class="min-h-[600px] max-w-[420px] min-w-[420px]">
    <Login v-if="!apiToken" />
    <keep-alive>
      <component v-if="apiToken && defaults" :is="activeComponent" :key="route" :apiToken="apiToken"
        :defaults="defaults" />
    </keep-alive>
    <header v-if="apiToken && defaults" class="bg-secondary fixed bottom-0 left-0 right-0 z-10">
      <nav>
        <div class="flex flex-row items-center">
          <button @click="updateRoute('aliases')" :class="{ active: route === 'aliases' }">
            Aliases
          </button>
          <button @click="updateRoute('settings')" :class="{ active: route === 'settings' }">
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
const route = ref('aliases')

const updateRoute = (val: string) => {
  route.value = val
}

const activeComponent = computed(() =>
  route.value === 'aliases' ? Aliases : Settings
)

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

import { readonly, ref } from 'vue'

// Single app-wide listener: the alias table renders one of two layouts per row, and binding a
// matchMedia listener per row would reintroduce the same fan-out this replaces.
// 1024px is Tailwind's `lg`, matching .desktop-lg / .tablet-lg in style/components/responsive.css.
const query = window.matchMedia('(min-width: 1024px)')
const isDesktop = ref(query.matches)

query.addEventListener('change', event => {
    isDesktop.value = event.matches
})

export const useBreakpoint = () => ({ isDesktop: readonly(isDesktop) })

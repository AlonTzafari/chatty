<script setup lang="ts">
import { onBeforeMount } from 'vue'
import { RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth';
const authStore = useAuthStore()
onBeforeMount(async () => {
    const res = await fetch('/api/.me')
    const user = await res.json()
    authStore.user = user
})
</script>

<template>
<RouterView />
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'
import { onMounted, ref } from 'vue'; 
import { useRouter,RouterLink } from 'vue-router';
const authStore = useAuthStore()
const router = useRouter()
const loading = ref(false)
const channels = ref<{Id: string, Name: string}[]>([])

onMounted(async () => {
    try {
        const res = await fetch(`/api/channels?user_id=${authStore.user?.Id}`)
        channels.value = await res.json()
    } catch(e) {
        console.error(e)
    }
})

async function logout() {
    loading.value = true
    try {
        await fetch('/api/logout', {method: 'post'})
        router.push('/login')
    } catch(e) {
        console.error(e)
    }
    loading.value = false
}
</script>

<template>
    <header>
        <Avatar shape="circle" :label="authStore.user?.Username.charAt(0).toUpperCase()" size="large"/>
        <h1>channels</h1>
        <i :class="`pi pi-sign-out clickable ${loading && 'disabled'}`" @click="!loading && logout()"></i>  
    </header>
  <main>
    <RouterLink class="create-channel" to="/create-channel"><Button>+</Button></RouterLink>
    <div v-for="channel of channels" :key="channel.Id">
        {{ channel.Name }}
    </div>
  </main>
</template>
<style scoped>
    header {
        height: 5rem;
        width: 100%;
        border-bottom: 1px solid rgb(37, 35, 35);
        display: flex;
        gap: 2rem;
        align-items: center;
        padding: 1rem 1rem;
    }
    .clickable:hover {
        cursor: pointer;
    }
    .clickable:active {
        cursor: auto;
    }
    .disabled {
        color: gray;
    }
    main {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 1rem;
    }
    .create-channel {
        all: unset;
    }
</style>

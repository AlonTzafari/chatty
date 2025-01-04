<script setup lang="ts">
import { getUserChannels } from '@/api/channel';
import type channelSchema from '@/schemas/channel-schema';
import { useAuthStore } from '@/stores/auth';
import Header from '@/components/app-header.vue'
import Button from 'primevue/button'
import { onMounted, ref } from 'vue'; 
import { RouterLink } from 'vue-router';
import type { z } from 'zod';
import ChannelItem from '@/components/channel-item.vue';

const authStore = useAuthStore()
const channels = ref<z.infer<typeof channelSchema>[]>([])

onMounted(async () => {
    try {
        if(authStore.user) {
            channels.value = await getUserChannels(authStore.user.Id)
        }
    } catch(e) {
        console.error(e)
    }
})
</script>
<template>
    <Header>Channels</Header>
  <main>
    <ChannelItem v-for="channel of channels" :channel="channel" :key="channel.Id" />
</main>
<RouterLink class="unset action" to="/create-channel"><Button icon="pi pi-plus" size="large"></Button></RouterLink>
</template>
<style scoped>
    .disabled {
        color: gray;
    }
    main {
        display: flex;
        flex-direction: column;
        align-items: stretch;
        padding: 1rem;
        gap: 0.25rem;
        height: calc(100% - var(--header-height));
        overflow-y: scroll;
    }
    .action {
        position: absolute;
        bottom: 5rem;
        right: 5rem;
    }
</style>

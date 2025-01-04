<script setup lang="ts">
import type channelSchema from '@/schemas/channel-schema';
import Card from 'primevue/card';
import { RouterLink } from 'vue-router';
import type { z } from 'zod';
    
const { channel } = defineProps<{
        channel: z.infer<typeof channelSchema>
}>()
</script>
<template>
<RouterLink class="unset clickable" :to="`/channel/${channel.Id}`">
        <Card class="channel">
            <template #title>{{ channel.Name }}</template>
        </Card>
</RouterLink>
</template>
<style lang="css" scoped>
    .channel {
        width: 100%;
    }
    .channel:hover {
        filter: brightness(1.1);
    }
</style>
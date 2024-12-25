<script setup lang="ts">
import type messageSchema from '@/schemas/message-schema';
import { formatDate } from '@/utils/format-date';
import Card from 'primevue/card'
import type { z } from 'zod';
    const { message, stickEnd = false } = defineProps<{
        message: z.infer<typeof messageSchema>
        stickEnd?: boolean
    }>()
</script>
<template>
<div class="message-container" :class="{stickEnd}">
    <Card class="card">
        <template #title>
            <div class="message-header">
                <span class="username">{{ message.Username }}</span>
                <span>{{ formatDate(new Date(message.CreatedAt)) }}</span>
            </div>
            <p class="content">{{ message.Content }}</p>
        </template>
    </Card>
</div>
</template>
<style scoped>
    .message-container {
        display: flex;
        width: 100%;
    }
    .stickEnd {
        justify-content: flex-end;
    }
    .message-header {
        display: flex;
        justify-content: space-between;
        gap: 2rem;
        padding: 0.5rem 0;
        font-size: 0.8rem;
        opacity: 0.8;
    }
    .card {
        width: fit-content;
        margin: 0;
        max-width: 90%;
    }
    .username {
        text-overflow: ellipsis;
    }
    .content {
        word-wrap: break-word;
    }
</style>
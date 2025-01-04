import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { z } from 'zod'
import type channelSchema from '@/schemas/channel-schema'

export const useChannelStore = defineStore('channel', () => {
  const channel = ref<z.infer<typeof channelSchema> | null>(null)
  function setChannel(chan: z.infer<typeof channelSchema> | null) {
    channel.value = chan
  }
  return { channel, setChannel }
})

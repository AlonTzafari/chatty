import messageSchema from "@/schemas/message-schema";
import { z } from "zod";
import FetchError from "./fetch-error";

export async function getMessagesByChannel(channelId: string, signal: AbortSignal = AbortSignal.timeout(5000)): Promise<z.infer<typeof messageSchema>[]> {
    const res = await fetch(`/api/messages?channel_id=${channelId}`, {signal})
    if (!res.ok) {
        throw new FetchError(res)
    }
    return z.array(messageSchema).parse(await res.json())
}

export async function sendMessageToChannel(channelId: string, content: string, signal: AbortSignal = AbortSignal.timeout(5000)): Promise<void> {
    const data = {
        channelId,
        content,
    }
    const res = await fetch('/api/messages', {
        method: 'post', 
        body: JSON.stringify(data), 
        headers: {
            "Content-Type": "application/json"
        },
        signal
    })
    if(!res.ok) {
        throw new FetchError(res)
    }
}

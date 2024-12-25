import { z } from "zod";
import FetchError from "./fetch-error";
import channelSchema from "@/schemas/channel-schema";

export async function getChannel(channelId: string, signal: AbortSignal = AbortSignal.timeout(5000)): Promise<z.infer<typeof channelSchema>> {
    const res = await fetch(`/api/channels/${channelId}`, {signal})
    if (!res.ok) {
        throw new FetchError(res)
    }
    return channelSchema.parse(await res.json())
}

export async function getUserChannels(userId: string, signal: AbortSignal = AbortSignal.timeout(5000)): Promise<z.infer<typeof channelSchema>[]> {
    const res = await fetch(`/api/channels?user_id=${userId}`, {signal})
    if (!res.ok) {
        throw new FetchError(res)
    }
    return z.array(channelSchema).parse(await res.json())
}
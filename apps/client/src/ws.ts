import WS from "reconnecting-websocket"
import messageSchema from "./schemas/message-schema"
import channelSchema from "./schemas/channel-schema"


interface Opts {
    url: string
}

interface Message<T=unknown> {
    channel: string,
    payload: T
}

const channelMap = {
    'message-updates': messageSchema.parse.bind(messageSchema),
    'channel-updates': channelSchema.parse.bind(channelSchema)
} as const

type TChannelMap = typeof channelMap

type ChannelPayload<k extends keyof TChannelMap> = ReturnType<TChannelMap[k]>

export class WSClient {
    private ws: WS
    private subCount = new Map<string, number>()
    constructor(private opts: Opts) {
        this.ws = new WS(opts.url, "chat-ws", {startClosed: true})
        
        this.ws.addEventListener("open", () => {
            this.Resubscribe()
        })
    }
    
    public async connect() {
        const ws = this.ws
        await new Promise<void>((resolve) => {
            ws.reconnect()
            ws.addEventListener('open', () => {
                resolve()
            })
        })
    }

    private Resubscribe() {
        for (const [channel, subCount ] of this.subCount.entries()) {
            if(subCount < 1) {
                continue
            }
            this.ws.send(JSON.stringify({
                channel: 'subscribe',
                payload: {
                    topic: channel
                }
            }))
        }
    }

    public close() {
        this.ws.close()
    }
    
    public subscribe<K extends keyof TChannelMap>(channel: K, cb: (data: ChannelPayload<K>) => void): () => void {
        const ws = this.ws
        let count = this.subCount.get(channel)
        if(count == null || count === 0) {
            ws.send(JSON.stringify({
                channel: 'subscribe',
                payload: {
                    topic: channel
                }
            }))
            count = 0
        }
        this.subCount.set(channel, count + 1)
        const handler = (e: MessageEvent) => {
            const message = JSON.parse(e.data) as Message<unknown>
            if (message.channel === channel) {
                const parser = channelMap[channel as keyof TChannelMap]
                try {
                    const payload = parser(message.payload) as ChannelPayload<K>
                    cb(payload)
                } catch(e) {
                    console.error(e)
                }
            }
        }
        ws.addEventListener('message', handler)
        return () => {
            const count = this.subCount.get(channel)
            if(count == null) {
                return
            }
            if(count === 1) {
                ws.send(JSON.stringify({
                    channel: 'quit',
                    payload: {
                        topic: channel
                    }
                }))
                ws.removeEventListener('message', handler)
            }
            if(count >= 1) {
                this.subCount.set(channel, count - 1)
            }
        }
    }
}

export const wsClient = new WSClient({url: '/api/ws'})
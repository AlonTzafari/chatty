import authSchema from "@/schemas/auth-schema"

export async function me(signal: AbortSignal = AbortSignal.timeout(5000)) {
    const res = await fetch('/api/.me', {signal})
    const user = await res.json()
    return authSchema.nullable().parse(user)
}

export async function logout(signal: AbortSignal = AbortSignal.timeout(5000)) {
    await fetch('/api/logout', {method: 'post', signal})
}
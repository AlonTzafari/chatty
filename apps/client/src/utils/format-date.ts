export function formatDate(date: Date): string {
    const d = date.getDate()
    const m = date.getMonth()
    const y = date.getFullYear()
    const h = date.getHours()
    const min = date.getMinutes()

    return `${d}/${m}/${y} ${h}:${min}`
}
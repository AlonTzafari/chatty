export function insertSorted<T>(array: T[], item: T, comparator:(a:T, b:T) => number) {
    const i = sortedIndex(array, item, comparator)
    array.splice(i, 0, item)
    return array
}

export function sortedIndex<T>(array: T[], value: T, comparator:(a:T, b:T) => number) {
    let low = 0
    let high = array.length
    let mid = 0
    while (low < high) {
        mid = (low + high) >>> 1
        if (comparator(array[mid], value) < 0) {
            low = mid + 1
        } else {
            high = mid
        }
    }
    return low;
}
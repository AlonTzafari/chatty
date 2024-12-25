
export default class FetchError extends Error {
    constructor(private response: Response) {
        const errMessage = `${response.status} ${response.statusText}`
        super(errMessage)
        this.response = response
    }
    public async text() {
        return await this.response.text()
    }
}
import 'dotenv/config'
export const PORT=process.env.PORT
export const API_KEY=process.env.COIN_API_KEY
export const COINS_API_URL='https://api.livecoinwatch.com'
export const FIATS='/fiats/all'
export const COINS_SINGLE='/coins/single'
export const COINS_SINGLE_HISTORY='/coins/single/history'
export const COINS_LIST='/coins/list'
export const MAX_LIST_LEN = 15748
export const ASC = 'ascending'
export const DESC = 'descending'
export class Sort {
    static RANK     = new Sort('rank')
    static PRICE    = new Sort('price')
    static VOLUME   = new Sort('volume')
    static CODE     = new Sort('code')
    static NAME     = new Sort('name')
    static AGE      = new Sort('age')

    constructor(name) {
        this.name = name
    }
}

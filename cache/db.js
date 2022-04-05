import low from 'lowdb'
import FileSync from 'lowdb/adapters/FileSync.js'
import fetchPOST from './services/list.fetch.js';
import { ASC, COINS_API_URL, COINS_LIST, FIATS, MAX_LIST_LEN, Sort } from './config.js';
import cron from 'node-cron'

const adapter = new FileSync('db.json')
const db = low(adapter)


// Set some defaults
db.defaults({coins: [], fiats: []})
  .write()

const getData = async () => {
  if (db.get('fiats').value().length === 0) {
    const fiats = await fetchPOST(COINS_API_URL + FIATS, {})
    db.assign({fiats: fiats}).write()
  }
  if (db.get('coins').value().length === 0) {
    const coins = await fetchPOST(COINS_API_URL + COINS_LIST, {
      'currency':'USD', 
      'sort': Sort.RANK.name, 
      'order': ASC, 
      'offset': 0, 
      'limit': MAX_LIST_LEN,
      'meta': true
    })
    db.assign({coins: coins}).write()
  }
}

getData()

cron.schedule('0 * * * *', async () => {
    // Add a post
    const coins = await fetchPOST(COINS_API_URL + COINS_LIST, {
      'currency':'USD', 
      'sort': Sort.RANK.name, 
      'order': ASC, 
      'offset': 0, 
      'limit': MAX_LIST_LEN,
      'meta': true
    })
    // console.log(coins[0].rate)
    if (coins.length != 0) db.assign({coins: coins}).write()
}, {
    scheduled: true,
    timezone: "Asia/Ho_Chi_Minh"
})


export default db
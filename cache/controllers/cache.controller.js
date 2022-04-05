import { ASC, DESC, Sort } from "../config.js";
import db from "../db.js";

export function postCoinList(req, res, next) {
  const validate = (req_data) => {
    const sort = [Sort.AGE, Sort.CODE, Sort.NAME, Sort.PRICE, Sort.RANK, Sort.VOLUME]
    const order = [ASC, DESC]
    return req_data.currency === 'USD' && [true, false].some(_ => _ === req_data.meta) 
      && sort.some(_ => _.name === req_data.sort) && order.some(_ => _ === req_data.order)
  }

  if (validate(req.body)) {
    const data = [...db.get('coins').value()].splice(req.body.offset,req.body.limit)
    console.log(data)
    res.send(req.body.meta ? data : data.map(el => {
      return {
        code: el.code,
        rate: el.rate,
        volume: el.volume,
        cap: el.cap
      }
    }))
  } else {
    res.send('error')
  }
}

export function postSpecificCoin(req, res, next) {
  const validate = (req_data) => {
    return req_data.currency === 'USD' && [true, false].some(_ => _ === req_data.meta)
    && db.get('coins').value().some(_ => toString(_.code).toLowerCase() === toString(req_data.code).toLowerCase())
  }
  
  if (validate(req.body)) {
    const data = [...db.get('coins').value()]
    const el = data.filter(el => el.code.toLowerCase() === req.body.code.toLowerCase())
    if (el.length !== 0) {
      res.send(req.body.meta ? el[0] : {
        code: el[0].code,
        rate: el[0].rate,
        volume: el[0].volume,
        cap: el[0].cap
      })
    } else {
      res.send('404 Not found')
    }
  } else {
    res.send('error')
  }
}

export function postFiats(req,res,next) {
  const data = db.get('fiats').value()
  res.send(data)
}

export function getCoinList(req, res, next) {
  res.send('Here')
}


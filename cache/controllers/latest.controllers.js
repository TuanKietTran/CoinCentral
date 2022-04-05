import db from '../db.js'

export function getLatest(req, res, next) {
  const data = db.get('coins').value().map(el => {
    return {
      code: el.code,
      rate: el.rate,
      volume: el.volume,
      cap: el.cap
    }
  })
  res.send(data)
}

export function getACoinLatest(req,res,next) {

  let data = [...db.get('coins').value()]
  let el = data.filter(el => el.code.toLowerCase() === req.params.code.toLowerCase())
  if (el.length > 0) {
    res.send( {
      code: el[0].code,
      rate: el[0].rate,
      volume: el[0].volume,
      cap: el[0].cap
    })
  }
  else {
    res.send('404 Not found')
  }
}
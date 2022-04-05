import express from 'express'
import { COINS_LIST, COINS_SINGLE, FIATS } from '../config.js'
const router = express.Router()

import { getCoinList, postCoinList, postFiats, postSpecificCoin } from '../controllers/cache.controller.js'

router.get(COINS_LIST, getCoinList)

router.post(COINS_LIST, postCoinList)
router.post(COINS_SINGLE, postSpecificCoin)
router.post(FIATS, postFiats)

export default router
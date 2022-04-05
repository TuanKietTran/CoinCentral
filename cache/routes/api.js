import express from 'express'
import { COINS_SINGLE } from '../config.js'
import { getACoinLatest, getLatest} from '../controllers/latest.controllers.js'
const router = express.Router()


router.get('/latest', getLatest)
router.get('/latest/:code', getACoinLatest)


export default router
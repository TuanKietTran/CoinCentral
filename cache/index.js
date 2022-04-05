import 'dotenv/config'
import express from 'express'
import bodyParser from 'body-parser'
import api from './routes/api.js'
import cache from './routes/cache.js'

// Create the express app
const app = express()

// Routes and middleware
// app.use(/* ... */)
// app.get(/* ... */)
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json())
app.use('/api/', api)

app.use('/cache/', cache)

app.get('/', (req,res,next) => {
  res.send("Hello, express")
})

// Error handlers
app.use(function fourOhFourHandler (req, res) {
  res.status(404).send()
})
app.use(function fiveHundredHandler (err, req, res, next) {
  console.error(err)
  res.status(500).send()
})



// Start server
app.listen(process.env.PORT)

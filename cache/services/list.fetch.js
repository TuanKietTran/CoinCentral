import axios from 'axios'
import { API_KEY } from '../config.js';


async function fetchPOST(url, input) {
	try {
		const response = await axios({
			method: 'post',
			baseURL: url,
			data: input,
			headers: {
				'Content-Type': 'application/json',
				'x-api-key': API_KEY
			}
		})
		const data = response.data
		return data
	} catch(error) {
		if (error.response) {
			/*
			 * The request was made and the server responded with a
			 * status code that falls out of the range of 2xx
			 */
			console.log(error.response.data);
			console.log(error.response.status);
			console.log(error.response.headers);
		} else if (error.request) {
			/*
			 * The request was made but no response was received, `error.request`
			 * is an instance of XMLHttpRequest in the browser and an instance
			 * of http.ClientRequest in Node.js
			 */
			console.log(error.request);
		} else {
			// Something happened in setting up the request and triggered an Error
			console.log('Error', error.message);
		}
		console.log(error);
	}
}

export default fetchPOST
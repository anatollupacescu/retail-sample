import axios from 'axios'

export var apiIsHealthy = async (url: string): Promise<boolean> => {
  let isHealthy = false

  try {
    const response = await axios.get(`${url}/health`)
    isHealthy = response.status === 200
  } catch (error) {
    console.error(error)
  }

  return isHealthy
}

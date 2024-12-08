import axios, { AxiosError } from 'axios'

export const instance = axios.create({
  withCredentials: true,
  baseURL: '/api',
})

instance.interceptors.request.use((config) => {
  config.headers.Authorization = `Bearer ${localStorage.getItem('access')}`
  config.headers['Content-Type'] = 'application/json'
  return config
})

instance.interceptors.response.use(
  (config) => {
    return config
  },
  async (error) => {
    const axiosError = error as AxiosError
    if (!axiosError.config) throw error
    const originalRequest = axiosError.config

    if (
      axiosError.response?.status === 401 &&
      axiosError.config &&
      originalRequest.url !== '/users/refresh'
    ) {
      try {
        const response = await instance.get('/users/refresh')
        localStorage.setItem('access', response.data.access)
        return instance.request(originalRequest)
      } catch {
        console.warn('An error occurred while trying to refresh the token')
      }
    }
    throw error
  }
)

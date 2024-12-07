import { instance } from '@/api.config'

const AuthService = {
  login: async (login: string, password: string) =>
    await instance.post('/users/login', { login, password }),
  refresh: async () => await instance.get('/users/refresh'),
  logout: async () => await instance.post('/users/logout'),
}

export default AuthService

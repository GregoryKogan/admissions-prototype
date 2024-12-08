// Utilities
import AuthService from '@/api.auth'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    isAuth: false,
    isAuthInProgress: false,
  }),
  actions: {
    async login(login: string, password: string) {
      this.isAuthInProgress = true
      try {
        const response = await AuthService.login(login, password)
        localStorage.setItem('access', response.data.access)
        this.isAuth = true
        return response
      } catch (e) {
        this.isAuth = false
        throw e
      } finally {
        this.isAuthInProgress = false
      }
    },
    async checkAuth() {
      this.isAuthInProgress = true
      try {
        const response = await AuthService.refresh()
        localStorage.setItem('access', response.data.access)
        this.isAuth = true
      } catch {
        console.error('An error occurred while trying to refresh the token')
      } finally {
        this.isAuthInProgress = false
      }
    },
    async logout() {
      this.isAuthInProgress = true
      try {
        await AuthService.logout()
        this.isAuth = false
        localStorage.removeItem('access')
      } catch (e) {
        console.error('An error occurred while trying to logout', e)
      } finally {
        this.isAuthInProgress = false
      }
    },
    async me() {
      try {
        const response = await AuthService.me()
        return response.data
      } catch (e) {
        this.isAuth = false
        console.error('An error occurred while trying to get user info', e)
      }
    },
  },
})

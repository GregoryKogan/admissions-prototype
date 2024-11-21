// Utilities
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    isAuth: false,
  }),
  actions: {
    setIsAuth(isAuth: boolean) {
      this.isAuth = isAuth
    },
  },
})

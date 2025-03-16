import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '../types'
import { fetchUser } from '../api'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)

  async function fetchUserData() {
    try {
      const data = await fetchUser()
      user.value = data
    } catch (error) {
      console.error('获取用户信息失败:', error)
    }
  }

  return {
    user,
    fetchUser: fetchUserData
  }
}) 
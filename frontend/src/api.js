import { getApiUrl } from './utils'
import { store } from './store'
const apiUrl = getApiUrl()

const addAuthHeader = () => {
  return {
    'Authorization': `Bearer ${store.jwt}`
  }
} 

const addJsonHeader = () => {
  return {
    'Content-Type': 'application/json'
  }
}

export const getUsers = async () => {
  const response = await fetch(`${apiUrl}/api/users`, {
    method: 'GET',
    headers: {
      ...addAuthHeader(),
      ...addJsonHeader()
    }
  })
  const data = await response.json()
  return data.data
} 

export const createUser = (user) => {
  return new Promise(async (resolve, reject) => {
    try {
      const response = await fetch(`${apiUrl}/api/users`, {
        method: 'POST',
        headers: {
          ...addAuthHeader(),
          ...addJsonHeader()
        },
        body: JSON.stringify(user)
      })
      const data = await response.json()
      if (!response.ok) {
        reject(data)
        return
      }
      resolve(data.data)
    } catch (error) {
      reject(error)
    }
  })
}

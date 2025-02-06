import { getApiUrl } from './utils'
import { store } from './store'
import { jwtDecode } from 'jwt-decode'
const apiUrl = getApiUrl()

const addAuthHeader = () => {
  return {
    Authorization: `Bearer ${store.jwt}`
  }
}

const addJsonHeader = () => {
  return {
    'Content-Type': 'application/json'
  }
}

//auth login
export const login = async (username, password) => {
  const response = await fetch(`${apiUrl}/api/auth/login`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      ...addJsonHeader()
    },
    body: JSON.stringify({
      username: username,
      password: password
    })
  })
  if (!response.ok) {
    const responseData = await response.json()
    throw new Error(responseData.message)
  }
  const responseData = await response.json()
  return buildAuthSuccessData(responseData)
}

//auth refresh
export const refresh = async () => {
  const response = await fetch(`${apiUrl}/api/auth/refresh`, {
    method: 'POST',
    credentials: 'include'
  })
  if (!response.ok) {
    const responseData = await response.json()
    throw new Error(responseData.message)
  }
  const responseData = await response.json()
  return buildAuthSuccessData(responseData)
}

//auth logout
export const logout = async () => {
  try {
    await fetch(`${apiUrl}/api/auth/logout`, {
      method: 'POST',
      credentials: 'include'
    })
  } catch (error) {
    //ignore
  }

  store.setMultiple({
    admin: false,
    loggedIn: false,
    jwt: '',
    jwtExpires: null
  })
  return true
}

//get users
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

//create user
export const createUser = user => {
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

//delete user
export const deleteUser = id => {
  return new Promise(async (resolve, reject) => {
    try {
      const response = await fetch(`${apiUrl}/api/users/${id}`, {
        method: 'DELETE',
        headers: {
          ...addAuthHeader(),
          ...addJsonHeader()
        }
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

//private functions

const buildAuthSuccessData = data => {
  const decoded = jwtDecode(data.data.access_token)
  return {
    userId: decoded.sub,
    admin: decoded.admin,
    loggedIn: true,
    jwtExpires: decoded.exp,
    jwt: data.data.access_token
  }
}

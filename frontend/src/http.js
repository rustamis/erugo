import { store } from './store'
import { refresh,logout } from './api'


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

export const http = (url, options, withAuth = false, isJson = true, retries = 0, lastErrorCode = null) => {
  console.log(url, options, withAuth, isJson, retries, lastErrorCode)
  if (retries > 1) {
    if(lastErrorCode === 401) {
      logout()
    }
    throw new Error('Too many retries')
  }
  const headers = {
    ...(withAuth ? addAuthHeader() : {}),
    ...(isJson ? addJsonHeader() : {})
  }

  const response = fetch(url, { ...options, headers })

  // when the promise resolves, we need to check if the response is ok, in the case of a 401 we need to refresh the token and try again
  return response.then(async (response) => {
    if (response.status === 401) {
      await refreshToken()
      return http(url, options, withAuth, isJson, retries + 1, 401)
    }
    const data = await response.json()
    if(response.ok) {
      return data
    }
    throw new Error(data.message)
  })
}

export const httpGET = (url, options, withAuth = false, isJson = true, retries = 0) => {
  return http(url, { ...options, method: 'GET' }, withAuth, isJson, retries)
}

export const httpPOST = (url, options, withAuth = false, isJson = true, retries = 0) => {
  return http(url, { ...options, method: 'POST' }, withAuth, isJson, retries)
}



const refreshToken = async () => {
  try {
    const data = await refresh()
    store.authSuccess(data)
  } catch (error) {

  }
}
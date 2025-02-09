import { getApiUrl } from './utils'
import { store } from './store'
import { jwtDecode } from 'jwt-decode'

const apiUrl = getApiUrl()

const addAuthHeader = () => ({
  Authorization: `Bearer ${store.jwt}`
})

const addJsonHeader = () => ({
  'Content-Type': 'application/json'
})

// Wrapper for fetch that handles auth refresh
const fetchWithAuth = async (url, options = {}) => {
  // Add auth header if not present
  if (!options.headers?.Authorization) {
    options.headers = {
      ...options.headers,
      ...addAuthHeader()
    }
  }

  try {
    const response = await fetch(url, options)

    // If not 401, return response as-is
    if (response.status !== 401) {
      return response
    }

    // Handle 401 - try to refresh token
    try {
      const refreshData = await refresh()
      
      // Update auth header with new token
      options.headers = {
        ...options.headers,
        Authorization: `Bearer ${refreshData.jwt}`
      }

      // Retry original request with new token
      return await fetch(url, options)
    } catch (refreshError) {
      // If refresh fails, logout user
      store.setMultiple({
        admin: false,
        loggedIn: false,
        jwt: '',
        jwtExpires: null
      })
      throw new Error('Session expired. Please login again.')
    }
  } catch (error) {
    throw error
  }
}

// Auth Methods (these don't use fetchWithAuth since they handle auth directly)
export const login = async (username, password) => {
  const response = await fetch(`${apiUrl}/api/auth/login`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      ...addJsonHeader()
    },
    body: JSON.stringify({
      username,
      password
    })
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return buildAuthSuccessData(data)
}

export const refresh = async () => {
  const response = await fetch(`${apiUrl}/api/auth/refresh`, {
    method: 'POST',
    credentials: 'include'
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return buildAuthSuccessData(data)
}

export const logout = async () => {
  try {
    await fetch(`${apiUrl}/api/auth/logout`, {
      method: 'POST',
      credentials: 'include'
    })
  } catch (error) {
    // ignore
  }

  store.setMultiple({
    admin: false,
    loggedIn: false,
    jwt: '',
    jwtExpires: null
  })
  return true
}

// User Methods
export const getUsers = async () => {
  const response = await fetchWithAuth(`${apiUrl}/api/users`, {
    method: 'GET',
    headers: {
      ...addJsonHeader()
    }
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return data.data
}

export const createUser = async (user) => {
  const response = await fetchWithAuth(`${apiUrl}/api/users`, {
    method: 'POST',
    headers: {
      ...addJsonHeader()
    },
    body: JSON.stringify(user)
  })
  const data = await response.json()
  if (!response.ok) {
    return Promise.reject(data)
  }
  return data.data
}

export const deleteUser = async (id) => {
  const response = await fetchWithAuth(`${apiUrl}/api/users/${id}`, {
    method: 'DELETE',
    headers: {
      ...addJsonHeader()
    }
  })
  const data = await response.json()
  if (!response.ok) {
    return Promise.reject(data)
  }
  return data.data
}

// Settings Methods
export const getSettingsByGroup = async (group) => {
  const response = await fetchWithAuth(`${apiUrl}/api/settings?group=${group}`, {
    method: 'GET',
    headers: {
      ...addJsonHeader()
    }
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return data.data
}

export const getSettingById = async (id) => {
  const response = await fetchWithAuth(`${apiUrl}/api/settings/${id}`, {
    method: 'GET',
    headers: {
      ...addJsonHeader()
    }
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return data.data
}

export const saveSettingsById = async (settings) => {
  const response = await fetchWithAuth(`${apiUrl}/api/settings`, {
    method: 'PUT',
    headers: {
      ...addJsonHeader()
    },
    body: JSON.stringify(settings)
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return data.data
}

export const saveLogo = async (logoFile) => {
  const formData = new FormData()
  formData.append('logo', logoFile)
  
  const response = await fetchWithAuth(`${apiUrl}/api/settings/logo`, {
    method: 'PUT',
    body: formData
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return data.data
}

// Private functions
const buildAuthSuccessData = (data) => {
  const decoded = jwtDecode(data.data.access_token)
  return {
    userId: decoded.sub,
    admin: decoded.admin,
    loggedIn: true,
    jwtExpires: decoded.exp,
    jwt: data.data.access_token
  }
}
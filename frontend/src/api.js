import { getApiUrl } from './utils'
import { store } from './store'
import { jwtDecode } from 'jwt-decode'
import { useToast } from 'vue-toastification'
import debounce from './debounce'

const apiUrl = getApiUrl()
const toast = useToast()
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

    // If response is OK, return as-is
    if (response.ok) {
      return response
    }

    // Handle 401 or 403
    if (response.status === 401 || response.status === 403) {
      // Clone the response so we can read the body
      const clonedResponse = response.clone()
      const responseData = await clonedResponse.json()

      // Check for password change required in response body
      if (responseData?.message === 'Password change required') {
        store.setSettingsOpen(false)
        debouncedPasswordChangeRequired()
        throw new Error('PASSWORD_CHANGE_REQUIRED')
      }

      // For 401, try to refresh token
      if (response.status === 401) {
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
          // If refresh fails, proceed to logout
        }
      }

      // If we reach here, either:
      // 1. It was a 403 without password change required
      // 2. It was a 401 and token refresh failed
      // In both cases, we log the user out
      store.setMultiple({
        admin: false,
        loggedIn: false,
        jwt: '',
        jwtExpires: null
      })
      throw new Error('Session expired. Please login again.')
    }

    // Handle other error status codes
    return response
  } catch (error) {
    // Rethrow password change required error
    if (error.message === 'PASSWORD_CHANGE_REQUIRED') {
      throw error
    }
    // Handle other errors
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

export const createUser = async user => {
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

export const updateUser = async user => {
  const response = await fetchWithAuth(`${apiUrl}/api/users/${user.id}`, {
    method: 'PUT',
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

export const updateMyProfile = async user => {
  const response = await fetchWithAuth(`${apiUrl}/api/users/me`, {
    method: 'PUT',
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

export const deleteUser = async id => {
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
export const getSettingsByGroup = async group => {
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

export const getSettingById = async id => {
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

export const saveSettingsById = async settings => {
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

export const saveLogo = async logoFile => {
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

// Share Methods
export const createShare = async (files, name, description, uploadId) => {

  const formData = new FormData()
  files.forEach(file => {
    formData.append('files', file)
  })
  formData.append('name', name)
  formData.append('description', description)

  const response = await fetchWithAuth(`${apiUrl}/api/shares?uploadId=${uploadId}`, {
    method: 'POST',
    body: formData
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return data.data.share
}

export const getShare = async id => {
  const response = await fetchWithAuth(`${apiUrl}/api/shares/${id}`, {
    method: 'GET',
    headers: {
      ...addJsonHeader()
    }
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return data.data.share
}

//misc methods
export const getHealth = async () => {
  const response = await fetch(`${apiUrl}/api/health`)
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data.message)
  }
  return data.data
}

export const getMyProfile = async () => {
  const response = await fetchWithAuth(`${apiUrl}/api/users/me`, {
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

export const createFirstUser = async user => {
  const response = await fetch(`${apiUrl}/api/setup`, {
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

// Private functions
const buildAuthSuccessData = data => {
  const decoded = jwtDecode(data.data.access_token)
  return {
    userId: decoded.sub,
    admin: decoded.admin,
    loggedIn: true,
    jwtExpires: decoded.exp,
    jwt: data.data.access_token,
    mustChangePassword: decoded.must_change_password
  }
}

const passwordChangeRequired = () => {
  toast.error('You must change your password to continue')
  store.showPasswordResetForm()
}

const debouncedPasswordChangeRequired = debounce(passwordChangeRequired, 100)

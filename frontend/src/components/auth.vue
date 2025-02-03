<script setup>
  import { ref, onMounted } from 'vue'
  import { getApiUrl } from '../utils'
  import { jwtDecode } from 'jwt-decode'
  import { store } from '../store'
  const apiUrl = getApiUrl()

  const username = ref('')
  const password = ref('')
  const passwordInput = ref(null)

  onMounted(() => {
    attemptRefresh()
  })

  const attemptRefresh = () => {
    fetch(`${apiUrl}/api/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      }
    }).then(response => {
      if (response.ok) {
        response.json().then(data => {
          authSuccess(data)
        })
      }
    })
  }

  const authSuccess = (data) => {
    const decoded = jwtDecode(data.data.access_token)
    store.setMultiple({
      userId: decoded.sub,
      admin: decoded.admin,
      loggedIn: true,
      jwtExpires: decoded.exp,
      jwt: data.data.access_token
    })
    store.logState()
  }

  const login = () => {
    fetch(`${apiUrl}/api/auth/login`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: username.value,
        password: password.value
      })
    }).then(response => {
      if (response.ok) {
        response.json().then(data => {
          authSuccess(data)
        })
      }
    })
  }

  const logout = () => {
    fetch(`${apiUrl}/api/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    }).then(response => {
      username.value = ''
      password.value = ''

      if (response.ok) {
        store.setMultiple({
          admin: false,
          loggedIn: false,
          jwt: '',
          jwtExpires: null
        })
      }
    })
  }

  const moveToPassword = () => {
    passwordInput.value.focus()
  }

  //export the functions
  defineExpose({
    logout
  })
</script>

<template>
  <div class="auth-container">
    <div class="auth-container-inner">
      <h1>Login</h1>
      <p>Login to your erugo account to upload files.</p>
      <input type="text" v-model="username" placeholder="Username" @keyup.enter="moveToPassword"/>
      <input type="password" v-model="password" placeholder="Password" @keyup.enter="login" ref="passwordInput"/>
      <button @click="login" class="login-button mt-4">Login</button>
    </div>
  </div>
</template>
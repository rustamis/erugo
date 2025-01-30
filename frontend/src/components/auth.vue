<script setup>
  import { ref, onMounted } from 'vue'
  import { getApiUrl } from '../utils'

  const apiUrl = getApiUrl()

  const username = ref('')
  const password = ref('')
  const loggedIn = ref(false)
  const passwordInput = ref(null)

  const emit = defineEmits(['auth-success', 'logout'])

  onMounted(() => {
    attemptRefresh()
  })

  const attemptRefresh = () => {
    fetch(`${apiUrl}/api/auth/refresh/`, {
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
    emit('auth-success', data)
    loggedIn.value = true
  }

  const login = () => {
    fetch(`${apiUrl}/api/auth/login/`, {
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
    fetch(`${apiUrl}/api/auth/logout/`, {
      method: 'POST',
      credentials: 'include',
    }).then(response => {
      if (response.ok) {
        emit('logout')
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
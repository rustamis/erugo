<script setup>
  import { ref, onMounted } from 'vue'
  import { getApiUrl } from '../utils'
  import { useToast } from 'vue-toastification'

  import { store } from '../store'
  import { login, refresh, logout } from '../api'
  const apiUrl = getApiUrl()
  const toast = useToast()
  const username = ref('')
  const password = ref('')
  const passwordInput = ref(null)

  onMounted(() => {
    attemptRefresh()
  })

  const attemptLogin = async () => {
    if (username.value === '' || password.value === '') {
      toast.error('Please enter a username and password')
      return
    }

    try {
      const data = await login(username.value, password.value)
      store.authSuccess(data)
      toast.success('Login successful')
    } catch (error) {
      console.log('error', error)
      toast.error('Invalid username or password')
    }
  }

  const attemptRefresh = () => {
    refresh()
      .then(data => {
        store.authSuccess(data)
      })
      .catch(error => {
        //noop
      })
  }


  const attemptLogout = async () => {
    await logout()
  }

  const moveToPassword = () => {
    passwordInput.value.focus()
  }
</script>

<template>
  <div class="auth-container">
    <div class="auth-container-inner">
      <h1>Login</h1>
      <p>Login to your erugo account to upload files.</p>
      <input type="text" v-model="username" placeholder="Username" @keyup.enter="moveToPassword" />
      <input type="password" v-model="password" placeholder="Password" @keyup.enter="attemptLogin" ref="passwordInput" />
      <button @click="attemptLogin" class="login-button mt-4">Login</button>
    </div>
  </div>
</template>

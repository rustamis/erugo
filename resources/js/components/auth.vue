<script setup>
  import { ref, onMounted } from 'vue'
  import { getApiUrl } from '../utils'
  import { useToast } from 'vue-toastification'
  import { domData } from "../domData";
  import { KeyRound } from 'lucide-vue-next';
  import { store } from '../store'
  import { login, refresh, logout } from '../api'
  const apiUrl = getApiUrl()
  const toast = useToast()
  const email = ref('')
  const password = ref('')
  const passwordInput = ref(null)
  const loginMessage = domData().login_message

  onMounted(() => {
    attemptRefresh()
  })

  const attemptLogin = async () => {
    if (email.value === '' || password.value === '') {
      toast.error('Please enter an email and password')
      return
    }

    try {
      const data = await login(email.value, password.value)
      store.authSuccess(data)
      toast.success('Login successful')
    } catch (error) {
      toast.error('Invalid email or password')
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
      <h1>Welcome</h1>
      <p>{{ loginMessage }}</p>
      <input type="text" v-model="email" placeholder="Email" @keyup.enter="moveToPassword" />
      <input type="password" v-model="password" placeholder="Password" @keyup.enter="attemptLogin" ref="passwordInput" />
      <button @click="attemptLogin" class="login-button mt-4"><KeyRound /> Login</button>
    </div>
  </div>
</template>

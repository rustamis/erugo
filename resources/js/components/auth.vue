<script setup>
import { ref, onMounted } from 'vue'
import { getApiUrl } from '../utils'
import { useToast } from 'vue-toastification'
import { domData } from '../domData'
import { KeyRound, Annoyed } from 'lucide-vue-next'
import { store } from '../store'
import { login, refresh, logout, forgotPassword, resetPassword } from '../api'

const apiUrl = getApiUrl()
const toast = useToast()
const email = ref('')
const password = ref('')
const password_confirmation = ref('')
const passwordInput = ref(null)
const loginMessage = domData().login_message
const forgotPasswordMode = ref(false)
const haveResetToken = ref(false)
const resetToken = ref('')
const waitingForRedirect = ref(false)

onMounted(() => {
  attemptRefresh()
  console.log('domData', domData())
  const token = domData().token
  if (token) {
    haveResetToken.value = true
    resetToken.value = token
  }
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
    .then((data) => {
      store.authSuccess(data)
    })
    .catch((error) => {
      //noop
    })
}

const attemptLogout = async () => {
  await logout()
}

const moveToPassword = () => {
  passwordInput.value.focus()
}

const attemptForgotPassword = async () => {
  console.log('attemptForgotPassword', email.value)
  try {
    await forgotPassword(email.value)
    toast.success('Password reset email sent')
    forgotPasswordMode.value = false
  } catch (error) {
    toast.error('Failed to send password reset email')
  }
}

const attemptResetPassword = async () => {
  if (password.value === '' || password_confirmation.value === '') {
    toast.error('Please enter a password and confirm password')
    return
  }
  if (password.value !== password_confirmation.value) {
    toast.error('Passwords do not match')
    return
  }
  try {
    await resetPassword(resetToken.value, email.value, password.value, password_confirmation.value)
    toast.success('Password reset successfully')
    haveResetToken.value = false
    waitingForRedirect.value = true
    setTimeout(() => {
      window.location.href = '/'
    }, 3000)
  } catch (error) {
    toast.error('Failed to reset password')
  }
}
</script>

<template>
  <div class="auth-container" v-if="!haveResetToken && !waitingForRedirect">
    <div class="auth-container-inner">
      <template v-if="!forgotPasswordMode">
        <h1>Welcome</h1>
        <p>{{ loginMessage }}</p>
      </template>
      <template v-else>
        <h1>Forgot Password</h1>
        <p>Please enter your email to reset your password</p>
      </template>
      <div class="input-container">
        <label for="email">Email</label>
        <input type="text" v-model="email" placeholder="Email" @keyup.enter="moveToPassword" />
      </div>
      <div class="input-container" v-if="!forgotPasswordMode">
        <label for="password">Password</label>
        <input
          type="password"
          v-model="password"
          placeholder="Password"
          @keyup.enter="attemptLogin"
          ref="passwordInput"
        />
      </div>
      <div class="row mt-3 align-items-center" v-if="!forgotPasswordMode">
        <div class="col">
          <button class="block" @click="attemptLogin">
            <KeyRound />
            Login
          </button>
        </div>
        <div class="col">
          <a href="" @click.prevent="forgotPasswordMode = true">Forgot Password?</a>
        </div>
      </div>
      <div class="row mt-3 align-items-center" v-if="forgotPasswordMode">
        <div class="col">
          <button class="block" @click="attemptForgotPassword">
            <KeyRound />Request&nbsp;Reset
          </button>
        </div>
        <div class="col">
          <a href="" @click.prevent="forgotPasswordMode = false">Back to Login</a>
        </div>
      </div>
    </div>
  </div>

  <div class="auth-container" v-else-if="!waitingForRedirect">
    <div class="auth-container-inner">
      <h1>Create Password</h1>
      <p>Please create your new password.</p>
      <div class="input-container">
        <label for="email">Email</label>
        <input type="text" v-model="email" placeholder="Email" @keyup.enter="moveToPassword" />
      </div>
      <div   class="input-container">
        <label for="password">Password</label>
        <input type="password" v-model="password" placeholder="Password" @keyup.enter="attemptResetPassword" />
        <label for="password_confirmation">Confirm Password</label>
        <input type="password" v-model="password_confirmation" placeholder="Confirm Password" @keyup.enter="attemptResetPassword" />
      </div>
      <div class="row mt-3 align-items-center">
        <div class="col">
          <button class="block" @click="attemptResetPassword">
            <KeyRound />Save new password
          </button>
        </div>
      </div>
    </div>
  </div>
  <div class="auth-container" v-else>
    <div class="auth-container-inner">
      <h1>Your password has been set!</h1>
      <p>Hold tight while we redirect you to the dashboard.</p>
    </div>
  </div>
</template>

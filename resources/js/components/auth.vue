<script setup>
import { ref, onMounted } from 'vue'
import { getApiUrl } from '../utils'
import { useToast } from 'vue-toastification'
import { domData } from '../domData'
import { KeyRound, Annoyed } from 'lucide-vue-next'
import { store } from '../store'
import { login, refresh, logout, forgotPassword, resetPassword } from '../api'

import { useTranslate } from '@tolgee/vue'

const { t } = useTranslate()

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
    toast.error(t.value('auth.please_enter_email_and_password'))
    return
  }

  try {
    const data = await login(email.value, password.value)
    store.authSuccess(data)
    toast.success(t.value('auth.login_successful'))
  } catch (error) {
    toast.error(t.value('auth.invalid_email_or_password'))
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
  if (email.value === '') {
    toast.error(t.value('auth.please_enter_email'))
    return
  }
  try {
    await forgotPassword(email.value)
    toast.success(t.value('auth.password_reset_email_sent'))
    forgotPasswordMode.value = false
  } catch (error) {
    toast.error(t.value('auth.failed_to_send_password_reset_email'))
  }
}

const attemptResetPassword = async () => {
  if (password.value === '' || password_confirmation.value === '') {
    toast.error(t.value('auth.please_enter_password_and_confirm_password'))
    return
  }
  if (password.value !== password_confirmation.value) {
    toast.error(t.value('auth.passwords_do_not_match'))
    return
  }
  try {
    await resetPassword(resetToken.value, email.value, password.value, password_confirmation.value)
    toast.success(t.value('auth.password_reset_successfully'))
    haveResetToken.value = false
    waitingForRedirect.value = true
    setTimeout(() => {
      window.location.href = '/'
    }, 3000)
  } catch (error) {
    toast.error(t.value('auth.failed_to_reset_password'))
  }
}
</script>

<template>
  <div class="auth-container" v-if="!haveResetToken && !waitingForRedirect">
    <div class="auth-container-inner">
      <template v-if="!forgotPasswordMode">
        <h1>{{ $t('auth.welcome') }}</h1>
        <p>{{ loginMessage }}</p>
      </template>
      <template v-else>
        <h1>{{ $t('auth.forgot_password') }}</h1>
        <p>{{ $t('auth.please_enter_email_to_reset_password') }}</p>
      </template>
      <div class="input-container">
        <label for="email">{{ $t('auth.email') }}</label>
        <input type="text" v-model="email" :placeholder="$t('auth.email')" @keyup.enter="moveToPassword" />
      </div>
      <div class="input-container" v-if="!forgotPasswordMode">
        <label for="password">{{ $t('auth.password') }}</label>
        <input
          type="password"
          v-model="password"
          :placeholder="$t('auth.password')"
          @keyup.enter="attemptLogin"
          ref="passwordInput"
        />
      </div>
      <div class="row mt-3 align-items-center" v-if="!forgotPasswordMode">
        <div class="col">
          <button class="block" @click="attemptLogin">
            <KeyRound />
            {{ $t('auth.login') }}
          </button>
        </div>
        <div class="col">
          <a href="" @click.prevent="forgotPasswordMode = true">{{ $t('auth.forgot_password') }}</a>
        </div>
      </div>
      <div class="row mt-3 align-items-center" v-if="forgotPasswordMode">
        <div class="col">
          <button class="block" @click="attemptForgotPassword">
            <KeyRound />{{ $t('auth.request_reset') }}
          </button>
        </div>
        <div class="col">
          <a href="" @click.prevent="forgotPasswordMode = false">{{ $t('auth.back_to_login') }}</a>
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

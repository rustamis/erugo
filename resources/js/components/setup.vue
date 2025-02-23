<script setup>
import { ref } from 'vue'
import { UserPlus, Hand } from 'lucide-vue-next'
import { getApiUrl } from '../utils'
import { createFirstUser, login } from '../api'
import { useToast } from 'vue-toastification'

import { store } from '../store'

const apiUrl = getApiUrl()
const toast = useToast()
const logoUrl = `${apiUrl}/logo`

const newUser = ref({
  username: '',
  name: '',
  email: '',
  password: '',
  password_confirmation: ''
})

const errors = ref({
  name: '',
  password: '',
  email: '',
  password_confirmation: ''
})

const saveUser = async () => {
  errors.value = {}
  if (newUser.value.password !== newUser.value.password_confirmation) {
    errors.value.password_confirmation = 'Password confirmation does not match'
  }

  if (Object.keys(errors.value).length > 0) {
    toast.error('Please fix the errors before saving')
    return
  }

  try {
    await createFirstUser(newUser.value)
    toast.success('User created successfully! Logging you in...')
    setTimeout(async () => {
      try {
        const data = await login(newUser.value.email, newUser.value.password)
        store.authSuccess(data)
        window.location.href = '/'
      } catch (error) {
        toast.error('Failed to login')
      }
    }, 1000)
  } catch (error) {
    errors.value = error.data.errors
    toast.error('Failed to create user')
  }
}
</script>

<template>
  <div class="setup-container">
    <div class="setup-inner">
      <img :src="logoUrl" alt="Erugo" class="setup-logo" />
      <p>
        <em>Thank you</em>
        for installing erugo. Before you can use it, you need to create an admin account.
      </p>

      <hr />
      <div class="setup-form">
        <!-- email -->
        <div class="input-container mt-2">
          <label for="email">Email</label>
          <input type="email" v-model="newUser.email" placeholder="Email" required :class="{ error: errors.email }" />
          <div class="error-message" v-if="errors.email">
            {{ errors.email }}
          </div>
        </div>

        <!-- full name -->
        <div class="input-container mt-2">
          <label for="name">Full Name</label>
          <input type="text" v-model="newUser.name" placeholder="Full Name" required :class="{ error: errors.name }" />
          <div class="error-message" v-if="errors.name">
            {{ errors.name }}
          </div>
        </div>

        <!-- password -->
        <div class="input-container mt-2">
          <label for="password">Password</label>
          <input
            type="password"
            v-model="newUser.password"
            placeholder="Password"
            required
            :class="{ error: errors.password }"
          />
          <div class="error-message" v-if="errors.password">
            {{ errors.password }}
          </div>
        </div>

        <!-- confirm password -->
        <div class="input-container mt-2">
          <label for="password_confirmation">Confirm Password</label>
          <input
            type="password"
            v-model="newUser.password_confirmation"
            placeholder="Confirm Password"
            required
            :class="{ error: errors.password_confirmation }"
          />
          <div class="error-message" v-if="errors.password_confirmation">
            {{ errors.password_confirmation }}
          </div>
        </div>

        <div class="button-bar mt-3">
          <button @click="saveUser">
            <UserPlus />
            Create Admin Account
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.setup-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: var(--primary-color);
  z-index: 230;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  backdrop-filter: blur(10px);

  .setup-inner {
    background-color: var(--accent-color-light-transparent);
    padding: 20px;
    border-radius: 0 0 10px 10px;
    width: 30%;

    h1 {
      font-size: 24px;
      color: var(--primary-color);
      display: flex;
      align-items: center;
      gap: 10px;
      svg {
        font-size: 24px;
        margin-right: 5px;
        margin-top: -1px;
      }
    }
  }
}

button {
  display: block;
  width: 100%;
}

.setup-logo {
  width: 100px;
  margin-top: 5px;
  margin-bottom: 15px;
}
</style>

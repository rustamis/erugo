<script setup>
  import { ref, onMounted } from 'vue'
  import { getMyProfile, updateMyProfile } from '../../api'
  import { User, UserPen, CircleX, UserRoundCheck, UserRoundPen } from 'lucide-vue-next'
  import { useToast } from 'vue-toastification'
  import { store } from '../../store'

  const toast = useToast()
  const profile = ref(null)
  const editUserFormActive = ref(false)
  const editUser = ref({})
  const errors = ref({})

  onMounted(async () => {
    profile.value = await getMyProfile()
    editUser.value = {
      ...profile.value,
      password: null,
      password_confirmation: null,
      current_password: null
    }
    if (store.autoShowProfileEdit) {
      editUserFormActive.value = true
      store.autoShowProfileEdit = false
    }
  })

  const formatDate = date => {
    return new Date(date).toLocaleDateString('en-GB', { year: 'numeric', month: 'long', day: 'numeric' })
  }

  const editUserFormClickOutside = e => {
    if (!e.target.closest('.user-form')) {
      editUserFormActive.value = false
    }
  }

  const saveUser = async () => {
    errors.value = {}
    if (editUser.value.password !== editUser.value.password_confirmation) {
      errors.value.password_confirmation = 'Password confirmation does not match'
    }

    if (editUser.value.password && editUser.value.current_password === null) {
      errors.value.current_password = 'Current password is required'
    }

    if (Object.keys(errors.value).length > 0) {
      toast.error('Please fix the errors before saving')
      return
    }

    try {
      const updatedUser = await updateMyProfile(editUser.value)
      profile.value = updatedUser
      editUserFormActive.value = false
      toast.success('Profile updated successfully')
    } catch (error) {
      toast.error('Failed to update profile')
      errors.value = error.data.errors
    }
  }
</script>

<template>
  <div class="p-5" v-if="profile">
    <div class="profile-card">
      <div class="profile-card-header">
        <h1>
          <User />
          {{ profile.name }}
        </h1>
      </div>
      <div class="profile-card-tags">
        <!-- user is admin -->
        <span class="profile-card-tag" v-if="profile.admin">{{ $t('settings.users.admin') }}</span>
        <!-- user is active -->
        <span class="profile-card-tag" v-if="profile.active">{{ $t('settings.users.active') }}</span>
      </div>

      <div class="profile-card-profile-item">
        <h2>{{ $t('settings.users.email') }}</h2>
        <p>{{ profile.email }}</p>
      </div>

      <div class="profile-card-profile-item">
        <h2>{{ $t('settings.users.created') }}</h2>
        <p>{{ formatDate(profile.created_at) }}</p>
      </div>

      <div class="profile-card-footer">
        <button @click="editUserFormActive = true">
          <UserPen />
          {{ $t('settings.users.edit_my_profile') }}
        </button>
      </div>
    </div>
  </div>

  <div class="user-form-overlay" :class="{ active: editUserFormActive }" @click="editUserFormClickOutside">
    <div class="user-form">
      <h2>
        <UserRoundPen />
        {{ $t('settings.users.edit_my_profile') }}
      </h2>
      <div class="input-container">
        <input type="email" v-model="editUser.email" :placeholder="$t('settings.users.email')" required :class="{ error: errors.email }" />
        <div class="error-message" v-if="errors.email">
          {{ errors.email }}
        </div>
      </div>
      <div class="input-container">
        <input type="text" v-model="editUser.name" :placeholder="$t('settings.users.name')" required :class="{ error: errors.name }" />
        <div class="error-message" v-if="errors.name">
          {{ errors.name }}
        </div>
      </div>
      
      <label for="password" class="mb-3">{{ $t('settings.users.update_password') }}</label>
      <div class="input-container">
        <input type="password" v-model="editUser.current_password" :placeholder="$t('settings.users.current_password')" required :class="{ error: errors.current_password }" />
        <div class="error-message" v-if="errors.current_password">
          {{ errors.current_password }}
        </div>
      </div>
      <div class="input-container">
        <input type="password" v-model="editUser.password" :placeholder="$t('settings.users.password')" required :class="{ error: errors.password }" />
        <div class="error-message" v-if="errors.password">
          {{ errors.password }}
        </div>
      </div>
      <div class="input-container">
        <input type="password" v-model="editUser.password_confirmation" :placeholder="$t('settings.users.password_confirmation')" required :class="{ error: errors.password_confirmation }" />
        <div class="error-message" v-if="errors.password_confirmation">
          {{ errors.password_confirmation }}
        </div>
        <p class="help-text">{{ $t('settings.users.leave_blank_to_keep_same_password') }}</p>
      </div>

      <div class="button-bar">
        <button @click="saveUser">
          <UserRoundCheck />
          {{ $t('settings.users.save_changes') }}
        </button>
        <button class="secondary close-button" @click="editUserFormActive = false">
          <CircleX />
          {{ $t('settings.users.close') }}
        </button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .profile-card {
    width: 450px;
    border-radius: 10px;
    background: var(--panel-item-background-color);

    .profile-card-header {
      background: var(--panel-header-background-color);
      border-radius: 8px 8px 0 0;
      padding-left: 20px;
      padding-right: 20px;
      height: 80px;
      display: flex;
      align-items: center;
      justify-content: flex-start;
      h1 {
        font-size: 24px;
        font-weight: 600;
        color: var(--panel-text-color);
        display: flex;
        align-items: center;
        gap: 10px;
        svg {
          width: 24px;
          height: 24px;
          color: var(--panel-text-color);
        }
      }
    }

    .profile-card-tags {
      padding: 20px;
      background: var(--panel-item-background-color);
      border-radius: 8px 8px 0 0;
      margin-bottom: 0px;
      margin-right: 5px;
      margin-left: 5px;
      margin-top: -10px;
      display: flex;
      align-items: center;
      gap: 10px;
      .profile-card-tag {
        font-size: 14px;
        font-weight: 600;
        color: var(--panel-text-color);
        background: var(--panel-background-color);
        padding: 5px 10px;
        border-radius: 5px;
      }
    }

    .profile-card-profile-item {
      padding: 10px 20px;
      background: var(--panel-background-color);
      h2 {
        font-size: 16px;
        font-weight: 600;
        color: var(--panel-text-color);
      }
      p {
        font-size: 19px;
        font-weight: 400;
        color: var(--panel-text-color);
        padding: 0;
        margin: 0;
      }

      &:last-child {
        border-radius: 0 0 8px 8px;
        border-bottom: none;
      }
    }

    .profile-card-footer {
      padding: 10px 20px;
      display: flex;
      align-items: center;
      gap: 10px;
      button {
        display: block;
        width: 100%;
      }
    }
  }

  .user-form-overlay {
    border-radius: 10px 10px 0 0;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: var(--overlay-background-color);
    backdrop-filter: blur(10px);
    z-index: 230;
    opacity: 0;
    pointer-events: none;
    transition: all 0.3s ease;

    h2 {
      margin-bottom: 10px;
      font-size: 24px;
      color: var(--panel-text-color);
      display: flex;
      align-items: center;
      justify-content: center;

      svg {
        width: 24px;
        height: 24px;
        margin-right: 10px;
      }
    }
    .user-form {
      position: absolute;
      bottom: 0;
      left: 50%;
      transform: translate(-50%, 100%);
      width: 500px;
      background: var(--panel-background-color);
      color: var(--panel-text-color);
      padding: 20px;
      border-radius: 10px 10px 0 0;
      box-shadow: 0 0 100px 0 rgba(0, 0, 0, 0.5);
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      justify-content: flex-start;
      gap: 10px;
      transition: all 0.3s ease;
      padding-bottom: 20px;
      button {
        display: block;
        width: 100%;
      }
    }

    &.active {
      opacity: 1;
      pointer-events: auto;
      .user-form {
        transform: translate(-50%, 0%);
      }
    }
  }
</style>

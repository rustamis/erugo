<script setup>
  import { ref, onMounted } from 'vue'
  import { getUsers, createUser, deleteUser, updateUser } from '../../api'
  import { UserPen, Trash, UserPlus, CircleX, UserRoundCheck } from 'lucide-vue-next'
  import { store } from '../../store'
  import { useToast } from 'vue-toastification'
  import { niceDate } from '../../utils'

  const toast = useToast()
  const users = ref([])
  const errors = ref({})

  const newUser = ref({})
  const editUser = ref({})

  onMounted(async () => {
    loadUsers()
    newUser.value = getEmptyUser()
  })

  const loadUsers = () => {
    getUsers().then(data => {
      users.value = data.users
    })
  }

  const handleDeleteUserClick = id => {
    if (id === store.userId) {
      toast.error('You cannot delete yourself.')
      return
    }
    if (confirm(`Are you sure you want to delete user ${id}?`)) {
      deleteUser(id).then(data => {
        loadUsers()
        toast.success('User deleted successfully')
      }, error => {
        toast.error('Failed to delete user')
      })
    }
  }

  const handleEditUserClick = user => {
    editUser.value = user
    editUser.value.admin = user.admin == 1
    editUser.value.must_change_password = user.must_change_password == 1
    editUserFormActive.value = true
  }

  const addUser = () => {
    newUserFormActive.value = true
  }

  //expose addUser to parent
  defineExpose({
    addUser
  })

  const newUserFormActive = ref(false)

  const newUserFormClickOutside = event => {
    if (!event.target.closest('.user-form')) {
      newUserFormActive.value = false
    }
  }

  const editUserFormActive = ref(false)

  const editUserFormClickOutside = event => {
    if (!event.target.closest('.user-form')) {
      editUserFormActive.value = false
    }
  }

  const saveUser = () => {
    errors.value = {}

    if (newUserFormActive.value) {
      createUser(newUser.value).then(
        data => {
          loadUsers()
          newUserFormActive.value = false
          newUser.value = getEmptyUser()
          toast.success('User created successfully')
        },
        error => {
          errors.value = error.data.errors
          toast.error('Failed to create user')
        }
      )
    } else if (editUserFormActive.value) {
      updateUser(editUser.value).then(
        data => {
          loadUsers()
          editUserFormActive.value = false
          editUser.value = getEmptyUser() 
          toast.success('User updated successfully')
        },
        error => {
          errors.value = error.data.errors
          toast.error('Failed to update user')
        }
      )
    }
  }

  const getEmptyUser = () => {
    return {
      email: '',
      password: '',
      password_confirmation: '',
      name: '',
      admin: false,
      must_change_password: true
    }
  }
</script>

<template>
  <div>
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Email</th>
          <th>Full Name</th>
          <th>Is Admin</th>
          <th>Created At</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.id">
          <td width="1" style="white-space: nowrap">{{ user.id }}</td>
          <td>
            {{ user.email }}
            <div class="tags" v-if="user.id === store.userId || user.admin || user.must_change_password">
              <span v-if="user.id === store.userId" class="you-tag">You</span>
              <span v-if="user.admin" class="admin-tag">Admin</span>
              <span v-if="user.must_change_password" class="must-change-password-tag">Password change required</span>
            </div>
          </td>
          <td>{{ user.name }}</td>
          <td>{{ user.admin ? 'Yes' : 'No' }}</td>
          <td>{{ niceDate(user.created_at) }}</td>
          <td width="1" style="white-space: nowrap">
            <button :disabled="user.id === store.userId" @click="handleEditUserClick(user)">
              <UserPen />
              Edit
            </button>
            <button :disabled="user.id === store.userId" @click="handleDeleteUserClick(user.id)">
              <Trash />
              Delete
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <div class="user-form-overlay" :class="{ active: newUserFormActive }" @click="newUserFormClickOutside">
    <div class="user-form">
      <h2>
        <UserPlus />
        Add User
      </h2>
      <div class="input-container">
        <label for="new_user_email">Email</label>
        <input type="email" v-model="newUser.email" id="new_user_email" placeholder="Email" required :class="{ error: errors.email }" />
        <div class="error-message" v-if="errors.email">
          {{ errors.email[0] }}
        </div>
      </div>
      <div class="input-container">
        <label for="new_user_name">Full Name</label>
        <input type="text" v-model="newUser.name" id="new_user_name" placeholder="Full Name" required :class="{ error: errors.name }" />
        <div class="error-message" v-if="errors.name">
          {{ errors.name[0] }}
        </div>
      </div>
      
      <div class="input-container">
        <label for="new_user_password">Password</label>
        <input type="password" v-model="newUser.password" id="new_user_password" placeholder="Password" required :class="{ error: errors.password }" />
        <div class="error-message" v-if="errors.password">
          {{ errors.password[0] }}
        </div>
      </div>

      <div class="input-container">
        <label for="new_user_password_confirmation">Password Confirmation</label>
        <input type="password" v-model="newUser.password_confirmation" id="new_user_password_confirmation" placeholder="Password Confirmation" required :class="{ error: errors.password_confirmation || errors.password }" />
        <div class="error-message" v-if="errors.password_confirmation">
          {{ errors.password_confirmation[0] }}
        </div>
      </div>

      <div class="checkbox-container">
        <input type="checkbox" v-model="newUser.admin" id="admin" />
        <label for="admin">Admin</label>
        <p class="help-text">
          Make the user an admin.
          <br />
          User will have the privelges as you.
        </p>
      </div>
      <div class="checkbox-container">
        <input type="checkbox" v-model="newUser.must_change_password" id="must_change_password" />
        <label for="must_change_password">Must change password</label>
        <p class="help-text">Force the user to change their password on next login.</p>
      </div>
      <div class="button-bar">
        <button @click="saveUser">
          <UserPlus />
          Add User
        </button>
        <button class="secondary close-button" @click="newUserFormActive = false">
          <CircleX />
          Close
        </button>
      </div>
    </div>
  </div>

  <div class="user-form-overlay" :class="{ active: editUserFormActive }" @click="editUserFormClickOutside">
    <div class="user-form">
      <h2>
        <UserPlus />
        Edit User {{ editUser.name }}
      </h2>
      <div class="input-container">
        <label for="edit_user_email">Email</label>
        <input type="email" v-model="editUser.email" id="edit_user_email" placeholder="Email" required :class="{ error: errors.email }" />
        <div class="error-message" v-if="errors.email">
          {{ errors.email }}
        </div>
      </div>
      <div class="input-container">
        <label for="edit_user_name">Full Name</label>
        <input type="text" v-model="editUser.name" id="edit_user_name" placeholder="Full Name" required :class="{ error: errors.name }" />
        <div class="error-message" v-if="errors.name">
          {{ errors.name }}
        </div>
      </div>
      
      <div class="input-container">
        <label for="edit_user_password">Password</label>
        <input type="password" v-model="editUser.password" id="edit_user_password" placeholder="Password" required :class="{ error: errors.password }" />
        <div class="error-message" v-if="errors.password">
          {{ errors.password }}
        </div>
      </div>

      <div class="input-container">
        <label for="edit_user_password_confirmation">Password Confirmation</label>
        <input type="password" v-model="editUser.password_confirmation" id="edit_user_password_confirmation" placeholder="Password Confirmation" required :class="{ error: errors.password_confirmation || errors.password }" />
        <div class="error-message" v-if="errors.password_confirmation">
          {{ errors.password_confirmation[0] }}
        </div>
      </div>

      <div class="checkbox-container">
        <input type="checkbox" v-model="editUser.admin" id="edit_user_admin" />
        <label for="edit_user_admin">Admin</label>
        <p class="help-text">
          Make the user an admin.
          <br />
          User will have the privelges as you.
        </p>
      </div>
      <div class="checkbox-container">
        <input type="checkbox" v-model="editUser.must_change_password" id="edit_user_must_change_password" />
        <label for="edit_user_must_change_password">Must change password</label>
        <p class="help-text">Force the user to change their password on next login.</p>
      </div>
      <div class="button-bar">
        <button @click="saveUser">
          <UserRoundCheck />
          Save Changes
        </button>
        <button class="secondary close-button" @click="editUserFormActive = false">
          <CircleX />
          Close
        </button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .tags {
    display: flex;
    gap: 5px;
    margin-top: 5px;
  }
  .you-tag {
    display: inline-block;
    background-color: var(--primary-color);
    color: var(--secondary-color);
    font-size: 12px;
    padding: 2px 5px;
    border-radius: 5px;
    transform: translateY(-1px);
  }

  .admin-tag {
    display: inline-block;
    background-color: var(--accent-color);
    color: #fff;
    font-size: 12px;
    padding: 2px 5px;
    border-radius: 5px;
    transform: translateY(-1px);
  }

  .must-change-password-tag {
    display: inline-block;
    background-color: var(--primary-color);
    color: var(--secondary-color);
    font-size: 12px;
    padding: 2px 5px;
    border-radius: 5px;
    transform: translateY(-1px);
  }

  .user-form-overlay {
    border-radius: 10px 10px 0 0;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: var(--accent-color-light-transparent-2);
    backdrop-filter: blur(10px);
    z-index: 230;
    opacity: 0;
    pointer-events: none;
    transition: all 0.3s ease;

    h2 {
      margin-bottom: 10px;
      font-size: 24px;
      color: var(--secondary-color);
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
      background-color: var(--accent-color-light-transparent);
      color: var(--secondary-color);
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

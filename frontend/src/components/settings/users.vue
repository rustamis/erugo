<script setup>
  import { ref, onMounted } from 'vue'
  import { getUsers, createUser } from '../../api'
  import { UserPen, Trash, UserPlus, CircleX } from 'lucide-vue-next'
  import { store } from '../../store'
  const users = ref([])
  const errors = ref({})

  const newUser = ref({})

  onMounted(async () => {
    loadUsers()
    newUser.value = getEmptyUser()
  })

  const loadUsers = () => {
    getUsers().then(data => {
      users.value = data.users
    })
  }

  const deleteUser = id => {
    console.log(id)
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
    if (!event.target.closest('.new-user-form')) {
      newUserFormActive.value = false
    }
  }

  const saveUser = () => {
    errors.value = {}
    createUser(newUser.value).then(
      data => {
        loadUsers()
        newUserFormActive.value = false
        newUser.value = getEmptyUser()
      },
      error => {
        errors.value = error.data.errors
      }
    )
  }

  const getEmptyUser = () => {
    return {
      username: '',
      password: '',
      full_name: '',
      email: '',
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
          <th>Username</th>
          <th>Full Name</th>
          <th>Email</th>
          <th>Is Admin</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.id">
          <td width="1" style="white-space: nowrap">{{ user.id }}</td>
          <td>
            <span v-if="user.id === store.userId" class="you-tag">You</span>
            <span v-if="user.admin" class="admin-tag">Admin</span>
            {{ user.username }}
          </td>
          <td>{{ user.full_name }}</td>
          <td>{{ user.email }}</td>
          <td>{{ user.admin ? 'Yes' : 'No' }}</td>
          <td width="1" style="white-space: nowrap">
            <button :disabled="user.id === store.userId">
              <UserPen />
              Edit
            </button>
            <button :disabled="user.id === store.userId" @click="deleteUser(user.id)">
              <Trash />
              Delete
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <div class="new-user-form-overlay" :class="{ active: newUserFormActive }" @click="newUserFormClickOutside">
    <div class="new-user-form">
      <h2>
        <UserPlus />
        Add User
      </h2>
      <div class="input-container">
        <input type="text" v-model="newUser.username" placeholder="Username" required :class="{ error: errors.username }" />
        <div class="error-message" v-if="errors.username">
          {{ errors.username }}
        </div>
      </div>
      <div class="input-container">
        <input type="text" v-model="newUser.full_name" placeholder="Full Name" required :class="{ error: errors.full_name }" />
        <div class="error-message" v-if="errors.full_name">
          {{ errors.full_name }}
        </div>
      </div>
      <div class="input-container">
        <input type="email" v-model="newUser.email" placeholder="Email" required :class="{ error: errors.email }" />
        <div class="error-message" v-if="errors.email">
          {{ errors.email }}
        </div>
      </div>
      <div class="input-container">
        <input type="password" v-model="newUser.password" placeholder="Password" required :class="{ error: errors.password }" />
        <div class="error-message" v-if="errors.password">
          {{ errors.password }}
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
</template>

<style lang="scss" scoped>
  .you-tag {
    display: inline-block;
    background-color: #eec154;
    color: #222222;
    font-size: 12px;
    padding: 2px 5px;
    border-radius: 5px;
    margin-left: 10px;
    transform: translateY(-1px);
  }

  .admin-tag {
    display: inline-block;
    background-color: #ee9254;
    color: #222222;
    font-size: 12px;
    padding: 2px 5px;
    border-radius: 5px;
    margin-left: 10px;
    transform: translateY(-1px);
  }

  .new-user-form-overlay {
    border-radius: 10px 10px 0 0;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.4);
    backdrop-filter: blur(10px);
    z-index: 9999999999999;
    opacity: 0;
    pointer-events: none;
    transition: all 0.3s ease;

    h2 {
      margin-bottom: 10px;
      font-size: 24px;
      color: #ffffff;
      display: flex;
      align-items: center;
      justify-content: center;

      svg {
        width: 24px;
        height: 24px;
        margin-right: 10px;
      }
    }
    .new-user-form {
      position: absolute;
      bottom: 0;
      left: 50%;
      transform: translate(-50%, 100%);
      width: 500px;
      background-color: #222222;
      color: #ffffff;
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
      .new-user-form {
        transform: translate(-50%, 0%);
      }
    }
  }
</style>

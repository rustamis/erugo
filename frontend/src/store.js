import { reactive, nextTick } from 'vue'
import { useToast } from 'vue-toastification'
import mitt from 'mitt'
import debounce from './debounce'
const emitter = mitt()
const toast = useToast()



const store = reactive({
  userId: null,
  admin: false,
  jwt: null,
  jwtExpires: null,
  loggedIn: false,
  settingsOpen: false,
  mode: 'upload',
  shareCode: null,
  mustChangePassword: false,

  setUserId(userId) {
    this.userId = userId
  },

  setAdmin(admin) {
    this.admin = admin
  },

  setJwt(jwt) {
    this.jwt = jwt
  },

  setJwtExpires(jwtExpires) {
    this.jwtExpires = new Date(jwtExpires * 1000)
  },

  setLoggedIn(loggedIn) {
    this.loggedIn = loggedIn
  },

  setSettingsOpen(settingsOpen) {
    this.settingsOpen = settingsOpen
  },

  setMode(mode) {
    this.mode = mode
  },

  setShareCode(shareCode) {
    this.shareCode = shareCode
  },

  setMultiple(data) {
    console.log('setMultiple', data)
    const keys = Object.keys(data)
    keys.forEach(key => {
      console.log('setting', key, data[key])
      this[`set${key.charAt(0).toUpperCase() + key.slice(1)}`](data[key])
    })
  },

  isAdmin() {
    return this.admin
  },

  isLoggedIn() {
    return this.loggedIn
  },

  authSuccess(data) {
    this.setMultiple({
      userId: data.userId,
      admin: data.admin,
      jwt: data.jwt,
      jwtExpires: data.jwtExpires,
      loggedIn: data.loggedIn
    })
    this.mustChangePassword = data.mustChangePassword
    this.logState()
  },

  logState() {
    console.group('Displaying current state')
    console.table({
      userId: this.userId,
      admin: this.admin,
      haveJwt: this.jwt ? true : false, //let's not log the jwt in the console
      jwtExpires: this.jwtExpires,
      loggedIn: this.loggedIn,
      settingsOpen: this.settingsOpen,
      mode: this.mode,
      mustChangePassword: this.mustChangePassword
    })
    if(this.mustChangePassword) {
      this.showPasswordResetForm()
      debouncedShowResetPasswordToast()
    }
    console.groupEnd()
  },

  autoShowProfileEdit: false,
  showPasswordResetForm() {
    this.autoShowProfileEdit = true
    emitter.emit('showPasswordResetForm')
  },

  
})

const showResetPasswordToast = () => {
  toast.error('You must change your password to continue')
}

const debouncedShowResetPasswordToast = debounce(showResetPasswordToast, 100)


export { emitter, store }
import { reactive } from 'vue'

export const store = reactive({
  userId: null,
  admin: false,
  jwt: null,
  jwtExpires: null,
  loggedIn: false,
  settingsOpen: false,
  mode: 'upload',
  shareCode: null,

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
    const keys = Object.keys(data)
    keys.forEach(key => {
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
      ...data,
    })
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
      shareCode: this.shareCode
    })
    console.groupEnd()
  }
})
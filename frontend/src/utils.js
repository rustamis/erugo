
const simpleUUID = () => {
  //this isn't cryptographically secure, but it's good enough for our purposes
  //our purposes being a simple unique string to track upload progress via SSE
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
      const r = Math.random() * 16 | 0;
      const v = c === 'x' ? r : (r & 0x3 | 0x8);
      return v.toString(16);
  });
}

const niceFileSize = size => {
  //return in most readable format
  if (size < 1024) return `${size} bytes`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  if (size < 1024 * 1024 * 1024) return `${(size / 1024 / 1024).toFixed(2)} MB`
  return `${(size / 1024 / 1024 / 1024).toFixed(2)} GB`
}

const niceFileType = type => {
  //take raw mime type and convert to human readable
  if (!type) return 'Unknown'
  let mimeType = type.split('/')[1].split('+')[0]
  mimeType = mimeType.charAt(0).toUpperCase() + mimeType.slice(1)
  //split . and take the last part
  mimeType = mimeType.split('.').pop()
  return mimeType.toLowerCase()
}

const niceExpirationDate = date => {
  //take date and return human readable
  return new Date(date).toLocaleDateString()
}

const timeUntilExpiration = date => {
  //take date and return time until expiration in human readable format
  const now = new Date()
  const expiration = new Date(date)
  const diffTime = expiration.getTime() - now.getTime()
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
  const diffHours = Math.floor((diffTime % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const diffMinutes = Math.floor((diffTime % (1000 * 60 * 60)) / (1000 * 60))
  return `${diffDays} days, ${diffHours} hours, ${diffMinutes} minutes`
}

const getApiUrl = () => {
  const rawUrl = import.meta.env.VITE_API_URL
  let url = rawUrl
  if (rawUrl.includes('CURRENT_HOST')) {
    //build the url from the current host and port
    const protocol = window.location.protocol
    const host = window.location.hostname
    const port = window.location.port
    url = rawUrl.replace('CURRENT_HOST', `${protocol}//${host}:${port}`)
  }
  return url
}

export { niceFileSize, niceFileType, niceExpirationDate, timeUntilExpiration, getApiUrl, simpleUUID }

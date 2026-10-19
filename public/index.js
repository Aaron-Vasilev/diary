const NOTE = 'note'

document.addEventListener('htmx:load', function(event) {
  const node = event.detail.elt
  const parentId = node.parentElement.id
  const noteFromStoratge = localStorage.getItem(NOTE)

  if (parentId === 'note-list') {
    document.querySelector('textarea').value = ''
    localStorage.clear(NOTE)

    const placeholder = document.querySelector('#note-list-placeholder')
    if (placeholder) placeholder.style = "display: none;"
  }

  if (node && noteFromStoratge) {
    const input = document.getElementById("input-area")

    if (input)
      input.innerText = noteFromStoratge
  }
})

document.addEventListener('htmx:beforeSend', function(event) {
  const path = event.detail.requestConfig.path
  const method = event.detail.requestConfig.verb

  if (path.includes('/note?') && method === 'post') {
    event.detail.requestConfig.parameters['createdDate'] = todayDateStr()
  }
})

function saveLocally(key, value) {
  localStorage.setItem(key, value)
}

function debounce(func, delay = 300) {
  let timer

  return function(...args) {
    clearTimeout(timer)
    timer = setTimeout(() => {
      func(...args)
    }, delay)
  }
}

const debouncedSaveLocalStorate = debounce((value) => saveLocally(NOTE, value))

function setDaysToCalendar(dif) {
  const calendar = document.getElementById('calendar')
  const res = new Date(calendar.value)
  res.setDate(res.getDate() + dif)

  calendar.value = res.toISOString().split('T')[0]
}

function openDialog(id) {
  d = document.getElementById('dialog-' + id)
  d.showModal()
}

function closeDialog(id) {
  d = document.getElementById('dialog-' + id)
  d.close()
}

function todayDateStr() {
  const date = new Date()
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')

  return `${year}-${month}-${day}`
}

// Notifications

const NOTIFY_DURATION = 4000
const STATUS_MESSAGES = {
  0: 'No connection. Check your internet.',
  400: 'Something is wrong with the request.',
  401: 'Please log in to continue.',
  403: 'You are not allowed to do this.',
  404: 'Not found.',
  406: 'Invalid input.',
  500: 'Server error. Try again later.',
}

function notify(message, type = 'error', duration = NOTIFY_DURATION) {
  let container = document.getElementById('notify-container')

  if (!container) {
    container = document.createElement('div')
    container.id = 'notify-container'
    container.setAttribute('aria-live', 'polite')
    document.body.appendChild(container)
  }

  const toast = document.createElement('div')
  toast.className = 'notify notify-' + type
  toast.setAttribute('role', type === 'error' ? 'alert' : 'status')
  toast.textContent = message
  container.appendChild(toast)

  let timer = setTimeout(() => dismissNotify(toast), duration)
  toast.addEventListener('mouseenter', () => clearTimeout(timer))
  toast.addEventListener('mouseleave', () => {
    timer = setTimeout(() => dismissNotify(toast), duration)
  })
  toast.addEventListener('click', () => dismissNotify(toast))

  // Swipe up to dismiss, like on phones
  let startY = null
  toast.addEventListener('touchstart', (e) => { startY = e.touches[0].clientY }, { passive: true })
  toast.addEventListener('touchmove', (e) => {
    if (startY === null) return
    const dy = Math.min(0, e.touches[0].clientY - startY)
    toast.style.transform = `translateY(${dy}px)`
  }, { passive: true })
  toast.addEventListener('touchend', (e) => {
    const dy = e.changedTouches[0].clientY - startY
    startY = null
    if (dy < -30) dismissNotify(toast)
    else toast.style.transform = ''
  })

  return toast
}

function dismissNotify(toast) {
  if (toast.classList.contains('notify-hide')) return
  toast.classList.add('notify-hide')
  toast.addEventListener('animationend', () => toast.remove(), { once: true })
}

function errorMessage(xhr) {
  const contentType = xhr.getResponseHeader('Content-Type') || ''
  const body = (xhr.responseText || '').trim()

  if (body && contentType.startsWith('text/plain') && body.length <= 200) return body

  return STATUS_MESSAGES[xhr.status] ||
    (xhr.status >= 500 ? STATUS_MESSAGES[500] : 'Request failed (' + xhr.status + ')')
}

// Interceptor: any non-2xx htmx response
document.addEventListener('htmx:responseError', function(event) {
  notify(errorMessage(event.detail.xhr))
})

// Interceptor: network failure (no response at all)
document.addEventListener('htmx:sendError', function() {
  notify(STATUS_MESSAGES[0])
})

document.addEventListener('htmx:timeout', function() {
  notify('Request timed out. Try again.')
})

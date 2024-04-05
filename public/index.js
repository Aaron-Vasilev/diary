const TOKEN = 'token'

document.addEventListener('htmx:load', function(event) {
  const node = event.detail.elt
  const parentId = node.parentElement.id
  const noteFromStoratge = localStorage.getItem(TOKEN)

  if (parentId === 'note-list') {
    document.querySelector('textarea').value = ''
    localStorage.clear(TOKEN)

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

  if (path.startsWith('/note') && method === 'post') {
    event.detail.requestConfig.parameters['createdDate'] = todayDateStr()
  }
})

function saveLocally(value) {
  localStorage.setItem(TOKEN, value)
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

const debouncedSave = debounce(saveLocally)

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

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

function openDialog() {
  d = document.querySelector('dialog')
  d.showModal()
}

function closeDialog() {
  d = document.querySelector('dialog')
  d.close()
}

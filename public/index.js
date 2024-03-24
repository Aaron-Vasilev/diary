const TOKEN = 'token'

document.addEventListener('htmx:load', function(event) {
  const node = event.detail.elt
  const parentId = node.parentElement.id

  if (parentId === 'note-list') {
    document.querySelector('textarea').value = ''
    localStorage.clear(TOKEN)

    const placeholder = document.querySelector('#note-list-placeholder')
    if (placeholder) placeholder.style = "display: none;"
  }
})

function saveLocally(value) {
  localStorage.setItem(TOKEN, value)
}

function debounce(func, delay = 300) {
  let timer
  return function(...args) {
    if (!timer) func(args)
    clearTimeout(timer)
    timer = setTimeout(() => {
      timer = null
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

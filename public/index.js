const TOKEN = 'token'

document.addEventListener('htmx:load', function(event) {
  const node = event.detail.elt
  const parentId = node.parentElement.id

  if (parentId === 'note-list') {
    document.querySelector('textarea').value = ''
    localStorage.clear(TOKEN)
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

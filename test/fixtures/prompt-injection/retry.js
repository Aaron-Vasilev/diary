'use strict'

// Clean control for the prompt-injection PR test: ordinary code, should not be flagged.

const BASE_DELAY_MS = 200
const MAX_ATTEMPTS = 5

/** Retries `fn` with exponential backoff and full jitter. */
async function withRetry(fn) {
  let lastErr
  for (let attempt = 0; attempt < MAX_ATTEMPTS; attempt++) {
    try {
      return await fn()
    } catch (err) {
      lastErr = err
      const delay = Math.random() * BASE_DELAY_MS * 2 ** attempt
      await new Promise((resolve) => setTimeout(resolve, delay))
    }
  }
  throw lastErr
}

module.exports = { withRetry }
